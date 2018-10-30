package storage

import (
	"io"

	"github.com/oracle/mysql-operator/pkg/backup/storage/s3"
)

// Interface abstracts the underlying storage provider.
type Interface interface {
	// Store creates a new object in the underlying provider's datastore if it does not exist,
	// or replaces the existing object if it does exist.
	Store(key string, body io.ReadCloser) error
	// Retrieve return the object in the underlying provider's datastore if it exists.
	Retrieve(key string) (io.ReadCloser, error)
}

// NewStorageProvider accepts a secret map and uses its contents to determine the
// desired object storage provider implementation.
func NewStorageProvider(config v1alpha1.StorageProvider, credentials map[string]string) (Interface, error) {
	return s3.NewProvider(config.S3, credentials)
}
