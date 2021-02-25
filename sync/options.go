package sync

import (
	"time"
)

type Options struct {
	Addresses []string
	Prefix    string
	Auth      bool
	Password  string

	// Blocked 阻塞 | 非阻塞
	Blocked bool
	TTL     time.Duration
}

type Option func(o *Options)

// WithAddresses sets the addresses to use
func WithAddresses(addresses ...string) Option {
	return func(o *Options) {
		o.Addresses = addresses
	}
}

// WithTTL set the timeout
func WithTTL(d time.Duration) Option {
	return func(o *Options) {
		o.TTL = d
	}
}

// WithAuth is the auth with connection
func WithAuth(auth bool, passwd string) Option {
	return func(o *Options) {
		o.Auth = auth
		o.Password = passwd
	}
}

// WithBlocked .
func WithBlocked() Option {
	return func(o *Options) {
		o.Blocked = true
	}
}

// WithPrefixPrefix sets a prefix to any lock ids used
func WithPrefix(p string) Option {
	return func(o *Options) {
		o.Prefix = p
	}
}

type LockOptions struct {
	TTL  time.Duration
	Wait time.Duration
}

type LockOption func(o *LockOptions)

// WithLockTTL sets the lock ttl
func WithLockTTL(t time.Duration) LockOption {
	return func(o *LockOptions) {
		o.TTL = t
	}
}

// WithLockWait sets the wait time
func WithLockWait(t time.Duration) LockOption {
	return func(o *LockOptions) {
		o.Wait = t
	}
}

// Leader provides leadership election
type Leader interface {
	// resign leadership
	Resign() error
	// status returns when leadership is lost
	Status() chan bool
}

type LeaderOptions struct{}

type LeaderOption func(o *LeaderOptions)
