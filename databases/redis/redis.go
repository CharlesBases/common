package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"

	"charlesbases/common/log"
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
		Addr: address,
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
