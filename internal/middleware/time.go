package middleware

import (
	"net/http"
	"time"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	myerror "github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/error"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ApiLockRange struct {
	Start   string
	End     string
	Message string
}

type (
	LockApiMiddleware struct {
		IsLocked bool
		location *time.Location
	}

	LockOption func(m *LockApiMiddleware)
)

func LockAPI(msg string, opts ...LockOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		location, _ := time.LoadLocation("Asia/Jakarta")
		lockApiMiddleware := LockApiMiddleware{
			IsLocked: false,
			location: location,
		}

		for _, opt := range opts {
			opt(&lockApiMiddleware)
		}

		if lockApiMiddleware.IsLocked {
			response.NewFailed(dto.MESSAGE_API_IS_LOCKED, myerror.New(msg, http.StatusForbidden)).
				SendWithAbort(ctx)
			return
		}

		ctx.Next()
	}
}

func LockAPIByKey(lockRanges map[string]ApiLockRange, key string) gin.HandlerFunc {
	lockRange, ok := lockRanges[key]
	if !ok || lockRange.Start == "" || lockRange.End == "" {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	msg := lockRange.Message
	if msg == "" {
		msg = "API is locked"
	}

	return LockAPI(msg, NotInRange(lockRange.Start, lockRange.End))
}

func NotBefore(t string) LockOption {
	return func(ml *LockApiMiddleware) {
		parsedTime, err := time.ParseInLocation("02-01-2006 15:04:05", t, ml.location)
		if err != nil {
			return
		}

		now := time.Now().In(ml.location)
		if now.Before(parsedTime) {
			ml.IsLocked = true
		}
	}
}

func NotAfter(t string) LockOption {
	return func(ml *LockApiMiddleware) {
		parsedTime, err := time.ParseInLocation("02-01-2006 15:04:05", t, ml.location)
		if err != nil {
			return
		}
		now := time.Now().In(ml.location)
		if now.After(parsedTime) {
			ml.IsLocked = true
		}
	}
}

func NotInRange(start, end string) LockOption {
	return func(ml *LockApiMiddleware) {
		startTime, err1 := time.ParseInLocation("02-01-2006 15:04:05", start, ml.location)
		endTime, err2 := time.ParseInLocation("02-01-2006 15:04:05", end, ml.location)

		if err1 != nil || err2 != nil {
			return
		}

		now := time.Now().In(ml.location)
		if now.Before(startTime) || now.After(endTime) {
			ml.IsLocked = true
		}
	}
}
