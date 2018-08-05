package middleware

import (
	"products/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Middleware provides necessary middleware to the http layer
type Middleware struct {
	config *config.Config
}

// NewMiddleware is the constructor of Middleware
func NewMiddleware(cf *config.Config) *Middleware {
	return &Middleware{
		config: cf,
	}
}

// Cors returns a handler to deal with CORS
func (m *Middleware) Cors() *gin.HandlerFunc {
	corsMid := cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:4200", "http://127.0.0.1:4200"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept-Encoding", "Accept-Language", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Cache-Control", "Connection",
			"Host", "Origin", "Pragma", "User-Agent", "X-Custom-Header", "access-control-allow-origin", "authorization", "Origin", "Content-Type", "Accept", "Key", "Keep-Alive", "User-Agent", "If-Modified-Since", "Cache-Control", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           48 * time.Hour,
	})

	return &corsMid
}
