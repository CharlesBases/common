package db

import (
	"fmt"
	"testing"
)

var (
	address = "192.168.1.88:6379"
)

type Peo struct {
	Name string
	Age  int
}

func TestRedis(t *testing.T) {
	Redis := GetRedis(address)
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

	err = Redis.Set("peo", &Peo{Name: "李四", Age: 18})
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
}
