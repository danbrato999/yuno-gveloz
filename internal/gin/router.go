package gin

import (
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetServer(orderService *domain.OrderService) *gin.Engine {
	ordersHandler := &OrdersHandler{
		orderService: orderService,
	}

	router := gin.Default()

	api := router.Group("/api/v1")

	orders := api.Group("/orders")
	orders.GET("", ordersHandler.List)
	orders.POST("", ordersHandler.Create)

	orders.GET("/:id", ordersHandler.Find)
	orders.PUT("/:id", ordersHandler.Update)
	orders.DELETE("/:id", ordersHandler.Delete)

	return router
}
