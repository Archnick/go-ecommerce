package api

import (
	"net/http"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin" // Import Gin
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UsersController holds the dependencies for user-related handlers.
type UsersController struct {
	db *gorm.DB
}

// NewUsersController creates a new instance of the UsersController.
func NewUsersController(db *gorm.DB) *UsersController {
	return &UsersController{db: db}
}

// handleRegisterUser now takes a *gin.Context.
func (c *UsersController) handleRegisterUser(ctx *gin.Context) {
	var payload models.UserPayload
	// Use ctx.BindJSON() to parse the request body.
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{Email: payload.Email, Password: string(hashedPassword)}
	result := c.db.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Use ctx.JSON() to send a response. gin.H is a shortcut for map[string]interface{}.
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// handleGetUsers now takes a *gin.Context.
func (c *UsersController) handleGetUsers(ctx *gin.Context) {
	var users []models.User

	result := c.db.Find(&users)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	type PublicUser struct {
		ID        uint      `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
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
