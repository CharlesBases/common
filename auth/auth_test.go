package auth

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	SecretKey = "shdkkj&(hkdksaYKBKDJah890uiojoiu0KNKSAdhka892hkj!@kndsajhd"
)

func TestAuth(t *testing.T) {
	user := User{
		UserId:    1,
		Timestamp: time.Now().UnixNano(),
	}

	token, err := GenToken(SecretKey, time.Hour*1, &user)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println("Token： " + token)

	value, err := ParToken(SecretKey, token)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println("value:  ", *value)
}
