package file_store

import (
	"context"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

// Minio Storage access object
type MinioStorage struct {
	address    string
	bucketName string
	userId     string
	accessKey  string
	secure     bool
}

// Constructs minio storage access object
func NewMinioStorage(ctx context.Context, address string, bucketName string, userId string, accessKey string) (*MinioStorage, error) {
	ms := &MinioStorage{
		address:    address,
		bucketName: bucketName,
		userId:     userId,
		accessKey:  accessKey,
		secure:     false,
	}

	client, err := minio.New(ms.address, &minio.Options{
		Creds:  credentials.NewStaticV4(ms.userId, ms.accessKey, ""),
		Secure: ms.secure,
	})
	if err != nil {
		return nil, err
	}
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return ms, nil
}

// Minio user registration
func (ms *MinioStorage) UserReg(ctx context.Context, userId string, accessKey string) error {
	mdmClnt, err := madmin.NewWithOptions(ms.address, &madmin.Options{
		Creds:     credentials.NewStaticV4(ms.userId, ms.accessKey, ""),
		Secure:    ms.secure,
		Transport: nil,
	})
	if err != nil {
		return err
	}

	err = mdmClnt.AddUser(ctx, userId, accessKey)
	if err != nil {
		return err
	}

	_, err = mdmClnt.AttachPolicy(ctx, madmin.PolicyAssociationReq{
		Policies: []string{"readwrite"},
		User:     userId,
		Group:    "",
	})
	if err != nil {
		mdmClnt.RemoveUser(ctx, userId)
		return err
	}

	return nil
}

// Download file from minio
func (ms *MinioStorage) Download(ctx context.Context, fileName string) (io.Reader, error) {
	minioClient, err := minio.New(ms.address, &minio.Options{
		Creds:  credentials.NewStaticV4(ms.userId, ms.accessKey, ""),
		Secure: ms.secure,
	})
	if err != nil {
		return nil, err
	}
	obj, err := minioClient.GetObject(ctx, ms.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Upload file to minio
func (ms *MinioStorage) Upload(ctx context.Context, reader io.Reader, fileName string) error {
	minioClient, err := minio.New(ms.address, &minio.Options{
		Creds:  credentials.NewStaticV4(ms.userId, ms.accessKey, ""),
		Secure: ms.secure,
	})
	if err != nil {
		return err
	}
	_, err = minioClient.PutObject(ctx, ms.bucketName, fileName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Delete file
func (ms *MinioStorage) Delete(ctx context.Context, fileName string) error {
	minioClient, err := minio.New(ms.address, &minio.Options{
		Creds:  credentials.NewStaticV4(ms.userId, ms.accessKey, ""),
		Secure: ms.secure,
	})
	if err != nil {
		return err
	}
	err = minioClient.RemoveObject(ctx, ms.bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
