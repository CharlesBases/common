package db

import (
	"encoding/json"
	"strconv"

	"github.com/gomodule/redigo/redis"

	"github.com/CharlesBases/common/log"
)

type Redis struct {
	redis.Conn
}

func GetRedis(address string) *Redis {
	conn, err := redis.Dial("tcp", address)
	if err != nil || conn == nil {
		log.Error(" - Redis连接失败 - ", err.Error())
		return nil
	}
	return &Redis{conn}
}

func (r *Redis) SetKeyExpire(key string, seconds int) error {
	_, err := r.Do("EXPIRE", key, seconds)
	return err
}

func (r *Redis) Set(key string, value interface{}, seconds ...int) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	args := []interface{}{
		key,
		string(bs),
	}
	for _, v := range seconds {
		args = append(args, "EX")
		args = append(args, strconv.Itoa(v))
	}
	_, err = r.Do("SET", args...)
	return err
}

func (r *Redis) Get(key string, values ...interface{}) (string, error) {
	jsonStr, err := redis.String(r.Do("GET", key))
	if err != nil {
		return "", err
	}
	for k := range values {
		return "", json.Unmarshal([]byte(jsonStr), values[k])
	}
	return jsonStr, err
}

func (r *Redis) Del(key string) error {
	_, err := r.Do("DEL", key)
	return err
}
