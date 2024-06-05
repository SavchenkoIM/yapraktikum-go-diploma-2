package server_store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgerrcode"
	uuid "github.com/jackc/pgx-gofrs-uuid"
	"golang.org/x/crypto/scrypt"
	"io"
	"passwordvault/internal/storage/file_store"
	"passwordvault/internal/utils"
	"strings"
)

func (s *Storage) UserRegister(ctx context.Context, login string, password string) error {

	salt := make([]byte, 32) // salt, 32 bytes len
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}
	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, 256) // hash, 256 bytes len
	if err != nil {
		return err
	}

	minioPassword := utils.GeneratePassword(25, true, true, true)

	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (login, password, salt, filestore_access_key) VALUES ($1, $2, $3, pgp_sym_encrypt($4, $5)) RETURNING id`

	var newUuid uuid.UUID

	if err = tx.QueryRow(ctx, query, login, hex.EncodeToString(hash), hex.EncodeToString(salt),
		minioPassword,
		s.config.Key).Scan(&newUuid); err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			s.logger.Sugar().Errorf("Login %s already exists in database", login)
			return fmt.Errorf("%s: %w", err.Error(), ErrUserAlreadyExists)
		}
		tx.Rollback(ctx)
		return err
	}

	uuidV, _ := newUuid.UUIDValue()

	mCli := file_store.NewMinioStorage(s.config.MinioEndPoint, "", s.config.MinioAdminId, s.config.MinioAdminKey)
	if err = mCli.UserReg(ctx, fmt.Sprintf("%x", uuidV.Bytes), minioPassword); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Storage) UserCheckLoggedIn(token string) (string, error) {
	ac := utils.AuthClaims{}
	err := ac.SetFromJWT(token, s.encKey)
	if err != nil {
		return "", ErrUserNotLoggedIn
	}

	return ac.UserID, nil
}

func (s *Storage) UserLogin(ctx context.Context, login string, password string) (string, error) {

	var err error
	query := `SELECT id, login, password, salt FROM users WHERE login=$1`
	row := s.dbConn.QueryRow(ctx, query, login)
	var (
		sUserID string
		sLogin  string
		sPassw  string
		sSalt   string
	)
	if err = row.Scan(&sUserID, &sLogin, &sPassw, &sSalt); err != nil {
		return "", err
	}

	xSalt, _ := hex.DecodeString(sSalt)

	var key []byte
	if key, err = scrypt.Key([]byte(password), xSalt, 1<<14, 8, 1, 256); err != nil {
		return "", err
	}

	if hex.EncodeToString(key) != sPassw {
		return "", ErrUserAuthFailed
	}

	ac := utils.AuthClaims{UserID: sUserID}
	jwt, err := ac.GetJWT(s.encKey)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
