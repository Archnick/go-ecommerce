package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{db: db}
}

type PublicProduct struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Stock       int             `json:"stock"`
	ShopID      uint            `json:"shop_id"`
	CategoryID  uint            `json:"category_id"`
	Images      []string        `json:"images"`
	Reviews     []models.Review `json:"reviews,omitempty"`
}

// handleGetProducts now takes a *gin.Context.
func (c *ProductController) handleGetProducts(ctx *gin.Context) {
	var products []models.Product

	result := c.db.Preload("Images").Find(&products)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	publicProducts := make([]PublicProduct, len(products))
	for i, product := range products {
		publicProducts[i] = PublicProduct{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			ShopID:      product.ShopID,
			CategoryID:  product.CategoryID,
			Images:      make([]string, len(product.Images)),
		}
		for j, image := range product.Images {
			publicProducts[i].Images[j] = image.ImageURL
		}
	}
	ctx.JSON(http.StatusOK, publicProducts)
}

func (c *ProductController) handleGetProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	result := c.db.Preload("Images").Preload("Reviews").First(&product, productID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		}
		return
	}

	publicProduct := PublicProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ShopID:      product.ShopID,
		CategoryID:  product.CategoryID,
		Images:      make([]string, len(product.Images)),
		Reviews:     product.Reviews,
	}
	for j, image := range product.Images {
		publicProduct.Images[j] = image.ImageURL
	}
	ctx.JSON(http.StatusOK, publicProduct)
}

func (c *ProductController) handleCreateProduct(ctx *gin.Context) {
	var payload models.ProductPayload
	// Bind the JSON payload to the ProductPayload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs.Translate(trans)})
		return
	}
	var user models.User
	userID, _ := ctx.Get("userID")
	if result := c.db.First(&user, userID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	shopId := user.Shop.ID

	// only admins can create products for other shops
	if payload.ShopID != 0 {
		if models.Role(user.Role) != models.AdminRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to decide the shop ID"})
			return
		}
		shopId = payload.ShopID
	}

	product := models.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		ShopID:      shopId,
		CategoryID:  payload.CategoryID,
	}

	result := c.db.Create(&product)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product_id": product.ID})
}

func (c *ProductController) handleUpdateProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	var user models.User
	userID, _ := ctx.Get("userID")
	if result := c.db.First(&user, userID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch the Product from the Database
	var product models.Product
	if result := c.db.First(&product, productID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if models.Role(user.Role) != models.AdminRole && user.Shop.ID != product.ShopID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this product"})
		return
	}

	// 2. Bind and Validate the Payload
	var payload models.UpdateProductPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs.Translate(trans)})
		return
	}

	if payload.ShopID != 0 && models.Role(user.Role) != models.AdminRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to modify products shop ID"})
		return
	}

	// Update the Product Record
	// GORM's Updates method will only update non-zero fields,
	// which works perfectly with our optional payload.
	if err := c.db.Model(&product).Updates(payload).Error; err != nil {
		slog.Error("failed to update product", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *ProductController) handleDeleteProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var user models.User
	userID, _ := ctx.Get("userID")
	if result := c.db.First(&user, userID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Check if the user is authorized to delete the product
	var count int64
	c.db.Model(&models.Product{}).Where("id = ?", productID).Where("shop_id = ?", user.Shop.ID).Count(&count)

	if count == 0 && models.Role(user.Role) != models.AdminRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this product"})
		return
	}

	result := c.db.Delete(&models.Product{}, productID)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
