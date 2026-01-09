package database

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/database/seed"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seed.ListUserSeeder(db); err != nil {
		return err
	}

	return nil
}
