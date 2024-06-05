package file_store

import (
	"context"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinioStorage struct {
	address    string
	bucketName string
	userId     string
	accessKey  string
	secure     bool
}

func NewMinioStorage(address string, bucketName string, userId string, accessKey string) *MinioStorage {
	return &MinioStorage{
		address:    address,
		bucketName: bucketName,
		userId:     userId,
		accessKey:  accessKey,
		secure:     false,
	}
}

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
