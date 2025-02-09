package gin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	orderService *services.OrderService
}

func (o *OrdersHandler) Create(c *gin.Context) {
	var body domain.NewOrder

	if err := c.BindJSON(&body); err != nil {
		return
	}

	order, err := o.orderService.CreateOrder(body)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (o *OrdersHandler) List(c *gin.Context) {
	var queryParams struct {
		Active bool `form:"active"`
	}

	if err := c.BindQuery(&queryParams); err != nil {
		return
	}

	var filters []domain.OrderFilterFn

	if queryParams.Active {
		filters = append(filters, domain.FilterActive)
	}

	orders, err := o.orderService.FindMany(filters...)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (o *OrdersHandler) Find(c *gin.Context) {
	id := c.Param("id")

	orderID, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := o.orderService.FindByID(uint(orderID))

	if err != nil {
		abortWithOrderError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

func (o *OrdersHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	status := domain.OrderStatus(c.Param("status"))

	orderID, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := o.orderService.UpdateStatus(uint(orderID), status)

	if err != nil {
		abortWithOrderError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

func abortWithOrderError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, domain.ErrOrderNotFound) {
		status = http.StatusNotFound
	}

	c.AbortWithStatus(status)
}
