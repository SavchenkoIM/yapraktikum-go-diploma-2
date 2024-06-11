package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"os"
	"passwordvault/cmd/test/testcontainers"
	"passwordvault/internal/config"
	"passwordvault/internal/grpc_server"
	"passwordvault/internal/http_server"
	"passwordvault/internal/storage/server_store"
	"passwordvault/internal/uni_client"
	"path/filepath"
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
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			minioEndpoint, err = cMinio.EndPoint()
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			cPostgres, err = testcontainers.NewPostgresContainer()
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			dbConnectionString, err = cPostgres.ConnectionString()
			if err != nil {
				assert.FailNow(t, err.Error())
			}
		}
	})

	if cMinio != nil && cPostgres != nil {
		defer cMinio.Close()
		defer cPostgres.Close()
	}

	logger, err := zap.NewDevelopment()

	logger.Sugar().Infof("Minio: %s, Postgres: %s", minioEndpoint, dbConnectionString)

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
		if err != nil {
			logger.Fatal(err.Error())
		}
		err = db.Init(ctx)
		if err != nil {
			logger.Error(err.Error())
		}

		gSrv := grpc_server.NewGRPCServer(db, &srvCfg, logger)
		gSrv.ListenAndServeAsync()

		hSrv := http_server.NewHttpServer(ctx, db, &srvCfg, logger)
		hSrv.ListenAndServeAsync()
	})

	t.Run("Wait_5_Seconds_For_Server_To_Start", func(t *testing.T) {
		time.Sleep(time.Second * 5)
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

	testLogic(ctx, t, uCli)
}

func testLogic(ctx context.Context, t *testing.T, client *uni_client.UniClient) {
	var err error

	t.Run("Unregistered_User_Login", func(t *testing.T) {
		_, err = client.UserLogin(ctx, "Victoria", "Victoria's secret")
		assert.Error(t, err)
	})

	t.Run("User_Create", func(t *testing.T) {
		_, err = client.UserCreate(ctx, "Victoria", "Victoria's secret")
		assert.NoError(t, err)
	})

	t.Run("Registered_User_Login_Wrong_Pass", func(t *testing.T) {
		_, err = client.UserLogin(ctx, "Victoria", "Victoria secret")
		assert.Error(t, err)
	})

	t.Run("Registered_User_Login", func(t *testing.T) {
		_, err = client.UserLogin(ctx, "Victoria", "Victoria's secret")
		assert.NoError(t, err)
	})

	fileOrig := "test_filestore_dir/document.test"
	testString := "this is test document"
	t.Run("Upload_File", func(t *testing.T) {
		os.MkdirAll(filepath.Dir(fileOrig), os.ModePerm)
		wrFile, err := os.OpenFile(fileOrig, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		defer func() {
			err = wrFile.Close()
			assert.NoError(t, err)
			err = os.Remove(fileOrig)
			assert.NoError(t, err)
		}()
		assert.NoError(t, err)
		_, err = wrFile.WriteString(testString)
		assert.NoError(t, err)
		err = client.UploadFile(ctx, "test_file", filepath.Base(fileOrig))
		assert.NoError(t, err)
	})

	t.Run("Download_File", func(t *testing.T) {
		err = client.DownloadFile(ctx, "test_file")
		assert.NoError(t, err)
		file1, err := os.OpenFile(fileOrig, os.O_RDONLY, os.ModePerm)
		assert.NoError(t, err)
		defer func() {
			err = file1.Close()
			assert.NoError(t, err)
			err = os.RemoveAll(filepath.Dir(fileOrig))
			assert.NoError(t, err)
		}()
		c1, err := io.ReadAll(file1)
		assert.NoError(t, err)
		assert.Equal(t, string(c1), testString)
	})
}
