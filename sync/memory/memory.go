package memory

import (
	gosync "sync"
	"time"

	"github.com/CharlesBases/common/sync"
)

// NewSync new sync of memory
func NewSync(opts ...sync.Option) sync.Sync {
	var options sync.Options
	for _, o := range opts {
		o(&options)
	}

	return &memorySync{
		options: options,
		locks:   make(map[string]*memoryLock),
	}
}

type memorySync struct {
	options sync.Options

	mtx   gosync.RWMutex
	locks map[string]*memoryLock
}

type memoryLock struct {
	id        string
	release   chan bool
	createdAt time.Time
	expiresAt time.Time
}

// load load from memory
func (m *memorySync) load(id string) (*memoryLock, bool) {
	m.mtx.Lock()
	lk, isExist := m.locks[id]
	m.mtx.Unlock()
	return lk, isExist
}

// Init init options
func (m *memorySync) Init(opts ...sync.Option) error {
	for _, o := range opts {
		o(&m.options)
	}
	return nil
}

// Options return Options
func (m *memorySync) Options() sync.Options {
	return m.options
}

// Lock .
func (m *memorySync) Lock(id string, opts ...sync.LockOption) error {
	var options sync.LockOptions
	for _, o := range opts {
		o(&options)
	}

	if m.options.Prefix != "" {
		id = m.options.Prefix + id
	}

	var timeout = time.Second * 3
	if m.options.Timeout != 0 {
		timeout = m.options.Timeout
	}

	var ttl = timeout
	if options.TTL != 0 {
		ttl = options.TTL
	}

lockloop:
	m.mtx.Lock()
	lk, isExist := m.locks[id]
	switch isExist {
	case true:
		// 锁未过期
		if time.Now().Before(lk.expiresAt) {
			m.mtx.Unlock()

			select {
			// 锁已释放
			case <-lk.release:
				goto lockloop
			// 超时
			case <-time.Tick(timeout):
				return sync.ErrLockTimeout
			}
		}
		// 锁已过期
		fallthrough
	default:
		m.locks[id] = &memoryLock{
			id:        id,
			release:   make(chan bool, 1),
			createdAt: time.Now(),
			expiresAt: time.Now().Add(ttl),
		}
		m.mtx.Unlock()
		return nil
	}
}

// Unlock .
func (m *memorySync) Unlock(id string) error {
	m.mtx.Lock()

	lk, isExist := m.locks[id]
	if isExist {
		delete(m.locks, id)

		select {
		case <-lk.release:
			break
		default:
			close(lk.release)
		}
	}

	m.mtx.Unlock()
	return nil
}

// String .
func (m *memorySync) String() string {
	return "memory"
}
