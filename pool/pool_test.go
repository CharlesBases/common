package pool

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func Test(t *testing.T) {
	config := &Config{
		MinConn: 5,                                                                        // 最小连接
		MaxConn: 30,                                                                       // 最大连接
		Factory: func() (interface{}, error) { return net.Dial("tcp", "127.0.0.1:8080") }, // 创建连接池
		Close:   func(v interface{}) error { return v.(net.Conn).Close() },                // 关闭连接
		Ping:    func(v interface{}) error { return nil },                                 // ping
		Timeout: 15 * time.Second,                                                         // 超时时间
	}

	pool, err := NewPool(config)
	if err != nil {
		fmt.Println("init pool error: ", err)
	}

	// 从连接池中取得一个连接
	conn, err := pool.Get()

	// do something

	// 将连接放回连接池中
	pool.Put(conn)

	// 释放连接池中的所有连接
	pool.Release()

	// 查看当前连接中的数量
	pool.Len()
}
