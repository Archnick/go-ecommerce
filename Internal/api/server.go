package api

import (
	"reflect"
	"strings"

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
	api.GET("/users", AuthMiddleware(), usersController.handleGetUsers)
	api.GET("/users/:id", AuthMiddleware(), usersController.handleGetUser)
	api.PUT("/users/:id", AuthMiddleware(), usersController.handleUpdateUser)
	api.DELETE("/users/:id", AuthMiddleware(), usersController.handleDeleteUser)
}

func (s *Server) getProductRoutes(api *gin.RouterGroup) {
	productController := NewProductController(s.db)
	api.GET("/products", productController.handleGetProducts)
	api.GET("/products/:id", productController.handleGetProduct)
	api.POST("/products", AuthMiddleware(), productController.handleCreateProduct)
	api.PUT("/products/:id", AuthMiddleware(), productController.handleUpdateProduct)
	api.DELETE("/products/:id", AuthMiddleware(), productController.handleDeleteProduct)
}
