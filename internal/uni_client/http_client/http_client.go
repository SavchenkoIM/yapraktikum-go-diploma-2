package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"passwordvault/internal/config"
	"path/filepath"
)

type HTTPClient struct {
	cfg       *config.ClientConfig
	logger    *zap.Logger
	token     string
	sendError chan error
}

func NewHTTPClient(cfg *config.ClientConfig, logger *zap.Logger) *HTTPClient {
	return &HTTPClient{
		cfg:       cfg,
		logger:    logger,
		sendError: make(chan error, 1),
	}
}

func (c *HTTPClient) getTlsClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ClientAuth:         tls.NoClientCert, // Server provides cert
				InsecureSkipVerify: true,             // Any server_store cert is accepted
			},
		},
	}
}

func (c *HTTPClient) SetToken(token string) {
	c.token = token
}

func (c *HTTPClient) DownloadFile(ctx context.Context, fileName string) error {
	r, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s/download", c.cfg.AddressHTTP),
		bytes.NewBuffer([]byte(fmt.Sprintf(`{ "filename": "%s" }`, fileName))))
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.token))

	resp, err := c.getTlsClient().Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("error downloading file, status: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	fFileName := ""
	if filepath.IsAbs(fileName) {
		fFileName = fileName
	} else {
		fFileName = filepath.Join(c.cfg.FilesDefaultDir, fileName)
	}

	file, err := os.OpenFile(fFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) UploadFile(ctx context.Context, fileName string) error {

	fFileName := ""
	if filepath.IsAbs(fileName) {
		fFileName = fileName
	} else {
		fFileName = filepath.Join(c.cfg.FilesDefaultDir, fileName)
	}

	file, err := os.Open(fFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	r, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s/upload", c.cfg.AddressHTTP), file)
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.token))
	r.Header.Add("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, filepath.Base(fFileName)))

	resp, err := c.getTlsClient().Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
