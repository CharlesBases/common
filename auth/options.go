package auth

import "time"

// Options .
type Options struct {
	// PrivateKey private key for generate token
	PrivateKey string
	// TTL default 4 Hours
	TTL time.Duration
}

// Option
type Option func(o *Options)

// WithTTL .
func WithTTL(d time.Duration) Option {
	return func(o *Options) {
		o.TTL = d
	}
}

// WithPrivateKey .
func WithPrivateKey(private string) Option {
	return func(o *Options) {
		o.PrivateKey = private
	}
}

// GenOptions
type GenOptions struct {
	// Metadata metadata with the account
	Metadata map[string]interface{}
	// Expiry is the time the token expires
	Expiry time.Time
	// TTL is the time until the token expires. 过期时间优先级排序: Expiry > TTL > Option.TTL
	TTL time.Duration
}

type GenOption func(o *GenOptions)

// WithGenExpiry .
func WithGenExpiry(t time.Time) GenOption {
	return func(o *GenOptions) {
		o.Expiry = t
	}
}

// WithGenTTL .
func WithGenTTL(d time.Duration) GenOption {
	return func(o *GenOptions) {
		o.TTL = d
	}
}

// ParOptions .
type ParOptions struct{}

type ParOption func(o *ParOptions)
