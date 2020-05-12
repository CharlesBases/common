package auth

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	user := User{
		UserID:    1,
		Timestamp: time.Now().UnixNano(),
	}

	token, err := GenToken(&user)
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
