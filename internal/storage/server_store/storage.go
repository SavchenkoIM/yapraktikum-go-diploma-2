// Database Postgres storage for service data

package server_store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"passwordvault/internal/config"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/utils"
	"sync"
	"time"
)

//////////////////////////
// Storage
//////////////////////////

var ErrUserAlreadyExists error = errors.New("this login already exists in database")
var ErrUserAuthFailed error = errors.New("authentication failed")
var ErrUserNotLoggedIn error = errors.New("user session has expired")
var ErrNoDataAffected error = errors.New("no data was affected. Check request parameters")
var ErrUnimplemented error = errors.New("unexpected request parameters")

// Postgres storage data
type Storage struct {
	dbConn      *pgxpool.Pool
	config      *config.ServerConfig
	encKey      string
	logger      *zap.Logger
	rStarter    sync.Once          // First Init call detector
	workersWg   *sync.WaitGroup    // WaitGroup for Storage Workers
	stopWorkers context.CancelFunc // Cancel function for Storage Workers Context
	workersCtx  context.Context    // Storage Workers Context
}

// Constructor of storage object
func New(config *config.ServerConfig, logger *zap.Logger) (*Storage, error) {
	encKey := make([]byte, 128)
	_, err := rand.Read(encKey)
	if err != nil {
		return nil, err
	}

	s := Storage{
		config:      config,
		logger:      logger,
		encKey:      hex.EncodeToString(encKey),
		stopWorkers: nil,
		workersCtx:  nil,
		workersWg:   &sync.WaitGroup{},
	}

	poolConfig, err := pgxpool.ParseConfig(s.config.DBConnectionString)
	if err != nil {
		s.logger.Sugar().Errorf("Unable to parse connection string: %s", err)
		return nil, err
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	s.dbConn, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		s.logger.Sugar().Errorf("Unable to create connection pool: %s", err)
		return nil, err
	}
	return &s, nil
}

// Initializes storage object before use
func (s *Storage) Init(ctx context.Context) error {
	firstInit := false
	s.rStarter.Do(func() {
		firstInit = true
	})

	if firstInit {
		s.workersCtx, s.stopWorkers = context.WithCancel(ctx)
	}

	var err error
	errs := make([]error, 0)

	queries := []string{
		queryCreateExtensionUUID,
		queryCreateExtensionPGCrypto,
		queryCreateUsers,
		getCreateDataQuery(proto.DataType_CREDENTIALS),
		getCreateDataQuery(proto.DataType_CREDIT_CARD),
		getCreateDataQuery(proto.DataType_TEXT_NOTE),
		getCreateDataQuery(proto.DataType_BLOB),
		getCreateMetaDataQuery(proto.DataType_CREDENTIALS),
		getCreateMetaDataQuery(proto.DataType_CREDIT_CARD),
		getCreateMetaDataQuery(proto.DataType_TEXT_NOTE),
		getCreateMetaDataQuery(proto.DataType_BLOB),
	}

	for _, query := range queries {
		_, err = s.dbConn.Exec(ctx, query)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if firstInit {
		s.workersWg.Add(1)
		go s.autoInit(s.workersCtx)
	}

	return errors.Join(errs...)
}

// Automatically checks for ping loss and reinitializes storage after database is recovered
func (s *Storage) autoInit(ctx context.Context) {
	defer func() { s.workersWg.Done() }()
	connPrev := true
	connected := false
	cw := utils.NewCtxCancelWaiter(ctx, 10*time.Second)

	for {
		if cw.Scan() != nil {
			s.logger.Info("autoInit worker stopped")
			return
		}
		if connected = s.dbConn.Ping(ctx) == nil; connected && !connPrev {
			err := s.Init(ctx)
			if err != nil {
				s.logger.Sugar().Errorf("Initialization error: %s", err.Error())
			} else {
				s.logger.Sugar().Warnf("Database restored after fault.")
			}
		}
		connPrev = connected
	}
}

// Shuts storage down
func (s *Storage) Close(ctx context.Context) {
	s.logger.Info("Stopping storage workers...")
	s.stopWorkers()
	s.workersWg.Wait()
	s.dbConn.Close()
}
