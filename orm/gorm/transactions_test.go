package gorm

import "github.com/jinzhu/gorm"

func Transactions() error {
	return Transaction(
		func(tx *gorm.DB) error {
			return tx.Error
		},
		func(tx *gorm.DB) error {
			return tx.Error
		},
	)
}
