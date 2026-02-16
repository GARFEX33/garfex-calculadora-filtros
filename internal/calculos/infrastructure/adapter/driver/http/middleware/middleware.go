// internal/calculos/infrastructure/adapter/driver/http/middleware/middleware.go
package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// corsMiddleware configura CORS.
// Por defecto usa "*" para desarrollo.
// En producción, configurar variable CORS_ORIGINS con lista de dominios separados por coma.
func CorsMiddleware() gin.HandlerFunc {
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLogger loguea todas las peticiones.
func RequestLogger() gin.HandlerFunc {
	return gin.Logger()
}
