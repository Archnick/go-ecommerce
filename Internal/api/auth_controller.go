package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// handleRegisterUser now takes a *gin.Context.
func (c *AuthController) handleRegisterUser(ctx *gin.Context) {
	var payload models.UserPayload
	// Use ctx.BindJSON() to parse the request body.
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		// The validation errors are now automatically translated.
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs.Translate(trans)})
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
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "This email is already in use"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Use ctx.JSON() to send a response. gin.H is a shortcut for map[string]interface{}.
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (c *AuthController) handleLogin(ctx *gin.Context) {
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

	expiry := time.Now().Add(7 * 24 * time.Hour)
	c.db.Model(&user).Updates(models.User{
		RefreshToken:          string(refreshToken),
		RefreshTokenExpiresAt: expiry,
	})

	ctx.JSON(http.StatusOK, gin.H{"access_token": tokenString, "refresh_token": refreshToken})
}

func (c *AuthController) handleRefreshToken(ctx *gin.Context) {
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
	if user.RefreshToken != body.RefreshToken || time.Now().After(user.RefreshTokenExpiresAt) {
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
