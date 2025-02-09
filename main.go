package main

import (
	"github.com/danbrato999/yuno-gveloz/internal/adapters"
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"github.com/danbrato999/yuno-gveloz/internal/gin"
)

func main() {
	orderStore := adapters.NewOrderStore()
	orderService := domain.NewOrderService(orderStore)

	server := gin.GetServer(orderService)
	server.Run(":9001")
}
