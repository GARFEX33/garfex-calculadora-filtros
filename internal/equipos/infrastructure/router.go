// internal/equipos/infrastructure/router.go
package infrastructure

import (
	equipohttp "github.com/garfex/calculadora-filtros/internal/equipos/infrastructure/adapter/driver/http"
	"github.com/gin-gonic/gin"
)

// RegisterEquiposRoutes mounts all equipo routes under the given RouterGroup.
// Call this from main.go passing the /api/v1 group.
func RegisterEquiposRoutes(rg *gin.RouterGroup, handler *equipohttp.EquipoHandler) {
	equipos := rg.Group("/equipos")
	{
		equipos.POST("", handler.Crear)
		equipos.GET("", handler.Listar)
		equipos.GET("/:id", handler.ObtenerPorID)
		equipos.PUT("/:id", handler.Actualizar)
		equipos.DELETE("/:id", handler.Eliminar)
	}
}
