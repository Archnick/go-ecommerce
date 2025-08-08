package api

import (
	"net/http"
	"strconv"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopController struct {
	db *gorm.DB
}

func NewShopController(db *gorm.DB) *ShopController {
	return &ShopController{db: db}
}

// handleGetShops retrieves all shops from the database.
func (c *ShopController) handleGetShops(ctx *gin.Context) {
	var shops []models.Shop
	result := c.db.Find(&shops)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shops"})
		return
	}
	publicShops := make([]models.ShopPayload, len(shops))
	for i, shop := range shops {
		publicShops[i] = models.ShopPayload{
			ID:           shop.ID,
			Name:         shop.Name,
			Description:  shop.Description,
			ShopImageUrl: shop.ShopImageUrl,
		}
	}
	ctx.JSON(http.StatusOK, publicShops)
}

func (c *ShopController) handleGetShop(ctx *gin.Context) {
	shopIDStr := ctx.Param("id")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil || shopID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var shop models.Shop
	if result := c.db.First(&shop, shopID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	response := models.ShopPayload{
		ID:           shop.ID,
		Name:         shop.Name,
		Description:  shop.Description,
		ShopImageUrl: shop.ShopImageUrl,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *ShopController) handleCreateShop(ctx *gin.Context) {
	var payload models.ShopPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	shop := models.Shop{
		Name:         payload.Name,
		Description:  payload.Description,
		ShopImageUrl: payload.ShopImageUrl,
	}

	if result := c.db.Create(&shop); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shop"})
		return
	}

	response := models.ShopPayload{
		ID:           shop.ID,
		Name:         shop.Name,
		Description:  shop.Description,
		ShopImageUrl: shop.ShopImageUrl,
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *ShopController) handleUpdateShop(ctx *gin.Context) {
	shopIDStr := ctx.Param("id")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil || shopID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var payload models.UpdateShopPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var shop models.Shop
	if result := c.db.First(&shop, shopID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	if payload.Name != "" {
		shop.Name = payload.Name
	}
	if payload.Description != "" {
		shop.Description = payload.Description
	}
	if payload.ShopImageUrl != "" {
		shop.ShopImageUrl = payload.ShopImageUrl
	}

	if result := c.db.Save(&shop); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shop"})
		return
	}

	response := models.ShopPayload{
		ID:           shop.ID,
		Name:         shop.Name,
		Description:  shop.Description,
		ShopImageUrl: shop.ShopImageUrl,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *ShopController) handleDeleteShop(ctx *gin.Context) {
	shopIDStr := ctx.Param("id")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil || shopID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	if result := c.db.Delete(&models.Shop{}, shopID); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shop"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
