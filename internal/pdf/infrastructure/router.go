// internal/pdf/infrastructure/router.go
package infrastructure

import (
	pdfhttp "github.com/garfex/calculadora-filtros/internal/pdf/infrastructure/adapter/driver/http"
	"github.com/gin-gonic/gin"
)

// RegisterPdfRoutes monta todas las rutas del módulo PDF bajo el RouterGroup dado.
// Invocar desde main.go pasando el grupo /api/v1.
func RegisterPdfRoutes(rg *gin.RouterGroup, handler *pdfhttp.PdfHandler) {
	pdf := rg.Group("/pdf")
	{
		pdf.POST("/memoria", handler.GenerarMemoria)
	}
}
