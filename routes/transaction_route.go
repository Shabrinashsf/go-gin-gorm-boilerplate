package routes

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/controller"
	"github.com/gin-gonic/gin"
)

func Transaction(route *gin.Engine, transactionController controller.TransactionController) {
	routes := route.Group("/transaction")
	{
		routes.POST("/webhook/tripay", transactionController.TripayWebhook)
	}
}
