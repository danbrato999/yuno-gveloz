package main

import (
	dbAdapter "github.com/danbrato999/yuno-gveloz/internal/adapters/gorm"
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"github.com/danbrato999/yuno-gveloz/internal/gin"
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
	orderService := domain.NewOrderService(orderStore)

	server := gin.GetServer(orderService)
	server.Run(":9001")
}
