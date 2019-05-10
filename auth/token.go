package auth

/*
 temp token
*/

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var securityKey = "shdkkj&(hkdksaYKBKDJah890uiojoiu0KNKSAdhka892hkj!@kndsajhd"

type Infor struct {
	User string
}

type infor struct {
	Infor Infor
	jwt.StandardClaims
}

func GenTempToken(infor *Infor) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(map[string]interface{}{
		// "iat":   time.Now().Unix(),
		// "exp":   time.Now().Add(time.Minute * 5).Unix(),
		"infor": infor,
	})).SignedString([]byte(securityKey))
}

func ParseTempToken(tokenString string) (*Infor, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &infor{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(securityKey), nil
	})
	if claims, ok := token.Claims.(*infor); ok && token.Valid {
		return &claims.Infor, nil
	}
	return nil, errors.New("token is error")
}
