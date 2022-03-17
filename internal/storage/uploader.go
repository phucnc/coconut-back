package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"io"
	"os"
)

type Config struct {
	BucketName string
	CredentialFile string
}

func LoadConfigFromEnv() *Config {
	c := &Config{
		BucketName:     os.Getenv("BUCKET_NAME"),
		CredentialFile: os.Getenv("CREDENTIAL_FILE"),
	}
	return c
}

type Uploader struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewUploader(ctx context.Context, config *Config) (*Uploader, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.CredentialFile))
	if err != nil {
		return nil, err
	}
	uploader := &Uploader{
		client: client,
		bucket: client.Bucket(config.BucketName),
	}
	return uploader, nil
}

func (u *Uploader) Upload(ctx context.Context, fileName string, fileSrc io.Reader) (string, error) {
	// Read the object1 from bucket.
	dest := u.bucket.Object(fileName).NewWriter(ctx)

	// Upload
	if _, err := io.Copy(dest, fileSrc); err != nil {
		return "", errors.Wrap(err, "io.Copy()")
	}

	// Finish upload
	err := dest.Close()
	if err != nil {
		return "", errors.Wrap(err, "Close()")
	}

	// Set ACL
	acl := u.bucket.Object(fileName).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", errors.Wrap(err, "acl.Set()")
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", os.Getenv("BUCKET_NAME"), fileName), nil
}
