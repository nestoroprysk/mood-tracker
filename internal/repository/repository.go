package repository

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// Repository interacts with the storage.
type Repository interface {
	Read(obj string) ([]byte, error)
	Override(obj string, body []byte) error
}

type repository struct {
	*storage.BucketHandle
}

var _ Repository = &repository{}

func New(bucket string) (Repository, error) {
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	return repository{
		BucketHandle: client.Bucket(bucket),
	}, nil
}

func (r repository) Read(obj string) ([]byte, error) {
	rc, err := r.Object(obj).NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func (r repository) Override(obj string, body []byte) error {
	wc := r.Object(obj).NewWriter(context.TODO())
	_, err := wc.Write(body)
	if err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
