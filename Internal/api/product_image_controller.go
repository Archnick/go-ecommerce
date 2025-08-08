package api

import (
	"net/http"
	"strconv"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductImageController struct {
	// Assuming you have a database connection or other dependencies here
	db *gorm.DB
}

// NewProductImageController creates a new instance of ProductImageController.
func NewProductImageController(db *gorm.DB) *ProductImageController {
	return &ProductImageController{db: db}
}

// handleGetProductImages handles the retrieval of product images.
func (c *ProductImageController) handleGetProductImages(ctx *gin.Context) {
	productIDStr := ctx.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var images []models.ProductImage
	result := c.db.Where("product_id = ?", productID).Find(&images)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product images"})
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (c *ProductImageController) handleCreateProductImage(ctx *gin.Context) {
	productIDStr := ctx.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var user models.User
	userID, _ := ctx.Get("userID")
	if result := c.db.First(&user, userID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var product models.Product
	// Check if the product exists
	if result := c.db.First(&product, productID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if models.Role(user.Role) != models.AdminRole && user.Shop.ID != product.ShopID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add images to this product"})
		return
	}

	// Bind the JSON payload to the ProductImagePayload struct
	var payload models.ProductImagePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	image := models.ProductImage{
		ImageURL:  payload.ImageURL,
		AltText:   payload.AltText,
		ProductID: uint(productID),
	}

	result := c.db.Create(&image)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product image"})
		return
	}

	ctx.JSON(http.StatusCreated, image)
}

func (c *ProductImageController) handleDeleteProductImage(ctx *gin.Context) {
	imageIDStr := ctx.Param("image_id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil || imageID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	var productImage models.ProductImage
	if result := c.db.First(&productImage, imageID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product image not found"})
		return
	}

	var user models.User
	userID, _ := ctx.Get("userID")
	if result := c.db.First(&user, userID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var product models.Product
	// Check if the product exists
	if result := c.db.First(&product, productImage.ProductID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if models.Role(user.Role) != models.AdminRole && user.Shop.ID != product.ShopID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add images to this product"})
		return
	}

	result := c.db.Delete(&models.ProductImage{}, imageID)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product image"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product image not found"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
