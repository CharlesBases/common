package algo

import (
	"fmt"
	"strings"
	"time"
)

/*
	0-00000000000000000000000000000000000000000-0000000000-000000000000
	- ----------------------------------------- ---------- ------------
	1                     2                          3           4

	1: 固定值, TraceID 为正数, 固定为 0
	2: 时间戳
	3: 机器 ID
	4: 序号
*/
var (
	sequence      int64
	lasttimestamp int64
)

func init() {
	lasttimestamp = timestamp()
}

func GetTraceID(machineid ...int64) string {
	machineID := func() int64 {
		if len(machineid) != 0 {
			// 机器 ID, 最多可部署 1024 台机器. [0, 1024]
			if machineid[0] > -1 && machineid[0] < 1024 {
				return machineid[0]
			}
		}
		return 0
	}()

	currtimestamp := timestamp()

	if currtimestamp == lasttimestamp {
		sequence++
		// 每毫秒最多产生 4096 个不同的 TraceID. [0, 4096]
		if sequence > 4095 {
			time.Sleep(time.Millisecond)

			currtimestamp = timestamp()
			lasttimestamp = currtimestamp
			sequence = 0
		}
	}

	if currtimestamp > lasttimestamp {
		sequence = 0
		lasttimestamp = currtimestamp
	}

	if currtimestamp < lasttimestamp {
		lasttimestamp = currtimestamp
		return ""
	}

	tarceID := currtimestamp << 22
	return DecHex(tarceID | machineID | sequence)
}

func timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// 十进制转换成二进制
func DecBin(dec int64) string {
	biner := strings.Builder{}
	if dec != 0 {
		biner.WriteString(fmt.Sprintf("%s%d", DecBin(dec/2), dec%2))
	}
	return strings.TrimSpace(biner.String())
}

// 十进制转换成十六进制 Hexadecimal
var hexadecimal = map[int64]string{
	0:  "0",
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "a",
	11: "b",
	12: "c",
	13: "d",
	14: "e",
	15: "f",
}

func DecHex(dec int64) string {
	hexer := strings.Builder{}
	if dec != 0 {
		hexer.WriteString(fmt.Sprintf("%s%s", DecHex(dec/16), hexadecimal[dec%16]))
	}
	return hexer.String()
}
