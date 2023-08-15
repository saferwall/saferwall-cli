// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package minio

import (
	"bytes"
	"context"
	"io"

	mio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	// ErrNoSuchKey is raised when the object is not found.
	ErrNoSuchKey = "The specified key does not exist."
)

// Service provides abstraction to cloud object storage.
type Service struct {
	// s3 client.
	client *mio.Client
}

func New(endpoint, accessKey, secretKey string) (Service, error) {

	// New returns an Amazon S3 compatible client object.
	// API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := mio.New(endpoint, &mio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return Service{}, nil
	}

	return Service{s3Client}, nil
}

// Upload uploads an object to s3.
func (s Service) Upload(ctx context.Context, bucket, key string,
	file io.Reader) error {

	// Get the size.
	buf := &bytes.Buffer{}
	size, err := io.Copy(buf, file)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(ctx, bucket, key, buf, size,
		mio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Download(ctx context.Context, bucket, key string,
	file io.Writer) error {

	reader, err := s.client.GetObject(ctx, bucket, key, mio.GetObjectOptions{})
	if err != nil {
		return err
	}

	stat, err := reader.Stat()
	if err != nil {
		return err
	}

	_, err = io.CopyN(file, reader, stat.Size)
	if err != nil {
		return err
	}
	return nil
}

// MakeBucket creates a new bucket with bucketName with a context to control
// cancellations and timeouts.
func (s Service) MakeBucket(ctx context.Context, bucketName,
	location string) error {
	err := s.client.MakeBucket(ctx, bucketName,
		mio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket
		// (which happens if you run this twice)
		exists, errBucketExists := s.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}

// Exists checks whether an object exists already in the object storage.
func (s Service) Exists(ctx context.Context, bucketName,
	key string) (bool, error) {
	opts := mio.GetObjectOptions{}
	_, err := s.client.StatObject(ctx, bucketName, key, opts)
	if err != nil {
		switch err.Error() {
		case ErrNoSuchKey:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// List enumerates the list of objects in a bucket.
func (s Service) List(ctx context.Context, bucketName string) ([]string, error) {
	var objKeys []string

	opts := mio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
	}

	for object := range s.client.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		objKeys = append(objKeys, object.Key)

	}

	return objKeys, nil
}

// Delete deletes an object from the remote storage.
func (s Service) Delete(ctx context.Context, bucket, key string) error {

	opts := mio.RemoveObjectOptions{}

	err := s.client.RemoveObject(ctx, bucket, key, opts)
	if err != nil {
		return err
	}

	return nil
}
