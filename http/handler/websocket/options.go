package websocket

import "time"

// Options .
type Options struct {
	Namespace   string        // 命名空间
	Buffer      int           // 缓冲区间. default: 32
	ContentType string        // 内容类型.
	Auth        bool          // 认证
	Timeout     time.Duration // 超时
}

type Option func(o *Options)

// WithNamespace .
func WithNamespace(namespace string) Option {
	return func(o *Options) {
		o.Namespace = namespace
	}
}

// WithContentType .
func WithContentType(t string) Option {
	return func(o *Options) {
		o.ContentType = t
	}
}

// WithBuffer .
func WithBuffer(cap int) Option {
	return func(o *Options) {
		o.Buffer = cap
	}
}

// WithAuth .
func WithAuth() Option {
	return func(o *Options) {
		o.Auth = true
	}
}

// WithTimeout .
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}
