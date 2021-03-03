package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	privateKey = "shdkkj&(hkdksaYKBKDJah890uiojoiu0KNKSAdhka892hkj!@kndsajhd"
	ttl        = 4 * time.Hour
)

func TestAuth(t *testing.T) {
	InitAuth(WithPrivateKey(privateKey), WithTTL(ttl))

	tokenString, err := auth.GenToken("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("token:\n", tokenString)

	account, err := auth.ParToken(tokenString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("account:\n", func() string {
		data, _ := json.MarshalIndent(account, "", "  ")
		return string(data)
	}())
}
