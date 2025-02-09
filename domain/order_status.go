package domain

type OrderStatus string

const OrderStatusPending OrderStatus = "pending"
const OrderStatusActive OrderStatus = "active"
const OrderStatusReady OrderStatus = "ready"
const OrderStatusDone OrderStatus = "done"
const OrderStatusCancelled OrderStatus = "cancelled"

var statusWeights = map[OrderStatus]uint{
	OrderStatusPending:   10,
	OrderStatusActive:    20,
	OrderStatusReady:     30,
	OrderStatusDone:      40,
	OrderStatusCancelled: 40,
}
