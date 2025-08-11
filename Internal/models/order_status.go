package models

type OrderStatus string

const (
	Cart       OrderStatus = "cart"
	Pending    OrderStatus = "pending"
	Processing OrderStatus = "processing"
	Completed  OrderStatus = "completed"
	Cancelled  OrderStatus = "cancelled"
)
