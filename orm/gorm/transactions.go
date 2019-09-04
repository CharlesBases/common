package gorm

import "github.com/jinzhu/gorm"

func Transaction(fs ...func(tx *gorm.DB) error) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, f := range fs {
		if err := f(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
