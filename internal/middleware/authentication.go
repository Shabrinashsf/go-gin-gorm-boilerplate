package middleware

import (
	"strings"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/constants"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.ErrTokenNotFound, nil).SendWithAbort(ctx)
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.ErrTokenNotValid, nil).SendWithAbort(ctx)
			return
		}
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.ErrTokenNotValid, nil).SendWithAbort(ctx)
			return
		}
		if !token.Valid {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.ErrDeniedAccess, nil).SendWithAbort(ctx)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.ErrTokenNotValid, nil).SendWithAbort(ctx)
			return
		}
		roleClaim, _ := claims[constants.CTX_KEY_ROLE_NAME]
		role, _ := roleClaim.(string)
		if role == "" {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.ErrDeniedAccess, nil).SendWithAbort(ctx)
			return
		}
		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err, nil).SendWithAbort(ctx)
			return
		}
		ctx.Set("token", authHeader)
		ctx.Set("user_id", userId)
		ctx.Set(constants.CTX_KEY_ROLE_NAME, role)
		ctx.Next()
	}
}
