package object

import (
	"context"
	"io"
	"webserver/services"

	"github.com/minio/minio-go/v7"
)

type store struct {
	s3         *minio.Client
	bucketName string
}

func NewStore(c *minio.Client, bucketName string) services.ObjectStore {
	return &store{c, bucketName}
}

func (s *store) PutObject(id string, r io.Reader, size int64) (int64, error) {
	inf, err := s.s3.PutObject(context.Background(), s.bucketName, id, r, size, minio.PutObjectOptions{})

	return inf.Size, err
}

func (s *store) GetObject(id string) (io.ReadSeekCloser, int64, error) {
	obj, err := s.s3.GetObject(context.Background(), s.bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}

	inf, err := obj.Stat()
	if err != nil {
		return nil, 0, err
	}

	return obj, inf.Size, err
}

func (s *store) DeleteObject(id string) error {
	return s.s3.RemoveObject(context.Background(), s.bucketName, id, minio.RemoveObjectOptions{})
}

func (s *store) HasObject(id string) bool {
	_, err := s.s3.StatObject(context.Background(), s.bucketName, id, minio.GetObjectOptions{})
	return err == nil
}
