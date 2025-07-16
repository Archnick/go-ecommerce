package api

import (
	"net/http"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin" // Import Gin
	"github.com/golang-jwt/jwt/v5"
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

func (c *UsersController) handleLogin(ctx *gin.Context) {
	var payload models.UserPayload
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	// Find the user by email
	result := c.db.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate the JWT
	tokenString, err := GenerateAccessToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process token"})
		return
	}

	expiry := time.Now().Add(7 * 24 * time.Hour)
	c.db.Model(&user).Updates(models.User{
		RefreshToken:          string(hashedRefreshToken),
		RefreshTokenExpiresAt: expiry,
	})

	ctx.JSON(http.StatusOK, gin.H{"access_token": tokenString, "refresh_token": refreshToken})
}

func (c *UsersController) handleRefreshToken(ctx *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 1. Parse the incoming refresh token to get the claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(body.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// 2. Fetch the specific user from the database using the ID from the token claims
	var user models.User
	result := c.db.First(&user, claims.UserID)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// 3. Compare the stored refresh token hash with the incoming one
	err = bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(body.RefreshToken))
	if err != nil || time.Now().After(user.RefreshTokenExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// 4. Generate a new access token
	newAccessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
