package redis

import (
	"log"
	gosync "sync"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/CharlesBases/common/sync"
)

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
	return nil
}

func (r *redisSync) Unlock(id string) error {
	return nil
}

func (r *redisSync) String() string {
	return "redis"
}
