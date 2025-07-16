package api

import (
	"github.com/gin-gonic/gin" // Import Gin
	"gorm.io/gorm"
)

// Server holds the dependencies for our API.
type Server struct {
	db     *gorm.DB
	router *gin.Engine // The router is now a Gin Engine
}

// NewServer creates a new Server instance with Gin.
func NewServer(db *gorm.DB) *Server {
	// gin.Default() creates a Gin router with default middleware (logger, recovery).
	router := gin.Default()
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
	usersController := NewUsersController(s.db)

	// Group routes under /api
	api := s.router.Group("/api")
	{
		api.POST("/register", usersController.handleRegisterUser)
		api.POST("/login", usersController.handleLogin)
		api.POST("/refresh", usersController.handleRefreshToken)

		api.GET("/users", AuthMiddleware(), usersController.handleGetUsers)

	}
}
