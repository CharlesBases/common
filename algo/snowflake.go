package algo

import (
	"fmt"
	"strings"

	"github.com/sony/sonyflake"

	"github.com/CharlesBases/common/log"
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
	sf *sonyflake.Sonyflake
)

func init() {
	var st sonyflake.Settings
	// st.StartTime = time.Now()

	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Fatal("sonyflake init failed")
	}
}

// NextID 返回十进制 id
func NextID() uint64 {
	id, _ := sf.NextID()
	return id
}

// DecBin 十进制转换成二进制
func DecBin(dec uint64) string {
	biner := strings.Builder{}
	if dec != 0 {
		biner.WriteString(fmt.Sprintf("%s%d", DecBin(dec/2), dec%2))
	}
	return strings.TrimSpace(biner.String())
}

// 十进制转换成十六进制 Hexadecimal
var hexadecimal = map[uint64]string{
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

// DecHex 十进制转十六进制
func DecHex(dec uint64) string {
	hexer := strings.Builder{}
	if dec != 0 {
		hexer.WriteString(fmt.Sprintf("%s%s", DecHex(dec/16), hexadecimal[dec%16]))
	}
	return hexer.String()
}
