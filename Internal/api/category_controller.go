package api

import (
	"net/http"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	db *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{db: db}
}

// handleGetCategories retrieves all categories from the database.
func (c *CategoryController) handleGetCategories(ctx *gin.Context) {
	var categories []models.Category
	result := c.db.Find(&categories)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	publicCategories := make([]models.CategoryPayload, len(categories))
	for i, category := range categories {
		publicCategories[i] = models.CategoryPayload{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
	}
	ctx.JSON(http.StatusOK, publicCategories)
}

// handleCreateCategory creates a new category in the database.
func (c *CategoryController) handleCreateCategory(ctx *gin.Context) {
	var payload models.CategoryPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	category := models.Category{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if result := c.db.Create(&category); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	response := models.CategoryPayload{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category": response})
}

// handleUpdateCategory updates an existing category in the database.
func (c *CategoryController) handleUpdateCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	var payload models.UpdateCategoryPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var category models.Category
	if err := c.db.First(&category, categoryID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if payload.Name != "" {
		category.Name = payload.Name
	}
	if payload.Description != "" {
		category.Description = payload.Description
	}

	if result := c.db.Save(&category); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	response := models.CategoryPayload{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "category": response})
}

// handleDeleteCategory deletes a category from the database.
func (c *CategoryController) handleDeleteCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if result := c.db.Delete(&models.Category{}, categoryID); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
