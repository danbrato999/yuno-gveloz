package domain

import "fmt"

var ErrOrderNotFound = fmt.Errorf("Order not found")
var ErrInvalidStatusUpdate = fmt.Errorf("Order cannot be updated to provided status")
var ErrCompleteOrderUpdate = fmt.Errorf("Completed order cannot be updated")
