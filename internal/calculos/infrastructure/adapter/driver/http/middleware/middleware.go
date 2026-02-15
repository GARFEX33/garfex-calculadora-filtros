// internal/calculos/infrastructure/adapter/driver/http/middleware/middleware.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// corsMiddleware configura CORS para desarrollo.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
