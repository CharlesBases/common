package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/CharlesBases/common/store"
)

// NewStore returns a redis store
func NewStore(opts ...store.Option) store.Store {
	var options store.Options
	for _, o := range opts {
		o(&options)
	}

	s := new(rkv)
	s.options = options

	if err := s.configure(); err != nil {
		log.Fatal(err)
	}

	return s
}

type rkv struct {
	options store.Options
	client  *redis.Client
}

func (r *rkv) configure() error {
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

func (r *rkv) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&r.options)
	}

	return r.configure()
}

func (r *rkv) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	options := store.ReadOptions{}
	options.Table = r.options.Table

	for _, o := range opts {
		o(&options)
	}

	rkey := fmt.Sprintf("%s%s", options.Table, key)

	var keys = make(map[string]struct{}, 0)
	keys[rkey] = struct{}{}

	// Prefix
	if options.Prefix {
		prefixKey := rkey + "*"
		pkeys, err := r.client.Keys(prefixKey).Result()
		if err != nil {
			return nil, err
		}
		for _, i := range pkeys {
			keys[i] = struct{}{}
		}
	}

	// Suffix
	if options.Suffix {
		suffixKey := "*" + rkey
		skeys, err := r.client.Keys(suffixKey).Result()
		if err != nil {
			return nil, err
		}
		for _, i := range skeys {
			keys[i] = struct{}{}
		}
	}

	records := make([]*store.Record, 0, len(keys))

	for rkey = range keys {
		val, err := r.client.Get(rkey).Bytes()

		if err != nil && err == redis.Nil {
			return nil, store.ErrNotFound
		} else if err != nil {
			return nil, err
		}

		if val == nil {
			return nil, store.ErrNotFound
		}

		d, err := r.client.TTL(rkey).Result()
		if err != nil {
			return nil, err
		}

		records = append(records, &store.Record{
			Key:    key,
			Value:  val,
			TTL:    d,
			Expiry: time.Now().Add(d),
		})
	}

	return records, nil
}

func (r *rkv) Delete(key string, opts ...store.DeleteOption) error {
	options := store.DeleteOptions{}
	options.Table = r.options.Table

	for _, o := range opts {
		o(&options)
	}

	rkey := fmt.Sprintf("%s%s", options.Table, key)
	return r.client.Del(rkey).Err()
}

func (r *rkv) Write(record *store.Record, opts ...store.WriteOption) error {
	options := store.WriteOptions{}
	options.Table = r.options.Table

	for _, o := range opts {
		o(&options)
	}

	if len(opts) > 0 {
		if !options.Expiry.IsZero() {
			record.TTL = time.Since(options.Expiry)
			record.Expiry = options.Expiry
		}
		if options.TTL != 0 {
			record.TTL = options.TTL
			record.Expiry = time.Now().Add(options.TTL)
		}
	}

	rkey := fmt.Sprintf("%s%s", options.Table, record.Key)
	return r.client.Set(rkey, record.Value, record.TTL).Err()
}

func (r *rkv) List(opts ...store.ListOption) ([]string, error) {
	options := store.ListOptions{}
	options.Table = r.options.Table

	for _, o := range opts {
		o(&options)
	}

	var pattern string
	if options.Prefix != "" {
		pattern = options.Prefix + "*"
	} else if options.Suffix != "" {
		pattern = "*" + options.Suffix
	} else {
		pattern = "*"
	}

	keys, err := r.client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *rkv) Close() error {
	return r.client.Close()
}

func (r *rkv) String() string {
	return "redis"
}

func (r *rkv) Options() store.Options {
	return r.options
}
