package routes

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/middleware"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("api/auth")
	{
		routes.POST("", userController.RegisterUser)
		routes.POST("/login", userController.Login)
		routes.POST("/send-verification-email", userController.SendVerificationEmail)
		routes.GET("/verify-email", userController.VerifyEmail)
		routes.POST("/forgot-password", userController.ForgotPassword)
		routes.POST("/reset-password", userController.ResetPassword)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.MeAuth)
		routes.PATCH("/update", middleware.Authenticate(jwtService), userController.UpdateUser)
	}
}
