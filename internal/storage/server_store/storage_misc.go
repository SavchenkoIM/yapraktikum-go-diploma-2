package server_store

import (
	"context"
)

// Returns minio password for logged user
func (s *Storage) GetFileStoreKey(ctx context.Context) (string, error) {
	userId := ctx.Value("LoggedUserId").(string)

	var minioKey string
	query := `SELECT pgp_sym_decrypt(filestore_access_key, $1) FROM users WHERE id = $2`
	err := s.dbConn.QueryRow(ctx, query, s.config.Key, userId).Scan(&minioKey)
	if err != nil {
		return "", err
	}

	return minioKey, nil
}
