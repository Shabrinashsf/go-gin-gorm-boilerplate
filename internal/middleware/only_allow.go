package middleware

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/constants"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func OnlyAllow(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.GetString(constants.CTX_KEY_ROLE_NAME)

		for _, role := range roles {
			if userRole == role {
				ctx.Next()
				return
			}
		}

		res := response.NewFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, dto.ErrRoleNotAllowed, nil)
		res.SendWithAbort(ctx)
	}
}
