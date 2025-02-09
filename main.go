package main

import (
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gin"
	dbAdapter "github.com/danbrato999/yuno-gveloz/internal/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	err = dbAdapter.Migrate(db)
	if err != nil {
		panic("failed to migrate database")
	}

	orderStore := dbAdapter.NewOrderStore(db)
	orderStatusStore := dbAdapter.NewOrderStatusStore(db)
	orderService := services.NewOrderService(orderStore, orderStatusStore)

	server := gin.GetServer(orderService)
	server.Run(":9001")
}
