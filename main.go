package main

import (
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gin"
	dbAdapter "github.com/danbrato999/yuno-gveloz/internal/gorm"
)

const dbName = "main"

func main() {
	db, err := dbAdapter.GetDBConnection(dbName)
	if err != nil {
		panic(err.Error())
	}

	orderStore := dbAdapter.NewOrderStore(db)
	orderStatusStore := dbAdapter.NewOrderStatusStore(db)
	orderService := services.NewOrderService(orderStore, orderStatusStore)

	server := gin.GetServer(orderService)
	server.Run(":9001")
}
