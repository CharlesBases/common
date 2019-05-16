package db

import (
	"testing"

	"github.com/jinzhu/gorm"
)

type User struct {
}

func TestTransaction(t *testing.T) {
	user := User{}
	Post(&user)
}

func Post(user *User) error {
	return Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}
