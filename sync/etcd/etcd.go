package etcd

import (
	"context"
	"errors"
	"log"
	"path"
	"strings"
	gosync "sync"

	client "github.com/coreos/etcd/clientv3"
	cc "github.com/coreos/etcd/clientv3/concurrency"

	"github.com/CharlesBases/common/sync"
)

// NewSync return a etcd sync
func NewSync(opts ...sync.Option) sync.Sync {
	var options sync.Options
	for _, o := range opts {
		o(&options)
	}

	var endpoints []string

	for _, addr := range options.Addresses {
		if len(addr) > 0 {
			endpoints = append(endpoints, addr)
		}
	}

	if len(endpoints) == 0 {
		endpoints = []string{"http://127.0.0.1:2379"}
	}

	// TODO: parse addresses
	c, err := client.New(client.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &etcdSync{
		path:    "/micro/sync",
		client:  c,
		options: options,
		locks:   make(map[string]*etcdLock),
	}
}

type etcdSync struct {
	options sync.Options
	path    string
	client  *client.Client

	mtx   gosync.Mutex
	locks map[string]*etcdLock
}

type etcdLock struct {
	s *cc.Session
	m *cc.Mutex
}

// Init init option
func (e *etcdSync) Init(opts ...sync.Option) error {
	for _, o := range opts {
		o(&e.options)
	}
	return nil
}

// Options return options
func (e *etcdSync) Options() sync.Options {
	return e.options
}

// Lock lock id
func (e *etcdSync) Lock(id string, opts ...sync.LockOption) error {
	var options sync.LockOptions
	for _, o := range opts {
		o(&options)
	}

	// make path
	path := path.Join(e.path, strings.Replace(e.options.Prefix+id, "/", "-", -1))

	var sopts []cc.SessionOption
	if options.TTL > 0 {
		sopts = append(sopts, cc.WithTTL(int(options.TTL.Seconds())))
	}

	s, err := cc.NewSession(e.client, sopts...)
	if err != nil {
		return err
	}

	m := cc.NewMutex(s, path)

	if err := m.Lock(context.TODO()); err != nil {
		return err
	}

	e.mtx.Lock()
	e.locks[id] = &etcdLock{
		s: s,
		m: m,
	}
	e.mtx.Unlock()
	return nil
}

// Unlock unlock id
func (e *etcdSync) Unlock(id string) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	v, ok := e.locks[id]
	if !ok {
		return errors.New("lock not found")
	}
	err := v.m.Unlock(context.Background())
	delete(e.locks, id)
	return err
}

// String .
func (e *etcdSync) String() string {
	return "etcd"
}
