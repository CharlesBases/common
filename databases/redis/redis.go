package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/CharlesBases/common/log"
)

const (
	STRING = iota + 1
	HASH
	LIST
	SET
	ZSET
)

var (
	wg sync.WaitGroup
)

type Redis struct {
	client *redis.Client
}

func GetRedis(address string) *Redis {
	redisClient := redis.NewClient(&redis.Options{
		// 连接信息
		Network:  "tcp",            // 网络类型，tcp or unix，默认tcp
		Addr:     "127.0.0.1:6379", // 主机名+冒号+端口，默认localhost:6379
		Username: "",               // 用户名
		Password: "",               // 密码
		DB:       0,                // redis数据库index

		// 重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		// 超时策略
		DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时

		// 连接池容量及闲置连接数量
		PoolSize:     15,              // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,              // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		MaxConnAge:   0 * time.Second, // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
		PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
		IdleTimeout:  5 * time.Minute, // 闲置超时，默认5分钟，-1表示取消闲置超时检查

		// 闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		TLSConfig:          nil,
		Limiter:            nil,

		// 可自定义连接函数
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial(network, addr)
		},
		// 钩子函数. 仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		OnConnect: func(conn *redis.Conn) error {
			fmt.Printf("conn=%v\n", conn)
			return nil
		},
	})
	return &Redis{client: redisClient}
}

func (r *Redis) Set(DataType int, key string, value interface{}, seconds ...int64) error {
	args := make([]interface{}, 0)
	switch DataType {
	case STRING:
		args = append(args,
			"SET",
			key,
			func() []byte {
				bytes, _ := json.Marshal(value)
				return bytes
			}())

		for _, second := range seconds {
			args = append(args,
				"EX",
				strconv.FormatInt(second, 10))
			break
		}

		return r.client.Do(args...).Err()
	case HASH:
		args = append(args,
			"HMSET",
			key,
			key,
			func() []byte {
				bytes, _ := json.Marshal(value)
				return bytes
			}())

		if err := r.client.Do(args...).Err(); err != nil {
			return err
		}

		for _, second := range seconds {
			go func() {
				err := r.Expire(key, second)
				if err != nil {
					log.Error(fmt.Sprintf("redis expire err: [key: %s] >> %s", key, err.Error()))
				}
			}()
			break
		}

		return nil
	case LIST:
	case SET:
	case ZSET:
	}

	return fmt.Errorf(fmt.Sprintf("redis data type unknown: %d", DataType))
}

func (r *Redis) Get(DataType int, key string, value interface{}) error {
	switch DataType {
	case STRING:
		bytes, err := r.client.Get(key).Bytes()
		if err != nil {
			log.Error(fmt.Sprintf("redis get error: [key: %s] >> %s", key, err.Error()))
			return err
		}
		return json.Unmarshal(bytes, value)
	case HASH:
		result, err := r.client.HGetAll(key).Result()
		if err != nil {
			log.Error(fmt.Sprintf("redis get error: [key: %s] >> %s", key, err.Error()))
			return err
		}
		return json.Unmarshal([]byte(result[key]), value)
	case LIST:
	case SET:
	case ZSET:
	}
	return fmt.Errorf(fmt.Sprintf("redis data type unknown: %d", DataType))
}

func (r *Redis) Exists(key string) bool {
	isExist := r.client.Exists(key).Val()
	if isExist != 0 {
		return true
	}
	return false
}

func (r *Redis) Expire(key string, second int64) error {
	return r.client.Do("EXPIRE", key, second).Err()
}

func (r *Redis) Del(key string) (delerr error) {
	// return r.client.Unlink(key...).Err()

	wg.Add(1)

	go func() {
		newkey := fmt.Sprintf("%s_delete_%d", key, time.Now().UnixNano())
		err := r.client.RenameNX(key, newkey).Err()
		if err != nil {
			delerr = fmt.Errorf("redis delete[renamenx] err: [key: %s] >> %s", key, err.Error())
		} else {
			go func() {
				if err := r.client.Del(newkey).Err(); err != nil {
					log.Error(fmt.Sprintf("redis delete[delete] err: [key: %s] >> %s", newkey, err.Error()))
				}
			}()
		}

		wg.Done()

	}()

	wg.Wait()

	return
}

func (r *Redis) Close() error {
	return r.client.Close()
}
