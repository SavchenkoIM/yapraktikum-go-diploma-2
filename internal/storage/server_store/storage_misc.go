package server_store

import (
	"context"
)

func getLoggedUserId(ctx context.Context) (string, error) {
	id, ok := ctx.Value("LoggedUserId").(string)
	if !ok {
		return "", ErrUserNotLoggedIn
	}
	return id, nil
}

// Returns minio password for logged user
func (s *Storage) GetFileStoreKey(ctx context.Context) (string, error) {
	userId, err := getLoggedUserId(ctx)
	if err != nil {
		return "", err
	}

	var minioKey string
	query := `SELECT pgp_sym_decrypt(filestore_access_key, $1) FROM users WHERE id = $2`
	err = s.dbConn.QueryRow(ctx, query, s.config.Key, userId).Scan(&minioKey)
	if err != nil {
		return "", err
	}

	return minioKey, nil
}
