package redis

import (
	"fmt"
	"log"
	gosync "sync"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/CharlesBases/common/sync"
)

// defaultFormat default format
const defaultFormat = "2006-01-02 15:04:05"

// NewStore returns a redis sync
func NewStore(opts ...sync.Option) sync.Sync {
	var options sync.Options
	for _, o := range opts {
		o(&options)
	}

	s := new(redisSync)
	s.options = options

	if err := s.configure(); err != nil {
		log.Fatal(err)
	}

	return s
}

// redisSync redis sync
type redisSync struct {
	options sync.Options
	client  *redis.Client

	mtx   gosync.RWMutex
	locks map[string]*redisLock
}

// redisLock redis lock
type redisLock struct {
	id  string
	ttl time.Duration
}

// redisLeader redis leader
type redisLeader struct {
}

func (r *redisSync) configure() error {
	var redisOptions *redis.Options
	addrs := r.options.Addresses

	if len(addrs) == 0 {
		addrs = []string{"redis://127.0.0.1:6379"}
	}

	redisOptions, err := redis.ParseURL(addrs[0])
	if err != nil {
		// Backwards compatibility
		redisOptions = &redis.Options{
			Addr:     addrs[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}

	if r.options.Auth {
		redisOptions.Password = r.options.Password
	}

	r.client = redis.NewClient(redisOptions)
	return nil
}

// lock .
func (r *redisSync) lock(id string, ttl time.Duration) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	locked, err := r.client.SetNX(id, time.Now().Format(defaultFormat), ttl).Result()
	if err != nil || !locked {
		return false
	}
	return true
}

func (r *redisLeader) Resign() error {
	return nil
}

func (r *redisLeader) Status() chan bool {
	return nil
}

func (r *redisSync) Leader(id string, opts ...sync.LeaderOption) (sync.Leader, error) {
	return nil, nil
}

func (r *redisSync) Init(opts ...sync.Option) error {
	for _, o := range opts {
		o(&r.options)
	}
	return nil
}

func (r *redisSync) Options() sync.Options {
	return r.options
}

func (r *redisSync) Lock(id string, opts ...sync.LockOption) error {
	var options sync.LockOptions
	for _, o := range opts {
		o(&options)
	}

	var ttl = time.Second * 3
	if options.TTL != 0 {
		ttl = options.TTL
	}

	if r.options.Prefix != "" {
		id = r.options.Prefix + id
	}

	switch r.options.Blocked {
	case false:
		if r.lock(id, ttl) {
			return nil
		}
		return fmt.Errorf("lock %[1]s failed: %[1]s's locking", id)
	case true:
		for {
			select {
			case <-time.Tick(ttl):
				return fmt.Errorf("lock %s failed: timeout", id)
			default:
				if r.lock(id, ttl) {
					return nil
				}
			}
		}
	}

	return nil
}

func (r *redisSync) Unlock(id string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.options.Prefix != "" {
		id = r.options.Prefix + id
	}

	affected, err := r.client.Del(id).Result()
	if err != nil || affected == 0 {
		log.Fatal(fmt.Sprintf(`unlock %[1]s failed: %[1]s's unlocked`, id))
		return nil
	}
	return nil
}

func (r *redisSync) String() string {
	return "redis"
}
