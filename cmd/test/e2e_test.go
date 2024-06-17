package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"passwordvault/cmd/test/testcontainers"
	"passwordvault/internal/config"
	"passwordvault/internal/grpc_server"
	"passwordvault/internal/http_server"
	"passwordvault/internal/storage/server_store"
	"passwordvault/internal/uni_client"
	"testing"
	"time"
)

var minioEndpoint string
var dbConnectionString string

func Test_E2E(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

	var cMinio *testcontainers.MinioContainer
	var cPostgres *testcontainers.PostgresContainer

	t.Run("Setting_Up_Test_Environment", func(t *testing.T) {
		if _, ok := os.LookupEnv("GITHUB_TEST_RUN"); ok {
			minioEndpoint = "minio:9000"
			dbConnectionString = "postgresql://postgres:postgres@postgres/postgres?sslmode=disable"
		} else {
			cMinio, err = testcontainers.NewMinioContainer()
			require.NoError(t, err)
			minioEndpoint, err = cMinio.EndPoint()
			require.NoError(t, err)
			cPostgres, err = testcontainers.NewPostgresContainer()
			require.NoError(t, err)
			dbConnectionString, err = cPostgres.ConnectionString()
			require.NoError(t, err)
		}
	})

	if cMinio != nil && cPostgres != nil {
		defer cMinio.Close()
		defer cPostgres.Close()
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	t.Logf("Minio: %s, Postgres: %s", minioEndpoint, dbConnectionString)

	t.Run("Starting_Server", func(t *testing.T) {
		assert.NoError(t, err)

		srvCfg := config.ServerConfig{
			GrpcEndPoint:       "localhost:8081",
			HttpEndPoint:       "localhost:8080",
			MinioEndPoint:      minioEndpoint,
			MinioAdminId:       "minioadmin",
			MinioAdminKey:      "minioadmin",
			DBConnectionString: dbConnectionString,
			Key:                "secret",
			UseKey:             true,
			CertFileName:       "../../data/cert/cert.pem",
			PKFileName:         "../../data/cert/priv.pem",
		}

		db, err := server_store.New(&srvCfg, logger)
		require.NoError(t, err)
		err = db.Init(ctx)
		require.NoError(t, err)

		gSrv := grpc_server.NewGRPCServer(db, &srvCfg, logger)
		gSrv.ListenAndServeAsync()

		hSrv := http_server.NewHttpServer(ctx, db, &srvCfg, logger)
		hSrv.ListenAndServeAsync()
	})

	t.Run("Wait_2_Seconds_For_Server_To_Start", func(t *testing.T) {
		time.Sleep(time.Second * 2)
	})

	var uCli *uni_client.UniClient
	t.Run("Init_Client", func(t *testing.T) {
		// Init client
		cliCfg := config.ClientConfig{
			AddressGRPC:     "localhost:8081",
			AddressHTTP:     "localhost:8080",
			FilesDefaultDir: "test_filestore_dir",
		}

		uCli = uni_client.NewUniClient(logger, cliCfg)
		uCli.Start(ctx)
	})

	testLogicUser(ctx, t, uCli)
	testLogicDataWrite(ctx, t, uCli)
	testLogicFile(ctx, t, uCli)
	testLogicDataCheck(ctx, t, uCli)
	testLogicDataDelete(ctx, t, uCli)

}
