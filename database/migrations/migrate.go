package migrations

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Transaction{},
	); err != nil {
		return err
	}

	return nil
}
