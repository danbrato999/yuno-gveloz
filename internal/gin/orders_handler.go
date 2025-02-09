package gin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	orderService *domain.OrderService
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
	orders, err := o.orderService.FindMany()

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

func (o *OrdersHandler) Update(c *gin.Context) {
	id := c.Param("id")

	orderID, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var update domain.OrderUpdate

	if err = c.BindJSON(&update); err != nil {
		return
	}

	order, err := o.orderService.Update(uint(orderID), update)

	if err != nil {
		abortWithOrderError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

func (o *OrdersHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	orderID, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, err := o.orderService.Cancel(uint(orderID))

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
