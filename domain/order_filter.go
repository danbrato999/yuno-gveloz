package domain

type OrderFilters struct {
	AnyStatus []OrderStatus
}

type OrderFilterFn = func(filter *OrderFilters)

var FilterActive OrderFilterFn = func(filter *OrderFilters) {
	filter.AnyStatus = []OrderStatus{OrderStatusPending, OrderStatusPreparing, OrderStatusReady}
}
