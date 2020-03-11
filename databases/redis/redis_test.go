package redis

import (
	"fmt"
	"testing"

	"charlesbases/common/log"
)

var (
	addr = "192.168.1.174:6379"
)

type Peo struct {
	Name string
	Age  int
}

func TestRedis(t *testing.T) {
	defer log.Flush()
	var err error

	Redis := GetRedis(addr)

	// ++++++++++++++++++++++ //
	// err = Redis.Set(STRING, "name", "张三")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// var name string
	// err = Redis.Get(STRING, "name", &name)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(name)

	// ++++++++++++++++++++++ //
	// err = Redis.Set(HASH, "peo", &Peo{
	// 	Name: "张三",
	// 	Age:  18,
	// }, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// var peo Peo
	// fmt.Println(Redis.Exists("peo"))
	// err = Redis.Get(HASH, "peo", &peo)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(peo)

	// ++++++++++++++++++++++ //
	err = Redis.Del("a")
	if err != nil {
		fmt.Println(err)
	}
}

/*func TestRedisPoll(t *testing.T) {

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
}*/
