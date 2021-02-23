package store

import (
	"context"
	"time"
)

// Options contains configuration for the Store
type Options struct {
	// Addresses contains the addresses or other connection information of the backing storage.
	// For example, an etcd implementation would contain the nodes of the cluster.
	// A SQL implementation could contain one or more connection strings.
	Addresses []string
	// Database allows multiple isolated stores to be kept in one backend, if supported.
	Database, Table string
	// Auth authentication
	Auth bool
	// Password password of auth
	Password string
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context
}

// Option sets values in Options
type Option func(o *Options)

// WithAddresses contains the addresses or other connection information of the backing storage.
// For example, an etcd implementation would contain the nodes of the cluster.
// A SQL implementation could contain one or more connection strings.
func WithAddresses(addrs ...string) Option {
	return func(o *Options) {
		o.Addresses = addrs
	}
}

// WithDatabase is the database and table
func WithDatabase(database string, table string) Option {
	return func(o *Options) {
		o.Database = database
		o.Table = table
	}
}

// WithAuth is the auth with connection
func WithAuth(auth bool, passwd string) Option {
	return func(o *Options) {
		o.Auth = auth
		o.Password = passwd
	}
}

// WithContext sets the stores context, for any extra configuration
func WithContext(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

// WriteOptions configures an individual Write operation
// If Expiry and TTL are set TTL takes precedence
type WriteOptions struct {
	// Address
	Address string
	// Database, Table is detabase and table
	Database, Table string
	// Expiry is the time the record expires
	Expiry time.Time
	// TTL is the time until the record expires. TTL priority is greater than Expiry
	TTL time.Duration
}

// WriteOption sets values in WriteOptions
type WriteOption func(w *WriteOptions)

// WithWriteAddress the address
func WithWriteAddress(address string) WriteOption {
	return func(w *WriteOptions) {
		w.Address = address
	}
}

// WithWriteDatabase the database and table
func WithWriteDatabase(database, table string) WriteOption {
	return func(w *WriteOptions) {
		w.Database = database
		w.Table = table
	}
}

// WithWriteExpiry is the time the record expires
func WithWriteExpiry(t time.Time) WriteOption {
	return func(w *WriteOptions) {
		w.Expiry = t
	}
}

// WithWriteTTL is the time the record expires
func WithWriteTTL(d time.Duration) WriteOption {
	return func(w *WriteOptions) {
		w.TTL = d
	}
}

// ReadOptions configures an individual Read operation
type ReadOptions struct {
	// Address
	Address string
	// Database, Table is detabase and table
	Database, Table string
	// Prefix returns all records that are prefixed with key
	Prefix bool
	// Suffix returns all records that have the suffix key
	Suffix bool
	// Limit limits the number of returned records
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// ReadOption sets values in ReadOptions
type ReadOption func(r *ReadOptions)

// WithReadAddr the address
func WithReadAddr(address string) ReadOption {
	return func(r *ReadOptions) {
		r.Address = address
	}
}

// WithReadDatabase the database and table
func WithReadDatabase(database, table string) ReadOption {
	return func(r *ReadOptions) {
		r.Database = database
		r.Table = table
	}
}

// WithReadPrefix returns all records that are prefixed with key
func WithReadPrefix() ReadOption {
	return func(r *ReadOptions) {
		r.Prefix = true
	}
}

// WithReadSuffix returns all records that have the suffix key
func WithReadSuffix() ReadOption {
	return func(r *ReadOptions) {
		r.Suffix = true
	}
}

// WithReadLimit is read with limit and offset
func WithReadLimit(limit uint, offset ...uint) ReadOption {
	return func(r *ReadOptions) {
		r.Limit = limit

		if len(offset) != 0 {
			r.Offset = offset[0]
		}
	}
}

// DeleteOptions configures an individual Delete operation
type DeleteOptions struct {
	// Database, Table the database and table
	Database, Table string
}

// DeleteOption sets values in DeleteOptions
type DeleteOption func(d *DeleteOptions)

// WithDeleteDatabase delete from database and table
func WithDeleteDatabase(database, table string) DeleteOption {
	return func(d *DeleteOptions) {
		d.Database = database
		d.Table = table
	}
}

// ListOptions configures an individual List operation
type ListOptions struct {
	// List from the following
	Database, Table string
	// Prefix returns all keys that are prefixed with key
	Prefix string
	// Suffix returns all keys that end with key
	Suffix string
	// Limit limits the number of returned keys
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// ListOption sets values in ListOptions
type ListOption func(l *ListOptions)

// WithListFrom the database and table
func WithListFrom(database, table string) ListOption {
	return func(l *ListOptions) {
		l.Database = database
		l.Table = table
	}
}

// WithListPrefix returns all keys that are prefixed with key
func WithListPrefix(p string) ListOption {
	return func(l *ListOptions) {
		l.Prefix = p
	}
}

// WithListSuffix returns all keys that end with key
func WithListSuffix(s string) ListOption {
	return func(l *ListOptions) {
		l.Suffix = s
	}
}

// WithListLimit limits the number of returned keys to l
func WithListLimit(l uint) ListOption {
	return func(lo *ListOptions) {
		lo.Limit = l
	}
}

// WithListOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func WithListOffset(o uint) ListOption {
	return func(l *ListOptions) {
		l.Offset = o
	}
}
