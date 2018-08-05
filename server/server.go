package server

import (
	"products/config"
	"products/controllers/v1"
	"products/handlers"
	"products/middleware"

	"github.com/gin-gonic/gin"
)

// Server is the http layer for role and user resource
type Server struct {
	config            *config.Config
	productController *controllers.ProductController
	middleware        *middleware.Middleware
	handlers          *handlers.HttpHandlers
}

// NewServer is the Server constructor
func NewServer(cf *config.Config,
	pc *controllers.ProductController,
	mid *middleware.Middleware,
	hand *handlers.HttpHandlers) *Server {

	return &Server{
		config:            cf,
		productController: pc,
		middleware:        mid,
		handlers:          hand,
	}
}

// Run loads server with its routes and starts the server
func (s *Server) Run() {
	// Instantiate a new router
	r := gin.Default()

	// cors
	r.Use(*s.middleware.Cors())

	// generic routes
	r.HandleMethodNotAllowed = false
	r.NoRoute(s.handlers.NotFound)

	// Product resource
	productApi := r.Group("/api/v1/product")
	{
		// Create a new product
		productApi.POST("", s.productController.CreateAction)
	}

	// Fire up the server
	r.Run(s.config.Host)
}
