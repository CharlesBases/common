package utils

import "time"

// 毫秒时间戳
func MilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// 零点毫秒时间戳
func ZeroMilliSecond(years, months, days int) int64 {
	zerotime, _ := time.ParseInLocation("2006-01-02", time.Now().AddDate(years, months, days).Format("2006-01-02"), time.Local)
	return zerotime.UnixNano() / 1e6
}
