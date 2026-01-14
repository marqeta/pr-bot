// Package gomemcache implements the adapters.Adapter interface for https://github.com/bradfitz/gomemcache/blob/master/memcache/memcache.go
package gomemcache

import (
	"github.com/alessandro-c/gomemcached-lock/adapters"
	"github.com/bradfitz/gomemcache/memcache"
	"math"
	"time"
)

type client struct {
	mc *memcache.Client
}

// New creates a new bradfitz/gomemcache implementation of lock.Client
func New(mc *memcache.Client) adapters.Adapter {
	return &client{mc}
}

// Add implements adapters.Add
func (c *client) Add(key string, value string, expiration time.Duration) (err error) {
	err = c.mc.Add(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(math.Abs(time.Since(time.Now().Add(expiration)).Seconds())),
	})
	// mapping adapters errors
	if err == memcache.ErrNotStored {
		err = adapters.ErrNotStored
	}
	return
}

// Get implements adapters.Add
func (c *client) Get(key string) (string, error) {
	var value string
	item, err := c.mc.Get(key)
	if err != nil {
		// mapping adapters errors
		if err == memcache.ErrCacheMiss {
			err = adapters.ErrNotFound
		}
	} else {
		value = string(item.Value)
	}
	return value, err
}

// Delete implements adapters.Delete
func (c *client) Delete(key string) error {
	err := c.mc.Delete(key)
	if err == memcache.ErrCacheMiss {
		// mapping adapters error
		err = adapters.ErrNotFound
	}
	return err
}
