package auth

import (
	"testing"

	"github.com/CharlesBases/common/log"
)

func TestToken(t *testing.T) {
	defer log.Flush()

	// 生成token
	infor := Infor{User: "JWT"}
	token, err := GenTempToken(&infor)
	if err != nil {
		log.Error(err)
	}
	log.Error(token)

	// 解析token
	value, err := ParseTempToken(token)
	if err != nil {
		log.Error(err)
	}
	log.Error(value)
}
