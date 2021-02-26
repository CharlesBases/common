package sync

import (
	"time"
)

type Options struct {
	Addresses []string
	Prefix    string

	// Auth authentication
	Auth     bool
	Password string

	// Blocked 阻塞 | 非阻塞
	Blocked bool
	// Timeout default 3 second
	Timeout time.Duration
}

type Option func(o *Options)

// WithAddresses sets the addresses to use
func WithAddresses(addresses ...string) Option {
	return func(o *Options) {
		o.Addresses = addresses
	}
}

// WithPrefixPrefix sets a prefix to any lock ids used
func WithPrefix(p string) Option {
	return func(o *Options) {
		o.Prefix = p
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

// WithTimeout set the timeout
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

type LockOptions struct {
	TTL time.Duration
}

type LockOption func(o *LockOptions)

// WithLockTTL sets the lock ttl
func WithLockTTL(d time.Duration) LockOption {
	return func(o *LockOptions) {
		o.TTL = d
	}
}
