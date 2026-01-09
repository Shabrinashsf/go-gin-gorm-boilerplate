package middleware

import (
	"net/http"
	"strings"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/constants"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !token.Valid {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_DENIED_ACCESS, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		roleClaim, _ := claims[constants.CTX_KEY_ROLE_NAME]
		role, _ := roleClaim.(string)
		if role == "" {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_DENIED_ACCESS, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			response := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		ctx.Set("token", authHeader)
		ctx.Set("user_id", userId)
		ctx.Set(constants.CTX_KEY_ROLE_NAME, role)
		ctx.Next()
	}
}
