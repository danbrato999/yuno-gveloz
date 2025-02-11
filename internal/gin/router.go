package gin

import (
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/gin-gonic/gin"
)

func addOrderRoutes(ordersHandler *OrdersHandler, api *gin.RouterGroup) {
	orders := api.Group("/orders")
	orders.GET("", ordersHandler.List)
	orders.POST("", ordersHandler.Create)

	orders.GET("/:id", ordersHandler.Find)
	orders.PUT("/:id/status/:status", ordersHandler.UpdateStatus)
}

func GetServer(orderService services.OrderService) *gin.Engine {
	ordersHandler := &OrdersHandler{
		orderService: orderService,
	}

	router := gin.Default()

	api := router.Group("/api/v1")
	addOrderRoutes(ordersHandler, api)
	return router
}
