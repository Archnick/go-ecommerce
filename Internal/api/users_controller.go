package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/Archnick/go-ecommerce/Internal/services"
	"github.com/gin-gonic/gin" // Import Gin
	"github.com/go-playground/validator/v10"
)

// UsersController holds the dependencies for user-related handlers.
type UsersController struct {
	service *services.UserService
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
func NewUsersController(service *services.UserService) *UsersController {
	return &UsersController{service: service}
}

// handleGetUsers now takes a *gin.Context.
func (c *UsersController) handleGetUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()

	if err != nil {
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		errorMessage := fmt.Sprintf("User with the id: %v not found", id)
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
	// Authorization Check
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

	// Bind and Validate the Payload
	var payload models.UpdateUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs.Translate(trans)})
		return
	}

	if payload.Role != "" {
		if models.Role(requestingUserRole.(string)) != models.AdminRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to modify users roles"})
			return
		}
	}

	// Fetch the User from the Database
	user, err := c.service.GetUserByID(uint(targetUserID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 4. Update the User Record
	err = c.service.UpdateUser(user, payload)
	if err != nil {
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

	err = c.service.DeleteUser(uint(targetUserID))
	if err != nil {
		slog.Error("failed to delete user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
