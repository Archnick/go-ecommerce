package models

type OrderPayload struct {
	TotalAmount     float64            `json:"total_amount"`
	ShippingAddress string             `json:"shipping_address"`
	UserID          uint               `json:"user_id"`
	OrderItems      []OrderItemPayload `json:"order_items"`
}

type OrderItemPayload struct {
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	ProductID uint    `json:"product_id"`
}

type UpdateOrderPayload struct {
	Status          string `json:"status"`
	ShippingAddress string `json:"shipping_address"`
}

type UpdateOrderItemPayload struct {
	Quantity int `json:"quantity"`
}
