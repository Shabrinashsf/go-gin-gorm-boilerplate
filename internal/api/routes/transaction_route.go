package routes

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/controller"
	"github.com/gin-gonic/gin"
)

func Transaction(route *gin.Engine, transactionController controller.TransactionController) {
	routes := route.Group("/api/transaction")
	{
		routes.POST("/webhook/tripay", transactionController.TripayWebhook)
	}
}
