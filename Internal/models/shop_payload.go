package models

type ShopPayload struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" binding:"required,min=2"`
	Description  string `json:"description" binding:"omitempty"`
	ShopImageUrl string `json:"shop_image_url" binding:"omitempty,url"`
}

type UpdateShopPayload struct {
	Name         string `json:"name" binding:"omitempty,min=2"`
	Description  string `json:"description" binding:"omitempty"`
	ShopImageUrl string `json:"shop_image_url" binding:"omitempty,url"`
}
