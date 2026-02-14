// internal/presentation/router.go
package presentation

import (
	"github.com/garfex/calculadora-filtros/internal/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/presentation/handler"
	"github.com/gin-gonic/gin"
)

// NewRouter crea y configura el router Gin.
func NewRouter(calcularMemoriaUC *usecase.CalcularMemoriaUseCase) *gin.Engine {
	router := gin.New()

	// Middlewares globales
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(requestLogger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Calculos
		calculos := v1.Group("/calculos")
		{
			calculoHandler := handler.NewCalculoHandler(calcularMemoriaUC)
			calculos.POST("/memoria", calculoHandler.CalcularMemoria)
		}
	}

	return router
}

// corsMiddleware configura CORS para desarrollo.
func corsMiddleware() gin.HandlerFunc {
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

// requestLogger loguea todas las peticiones.
func requestLogger() gin.HandlerFunc {
	return gin.Logger()
}
