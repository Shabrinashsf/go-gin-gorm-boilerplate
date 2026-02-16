package seeders

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/database/seeders/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}

	return nil
}
