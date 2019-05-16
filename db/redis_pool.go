package db

import "github.com/gomodule/redigo/redis"

var address string

var (
	MaxActive int  // 最大连接数
	MaxIdle   int  // 最大空闲连接数
	Wait      bool // true
)

type RedisPool struct {
	*redis.Pool
}

func GetRedisPool() *RedisPool {
	return &RedisPool{
		&redis.Pool{
			MaxActive: MaxActive,
			MaxIdle:   MaxIdle,
			Wait:      Wait,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial("tcp", address)
			},
		},
	}
}

func (pool *RedisPool) GetRedisClient() *Redis {
	return &Redis{pool.Get()}
}
