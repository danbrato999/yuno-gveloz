package domain

import "time"

type OrderStatusHistory struct {
	Status    OrderStatus `json:"status"`
	Timestamp *time.Time  `json:"timestamp"`
}
