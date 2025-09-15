package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/Archnick/go-ecommerce/Internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	service *services.UserService
}

func NewAuthController(db *services.UserService) *AuthController {
	return &AuthController{service: db}
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

	user, err := c.service.RegisterUser(payload)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "This email is already in use"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Use ctx.JSON() to send a response. gin.H is a shortcut for map[string]interface{}.
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func (c *AuthController) handleLogin(ctx *gin.Context) {
	var payload models.UserPayload
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := c.service.AuthorizeUser(payload)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate the JWT
	tokenString, err := GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.service.SetRefreshToken(user, string(refreshToken))

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

	// Parse the incoming refresh token to get the claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(body.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Fetch the specific user from the database using the ID from the token claims
	user, err := c.service.GetUserByID(claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Compare the stored refresh token hash with the incoming one
	if user.RefreshToken != body.RefreshToken || time.Now().After(user.RefreshTokenExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Generate a new access token
	newAccessToken, err := GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	// Generate a NEW refresh token.
	newRefreshToken, err := GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	// Update the database with the new refresh token.
	c.service.SetRefreshToken(user, newRefreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (c *AuthController) handleLogout(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "You are not logged in"})
		return
	}

	// Find the user and clear their refresh token fields.
	// TODO: I still need to find a solution for invalidating the Auth token
	user, err := c.service.GetUserByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.service.SetRefreshToken(user, "")
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
