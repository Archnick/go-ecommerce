package api

import (
	"net/http"
	"strconv"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewController struct {
	db *gorm.DB
}

func NewReviewController(db *gorm.DB) *ReviewController {
	return &ReviewController{db: db}
}

type PublicReview struct {
	ID      uint   `json:"id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (c *ReviewController) handleGetReviewsForProduct(ctx *gin.Context) {
	productIDStr := ctx.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var reviews []models.Review
	result := c.db.Where("product_id = ?", productID).Find(&reviews)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product reviews"})
		return
	}

	var publicReviews []PublicReview
	for _, review := range reviews {
		publicReviews = append(publicReviews, PublicReview{
			ID:      review.ID,
			Rating:  review.Rating,
			Comment: review.Comment,
		})
	}

	ctx.JSON(http.StatusOK, publicReviews)
}

func (c *ReviewController) handleGetReviewsForUser(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var reviews []models.Review
	result := c.db.Where("user_id = ?", userID).Find(&reviews)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user reviews"})
		return
	}

	var publicReviews []PublicReview
	for _, review := range reviews {
		publicReviews = append(publicReviews, PublicReview{
			ID:      review.ID,
			Rating:  review.Rating,
			Comment: review.Comment,
		})
	}

	ctx.JSON(http.StatusOK, publicReviews)
}

func (c *ReviewController) handleCreateReview(ctx *gin.Context) {
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

	if user.Shop.ID == product.ShopID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to review products from your own shop"})
		return
	}

	// Bind the JSON payload to the ReviewPayload struct
	var payload models.ReviewPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	review := models.Review{
		Rating:    payload.Rating,
		Comment:   payload.Comment,
		UserID:    user.ID,
		ProductID: product.ID,
	}

	result := c.db.Create(&review)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Review created successfully", "review_id": review.ID})
}

func (c *ReviewController) handleUpdateReview(ctx *gin.Context) {
	reviewIDStr := ctx.Param("review_id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil || reviewID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var user models.User
	userID, _ := ctx.Get("userID")
	if result := c.db.First(&user, userID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var review models.Review
	if result := c.db.First(&review, reviewID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.UserID != user.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this review"})
		return
	}

	var payload models.UpdateReviewPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}
	if payload.Comment != "" {
		review.Comment = payload.Comment
	}

	result := c.db.Save(&review)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review updated successfully", "review_id": review.ID})
}
