package testcontainers

import (
	"context"
	"fmt"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

// Docker container with Minio
type MinioContainer struct {
	instance testcontainers.Container
}

// Constructor for MinioContainer
func NewMinioContainer() (*MinioContainer, error) {

	//time.Sleep(time.Duration(rand.Int31n(10)*5000) * time.Millisecond) // Temp workaround for go test ./... bug

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	//g := CustomLogConsumer{}
	testcontainers.Logger = log.New(&ioutils.NopWriter{}, "", 0)
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio",
		ExposedPorts: []string{"9000/tcp"},
		Env:          map[string]string{},
		WaitingFor:   wait.ForListeningPort("9000/tcp"),
		Cmd:          []string{"server", "/data"},
		Entrypoint:   []string{"minio"},
		/*LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts:      []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10 * time.Second)},
			Consumers: []testcontainers.LogConsumer{&g},
		},*/
	}
	minio, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}
	return &MinioContainer{
		instance: minio,
	}, nil
}

// Address:port for containerized Minio instance
func (db *MinioContainer) EndPoint() (string, error) {
	port, err := db.Port()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("127.0.0.1:%d", port), nil
}

// Returns mapped Minio 9000 port
func (db *MinioContainer) Port() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "9000")
	if err != nil {
		return 0, err
	}
	return p.Int(), nil
}

// Destructor for MinioContainer
func (db *MinioContainer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return db.instance.Terminate(ctx)
}

// Returns containerized Minio host
func (db *MinioContainer) Host() string {
	return "localhost"
}
