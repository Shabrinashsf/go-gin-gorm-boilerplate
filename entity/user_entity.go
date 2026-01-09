package entity

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`

	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	Instansi   string   `json:"instansi"`
	NoTelp     string   `json:"no_telp"`
	Role       UserRole `json:"role" gorm:"default:user"`
	IsVerified bool     `json:"is_verified"`

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var err error
	// u.ID = uuid.New()
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
