package domain

import "fmt"

var ErrOrderNotFound = fmt.Errorf("Order not found")
var ErrInvalidOrderUpdate = fmt.Errorf("Order updated is incorrect")
var ErrCompleteOrderUpdate = fmt.Errorf("Completed order cannot be updated")
var ErrIncorrectOrderQueueing = fmt.Errorf("Order queue operation is not valid")
