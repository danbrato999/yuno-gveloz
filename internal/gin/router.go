package gin

import (
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/gin-gonic/gin"
)

func addOrderRoutes(ordersHandler *OrdersHandler, api *gin.RouterGroup) {
	orders := api.Group("/orders")
	orders.GET("", ordersHandler.List)
	orders.POST("", ordersHandler.Create)

	order := orders.Group("/:id")
	order.GET("", ordersHandler.Find)
	order.PUT("", ordersHandler.UpdateContent)
	order.PUT("/status/:status", ordersHandler.UpdateStatus)
	order.PUT("/prioritize", ordersHandler.Prioritize)
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
