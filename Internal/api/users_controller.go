package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin" // Import Gin
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// UsersController holds the dependencies for user-related handlers.
type UsersController struct {
	db *gorm.DB
}

type PublicUser struct {
	ID              uint      `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	ProfileImageURL string    `json:"photo"`
	Role            string    `json:"role"`
	CreatedAt       time.Time `json:"createdAt"`
}

// NewUsersController creates a new instance of the UsersController.
func NewUsersController(db *gorm.DB) *UsersController {
	return &UsersController{db: db}
}

// handleGetUsers now takes a *gin.Context.
func (c *UsersController) handleGetUsers(ctx *gin.Context) {
	var users []models.User

	result := c.db.Find(&users)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	publicUsers := make([]PublicUser, len(users))

	for i, user := range users {
		publicUsers[i] = PublicUser{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			ProfileImageURL: user.ProfileImageURL,
			Role:            user.Role,
			CreatedAt:       user.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, publicUsers)
}

func (c *UsersController) handleGetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var user models.User
	result := c.db.First(&user, id)

	if result.Error != nil {
		errorMessage := fmt.Sprintf("User with the id: %s not found", id)
		ctx.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	publicUser := PublicUser{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		ProfileImageURL: user.ProfileImageURL,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, publicUser)
}

func (c *UsersController) handleUpdateUser(ctx *gin.Context) {
	// 1. Authorization Check (same as handleGetUser)
	requestingUserID, _ := ctx.Get("userID")
	requestingUserRole, _ := ctx.Get("role")

	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.Atoi(targetUserIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if models.Role(requestingUserRole.(string)) != models.AdminRole && requestingUserID.(uint) != uint(targetUserID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this user"})
		return
	}

	// 2. Bind and Validate the Payload
	var payload models.UpdateUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs.Translate(trans)})
		return
	}

	if payload.Role != "" {
		requestingUserRole, _ = ctx.Get("role")
		if models.Role(requestingUserRole.(string)) != models.AdminRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to modify users roles"})
			return
		}

	}

	// 3. Fetch the User from the Database
	var user models.User
	if result := c.db.First(&user, targetUserID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 4. Update the User Record
	// GORM's Updates method will only update non-zero fields,
	// which works perfectly with our optional payload.
	if err := c.db.Model(&user).Updates(payload).Error; err != nil {
		slog.Error("failed to update user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *UsersController) handleDeleteUser(ctx *gin.Context) {
	// 1. Authorization Check (same as handleGetUser and handleUpdateUser)
	requestingUserID, _ := ctx.Get("userID")
	requestingUserRole, _ := ctx.Get("role")

	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.Atoi(targetUserIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if models.Role(requestingUserRole.(string)) != models.AdminRole && requestingUserID.(uint) != uint(targetUserID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this user"})
		return
	}

	// GORM will perform a soft delete
	if result := c.db.Delete(&models.User{}, targetUserID); result.Error != nil {
		slog.Error("failed to delete user", "error", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
