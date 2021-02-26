// Package sync is an interface for distributed synchronization
package sync

import (
	"errors"
)

var (
	// ErrLockTimeout timeout
	ErrLockTimeout = errors.New("lock timeout")
)

// Sync is an interface for distributed synchronization
type Sync interface {
	// Initialise options
	Init(...Option) error
	// Return the options
	Options() Options
	// Lock acquires a lock
	Lock(id string, opts ...LockOption) error
	// Unlock releases a lock
	Unlock(id string) error
	// Sync implementation
	String() string
}
