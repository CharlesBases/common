package redis

// var address string
//
// var (
// 	MaxActive int  // 最大连接数
// 	MaxIdle   int  // 最大空闲连接数
// 	Wait      bool // true
// )
//
// type RedisPool struct {
// 	*redis.Pool
// }
//
// func GetRedisPool() *RedisPool {
// 	return &RedisPool{
// 		&redis.Pool{
// 			IdleTimeout: time.Hour,
// 			MaxActive:   MaxActive,
// 			MaxIdle:     MaxIdle,
// 			Wait:        Wait,
// 			Dial: func() (conn redis.Conn, e error) {
// 				return redis.Dial("tcp", address)
// 			},
// 			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
// 				if _, err := conn.Do("PING"); err != nil {
// 					return fmt.Errorf("redis error: %v", err)
// 				}
// 				return nil
// 			},
// 		},
// 	}
// }
//
// func (pool *RedisPool) GetRedisClient() *Redis {
// 	return &Redis{pool.Get()}
// }
