package middleware

import (
	"fmt"
	"net/http"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/constants"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/response"
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

		err := fmt.Sprintf(dto.ErrRoleNotAllowed.Error(), userRole)
		response := response.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, err, nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}
}
