package routes

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/middleware"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("api/v1/auth")
	{
		routes.POST("", userController.RegiterUser)
		routes.POST("/login", userController.Login)
		routes.POST("/send-verification-email", userController.SendVerificationEmail)
		routes.POST("/verify-email", userController.VerifyEmail)
		routes.POST("/forgot-password", userController.ForgotPassword)
		routes.POST("/reset-password", userController.ResetPassword)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.MeAuth)
		routes.PUT("/update", middleware.Authenticate(jwtService), userController.UpdateUser)
	}
}
