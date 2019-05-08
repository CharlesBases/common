package db

import (
	"encoding/json"
	"os"

	"common/log"

	"github.com/gomodule/redigo/redis"
)

var (
	Redis redis.Conn
)

func InitRedis(address string) redis.Conn {
	openRedis(address)
	return Redis
}

func openRedis(address string) {
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		log.Error(" - Redis连接失败 - ", err.Error())
		os.Exit(0)
	}

	Redis = conn
}

// key seconds
func SetKeyExpire(key string, seconds int) error {
	_, err := Redis.Do("EXPIRE", key, seconds)
	return err
}

// del key
func DelKey(key string) error {
	_, err := Redis.Do("DEL", key)
	return err
}

// set key-value for string
func SetKey(key string, value string) error {
	_, err := Redis.Do("SET", key, value)
	return err
}

// get key-value for string
func GetKey(key string) (value string, err error) {
	return redis.String(Redis.Do("GET", key))
}

// set key-value for interface
func SetKeyInterface(key string, s interface{}) error {
	jsonByte, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = Redis.Do("SET", key, string(jsonByte))
	return err
}

// get key-value for interface
func GetKeyInterface(key string, value interface{}) error {
	jsonStr, err := redis.String(Redis.Do("GET", key))
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), value)
}
