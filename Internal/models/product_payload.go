package models

type ProductPayload struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	ShopID      uint    `json:"shop_id" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateProductPayload struct {
	Name        string  `json:"name" binding:"omitempty,min=2"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	ShopID      uint    `json:"shop_id" binding:"omitempty"`
	CategoryID  uint    `json:"category_id" binding:"omitempty"`
}

type ProductImagePayload struct {
	ProductID uint   `json:"product_id" binding:"required"`
	AltText   string `json:"alt_text" binding:"omitempty"`
	ImageURL  string `json:"image_url" binding:"required,url"`
}
