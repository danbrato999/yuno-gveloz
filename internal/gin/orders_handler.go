package gin

import (
	"net/http"
	"strconv"

	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	orderService *domain.OrderService
}

func (o *OrdersHandler) Create(context *gin.Context) {
	var body domain.NewOrder

	if err := context.BindJSON(&body); err != nil {
		return
	}

	order, err := o.orderService.CreateOrder(body)

	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusCreated, order)
}

func (o *OrdersHandler) List(context *gin.Context) {
	orders, err := o.orderService.FindMany()

	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (o *OrdersHandler) Find(context *gin.Context) {
	id := context.Param("id")

	orderID, err := strconv.Atoi(id)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := o.orderService.FindByID(uint(orderID))

	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, order)
}

func (o *OrdersHandler) Update(context *gin.Context) {
	id := context.Param("id")

	orderID, err := strconv.Atoi(id)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var update domain.OrderUpdate

	if err = context.BindJSON(&update); err != nil {
		return
	}

	order, err := o.orderService.Update(uint(orderID), update)

	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, order)
}

func (o *OrdersHandler) Delete(context *gin.Context) {
	id := context.Param("id")

	orderID, err := strconv.Atoi(id)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := o.orderService.Cancel(uint(orderID))

	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, order)
}
