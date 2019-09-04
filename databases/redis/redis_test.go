package redis

import (
	"fmt"
	"testing"
)

var (
	addr = "192.168.1.88:6379"
)

type Peo struct {
	Name string
	Age  int
}

func TestRedis(t *testing.T) {
	Redis := GetRedis(addr)
	// ++++++++++++++++++++++ //
	err := Redis.Set("name", "张三")
	if err != nil {
		fmt.Println(err)
		return
	}
	str, err := Redis.Get("name")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("name:  ", str)
	Redis.Del("name")
	// ++++++++++++++++++++++ //
	err = Redis.Set("peo", &Peo{"李四", 18})
	if err != nil {
		fmt.Println(err)
		return
	}
	var peo Peo
	_, err = Redis.Get("peo", &peo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("peo:   ", peo)
	Redis.Del("peo")
	// ++++++++++++++++++++++ //
	err = Redis.Set("name", "李四", 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	str, err = Redis.Get("name")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("name:  ", str)
	Redis.Del("name")
}

func TestRedisPoll(t *testing.T) {

	redisPoll := GetRedisPool()
	defer redisPoll.Close()

	redisClient := redisPoll.GetRedisClient()
	defer redisClient.Close()

	err := redisClient.Set("name", "张三")
	if err != nil {
		fmt.Println(err)
		return
	}
	str, err := redisClient.Get("name")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("name:  ", str)
	redisClient.Del("name")
}
