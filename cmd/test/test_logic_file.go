package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"passwordvault/internal/uni_client"
	"path/filepath"
	"testing"
)

func testLogicFile(ctx context.Context, t *testing.T, client *uni_client.UniClient) {

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
		err := client.DownloadFile(ctx, "test_file")
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
