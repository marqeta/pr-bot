package adapters

import (
	"errors"
	"time"
)

var (
	// ErrNotStored is returned when an Add operation failed to respect its condtions
	ErrNotStored = errors.New("adapter: not stored")

	// ErrNotFound is returned when attempting to Get a key that does not exist
	ErrNotFound = errors.New("adapter: not found")
)

// Adapter defines the interface for memchached client adapters to be implemented
type Adapter interface {

	// Add will attempt to add a new item in memcached.
	// If the item already exists returns ErrNotStored.
	Add(key string, value string, expiration time.Duration) error

	// Get an existing item from memcached
	// returns ErrNotFound if the correspondent key does not exist
	Get(key string) (value string, err error)

	// Delete an existing item in memcached.
	// returns ErrNotFound if the correspondent key does not exist
	Delete(key string) error
}
