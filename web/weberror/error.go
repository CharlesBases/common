package weberror

import (
	"fmt"
	"strings"
)

const (
	Code = 100000
)

var err = []string{
	Code: "错误说明",
}

type WebError map[string]interface{}

func Success(msgs ...interface{}) WebError {
	if msgs != nil {
		return WebError{"code": 0, "data": msgs}
	}
	return WebError{"code": 0, "data": "success"}
}

func NewWebrror(code int, msgs ...string) WebError {
	if msgs != nil {
		eStr := fmt.Sprintf("%s -- %s", err[code], strings.Join(msgs, ","))
		return WebError{"code": code, "msg": eStr}
	}
	return WebError{"code": code, "msg": err[code]}
}
