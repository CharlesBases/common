package pool

import (
	"fmt"
	"sync"
	"time"
)

type Config struct {
	MaxConn int
	MinConn int
	Factory func() (interface{}, error)
	Close   func(interface{}) error
	Ping    func(interface{}) error
	Timeout time.Duration
}

type Pool interface {
	Get() (interface{}, error)
	Put(interface{}) error
	Close(interface{}) error
	Release()
	Len() int
}

type pool struct {
	mutex   sync.RWMutex
	conns   chan *conn
	factory func() (interface{}, error)
	close   func(interface{}) error
	ping    func(interface{}) error
	timeout time.Duration
}

type conn struct {
	connect interface{}
	time    time.Time
}

// NewPool init pool
func NewPool(config *Config) (Pool, error) {
	if config.MinConn < 0 || config.MaxConn <= 0 || config.MinConn > config.MaxConn {
		return nil, fmt.Errorf("invalid capacity settings")
	}
	if config.Factory == nil {
		return nil, fmt.Errorf("invalid factory func settings")
	}
	if config.Close == nil {
		return nil, fmt.Errorf("invalid close func settings")
	}
	pool := &pool{
		conns:   make(chan *conn, config.MaxConn),
		factory: config.Factory,
		close:   config.Close,
		timeout: config.Timeout,
	}
	if config.Ping != nil {
		pool.ping = config.Ping
	}
	for i := 0; i < config.MinConn; i++ {
		connect, err := pool.factory()
		if err != nil {
			pool.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		pool.conns <- &conn{connect: connect, time: time.Now()}
	}
	return pool, nil
}

// get get all conn
func (pool *pool) get() chan *conn {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	conns := pool.conns
	return conns
}

// Get take a connection from pool
func (pool *pool) Get() (interface{}, error) {
	conns := pool.get()
	if conns == nil {
		return nil, fmt.Errorf("pool is closed")
	}
	for {
		select {
		case conn := <-conns:
			if conn == nil {
				return nil, fmt.Errorf("pool is closed")
			}
			// timeout
			if timeout := pool.timeout; timeout > 0 {
				if conn.time.Add(timeout).Before(time.Now()) {
					pool.Close(conn.connect)
					continue
				}
			}
			// ping
			if pool.ping != nil {
				if err := pool.Ping(conn.connect); err != nil {
					fmt.Println("connect is not able to be connected: ", err)
					continue
				}
			}
			return conn.connect, nil
		default:
			pool.mutex.RLock()
			defer pool.mutex.RUnlock()

			if pool.factory == nil {
				pool.mutex.RUnlock()
				continue
			}
			connect, err := pool.factory()
			if err != nil {
				return nil, err
			}
			return connect, nil
		}
	}
}

// Put place the connection back in the pool
func (pool *pool) Put(connect interface{}) error {
	if connect == nil {
		return fmt.Errorf("connection is nil. rejecting")
	}
	pool.mutex.RLock()
	if pool.conns == nil {
		pool.mutex.RUnlock()
		return pool.Close(connect)
	}
	select {
	case pool.conns <- &conn{connect: connect, time: time.Now()}:
		pool.mutex.RUnlock()
		return nil
	default:
		pool.mutex.RUnlock()
		return pool.Close(connect)
	}
}

// Close close single connection
func (pool *pool) Close(connect interface{}) error {
	if connect == nil {
		return fmt.Errorf("connection is nil. refuse")
	}

	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	if pool.close == nil {
		return nil
	}
	return pool.close(connect)
}

// Ping
func (pool *pool) Ping(connect interface{}) error {
	if connect == nil {
		return fmt.Errorf("connection is nil. refuse")
	}
	return pool.ping(connect)
}

// Release release all connections in connection pool
func (pool *pool) Release() {
	pool.mutex.RLock()

	conns := pool.conns
	poolclose := pool.close

	pool.factory = nil
	pool.conns = nil
	pool.close = nil
	pool.ping = nil

	pool.mutex.RUnlock()
	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		poolclose(conn.connect)
	}
}

// Len
func (pool *pool) Len() int {
	return len(pool.get())
}
