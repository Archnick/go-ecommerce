package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin" // Import Gin
	"gorm.io/gorm"
)

// UsersController holds the dependencies for user-related handlers.
type UsersController struct {
	db *gorm.DB
}

type PublicUser struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
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
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, publicUsers)
}

func (c *UsersController) handleGetUser(ctx *gin.Context) {
	// get users id
	requestingUserID, _ := ctx.Get("userID")
	requestingUserRole, _ := ctx.Get("role")

	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.Atoi(targetUserIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if models.Role(requestingUserRole.(string)) != models.AdminRole && requestingUserID.(uint) != uint(targetUserID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this user"})
		return
	}

	id := ctx.Param("id")

	var user models.User
	result := c.db.First(&user, id)

	if result.Error != nil {
		errorMessage := fmt.Sprintf("User with the id: %s not found", id)
		ctx.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	publicUser := PublicUser{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, publicUser)
}
