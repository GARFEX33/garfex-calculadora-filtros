// internal/calculos/infrastructure/router.go
package infrastructure

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driver/http"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driver/http/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter crea y configura el router Gin.
func NewRouter(
	calcularCorrienteUC *usecase.CalcularCorrienteUseCase,
	ajustarCorrienteUC *usecase.AjustarCorrienteUseCase,
	seleccionarConductorAlimentacionUC *usecase.SeleccionarConductorAlimentacionUseCase,
	seleccionarConductorTierraUC *usecase.SeleccionarConductorTierraUseCase,
) *gin.Engine {
	router := gin.New()

	// Middlewares globales
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.RequestLogger())

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
			calculoHandler := http.NewCalculoHandler(calcularCorrienteUC, ajustarCorrienteUC)
			calculos.POST("/amperaje", calculoHandler.CalcularAmperaje)
			calculos.POST("/corriente-ajustada", calculoHandler.CalcularCorrienteAjustada)

			// Conductor de alimentacion
			conductorAlimentacionHandler := http.NewConductorAlimentacionHandler(seleccionarConductorAlimentacionUC)
			calculos.POST("/conductor-alimentacion", conductorAlimentacionHandler.SeleccionarConductorAlimentacion)

			// Conductor de tierra
			conductorTierraHandler := http.NewConductorTierraHandler(seleccionarConductorTierraUC)
			calculos.POST("/conductor-tierra", conductorTierraHandler.SeleccionarConductorTierra)
		}
	}

	return router
}
