// Package lock provides a simple lock/release library to be used on a single memcached instance.
//
// This package is compatible with https://github.com/bradfitz/gomemcache but
// really any client will do, just implement the ./adapters.Adapter interface.
package lock

import (
	"errors"
	"fmt"
	"github.com/alessandro-c/gomemcached-lock/adapters"
	"github.com/thanhpk/randstr"
	"time"
)

var (
	// ErrNotAcquired is returned when a locker tries to lock an already locked resource
	ErrNotAcquired = errors.New("locker: not acquired")

	// ErrNotFound is returned when a lock does not exist in memcached
	ErrNotFound = errors.New("locker: not found")

	// ErrForbidden is returned when an owner attempts to release a non-owned locked
	ErrForbidden = errors.New("locker: forbidden")
)

// Locker is the main entrypoint for locking operations
type Locker struct {
	c     adapters.Adapter
	name  string
	owner string
}

// New creates a new Locker instance
func New(c adapters.Adapter, name, owner string) *Locker {
	if len(owner) == 0 {
		owner = randstr.String(8)
	}
	return &Locker{
		name:  name,
		owner: owner,
		c:     c,
	}
}

// Lock attempts to lock for the given TTL.
// returns ErrNotAcquired if the resource was already locked.
func (l *Locker) Lock(ttl time.Duration) (err error) {
	err = l.c.Add(l.name, l.owner, ttl)
	if err == nil {
		// it looks like the lock was obtained but
		// sometimes, under heavy load, it is possible
		// for 2 or more "add" operations to succeed.
		// enforcing ownership even after successful locks
		// will increase locking reliability.
		owner, _ := l.GetCurrentOwner()
		if owner != l.owner {
			// RACE CONDITION! leave the lock to the actual owner
			err = ErrNotAcquired
		}
	} else if err == adapters.ErrNotStored {
		err = ErrNotAcquired
	}
	return
}

// Release attempts to release a lock
// return ErrNotFound if the lock does not exist
// return ErrForbidden if the locker does not own the lock
func (l *Locker) Release() (err error) {
	// attempts to retrieve current owner
	owner, err := l.GetCurrentOwner()
	if err == nil {
		// enforce lock ownership
		if l.owner == owner {
			err = l.c.Delete(l.name)
		} else {
			err = ErrForbidden
		}
	} else if err == ErrNotFound {
		err = ErrNotFound
	} else {
		err = fmt.Errorf("locker: release - something went wrong : '%s'", err.Error())
	}
	return
}

// GetCurrentOwner returns the current lock owner
// return ErrNotFound if lock does not exist
func (l *Locker) GetCurrentOwner() (string, error) {
	owner, err := l.c.Get(l.name)
	if err == adapters.ErrNotFound {
		return "", ErrNotFound
	}
	return owner, nil
}
