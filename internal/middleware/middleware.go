package middleware

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/service"
	"gorm.io/gorm"
)

type ApiLockRange struct {
	Start   string
	End     string
	Message string
}

type Middleware struct {
	db         *gorm.DB
	jwtService service.JWTService
	lockRanges map[string]ApiLockRange
}

func New(db *gorm.DB, jwtService service.JWTService, lockRanges map[string]ApiLockRange) Middleware {
	return Middleware{db, jwtService, lockRanges}
}
