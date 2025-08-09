package api

import (
	"reflect"
	"strings"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/gin-gonic/gin" // Import Gin
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"gorm.io/gorm"
)

var trans ut.Translator

// Server holds the dependencies for our API.
type Server struct {
	db     *gorm.DB
	router *gin.Engine // The router is now a Gin Engine
}

// NewServer creates a new Server instance with Gin.
func NewServer(db *gorm.DB) *Server {
	// gin.Default() creates a Gin router with default middleware (logger, recovery).
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Use JSON tag name for field names in errors
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Setup translator
		en := en.New()
		uni := ut.New(en, en)
		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, trans)
	}

	s := &Server{
		db:     db,
		router: router,
	}
	s.routes()
	return s
}

// Start runs the HTTP server using the Gin engine.
func (s *Server) Start(addr string) error {
	// The Run method starts the server.
	return s.router.Run(addr)
}

// routes sets up all the routing for the application using Gin's syntax.
func (s *Server) routes() {
	api := s.router.Group("/api")
	s.getAuthRoutes(api)
	s.getUserRoutes(api)
	s.getProductRoutes(api)
	s.getCategoryRoutes(api)
	s.getShopRoutes(api)
}

func (s *Server) getAuthRoutes(api *gin.RouterGroup) {
	authController := NewAuthController(s.db)
	api.POST("/register", authController.handleRegisterUser)
	api.POST("/login", authController.handleLogin)
	api.POST("/refresh", authController.handleRefreshToken)
	api.POST("/logout", AuthMiddleware(), authController.handleLogout)
}

func (s *Server) getUserRoutes(api *gin.RouterGroup) {
	usersController := NewUsersController(s.db)
	ReviewController := NewReviewController(s.db)
	api.GET("/users", AuthMiddleware(), usersController.handleGetUsers)
	api.GET("/users/:id", AuthMiddleware(), usersController.handleGetUser)
	api.GET("/users/:user_id/reviews", ReviewController.handleGetReviewsForUser)
	api.PUT("/users/:id", AuthMiddleware(), usersController.handleUpdateUser)
	api.DELETE("/users/:id", AuthMiddleware(), usersController.handleDeleteUser)
}

func (s *Server) getProductRoutes(api *gin.RouterGroup) {
	productController := NewProductController(s.db)
	productImageController := NewProductImageController(s.db)
	reviewController := NewReviewController(s.db)
	api.GET("/products", productController.handleGetProducts)
	api.GET("/products/:id", productController.handleGetProduct)
	api.GET("/products/:product_id/images", productImageController.handleGetProductImages)
	api.GET("/products/:product_id/reviews", reviewController.handleGetReviewsForProduct)
	api.POST("/products", AuthMiddleware(), productController.handleCreateProduct)
	api.POST("/products/:product_id/images", AuthMiddleware(), productImageController.handleCreateProductImage)
	api.POST("/products/:product_id/reviews", AuthMiddleware(), reviewController.handleCreateReview)
	api.PUT("/products/:id", AuthMiddleware(), productController.handleUpdateProduct)
	api.DELETE("/products/:id", AuthMiddleware(), productController.handleDeleteProduct)

	api.DELETE("/product_images/:image_id", AuthMiddleware(), productImageController.handleDeleteProductImage)
	api.PUT("/reviews/:review_id", AuthMiddleware(), reviewController.handleUpdateReview)
}

func (s *Server) getCategoryRoutes(api *gin.RouterGroup) {
	categoryController := NewCategoryController(s.db)
	api.GET("/categories", categoryController.handleGetCategories)
	api.POST("/categories", AuthMiddleware(), RoleMiddleware(models.AdminRole), categoryController.handleCreateCategory)
	api.PUT("/categories/:id", AuthMiddleware(), RoleMiddleware(models.AdminRole), categoryController.handleUpdateCategory)
	api.DELETE("/categories/:id", AuthMiddleware(), RoleMiddleware(models.AdminRole), categoryController.handleDeleteCategory)
}

func (s *Server) getShopRoutes(api *gin.RouterGroup) {
	shopController := NewShopController(s.db)
	api.GET("/shops", shopController.handleGetShops)
	api.GET("/shops/:id", shopController.handleGetShop)
	api.POST("/shops", AuthMiddleware(), RoleMiddleware(models.AdminRole), shopController.handleCreateShop)
	api.PUT("/shops/:id", AuthMiddleware(), RoleMiddleware(models.AdminRole), shopController.handleUpdateShop)
	api.DELETE("/shops/:id", AuthMiddleware(), RoleMiddleware(models.AdminRole), shopController.handleDeleteShop)
}
