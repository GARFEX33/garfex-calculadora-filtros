// internal/pdf/infrastructure/adapter/driver/http/pdf_handler.go
package http

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/garfex/calculadora-filtros/internal/pdf/application/dto"
	"github.com/garfex/calculadora-filtros/internal/pdf/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/pdf/domain"
	"github.com/gin-gonic/gin"
)

// reNonAlphaNum filtra caracteres no alfanuméricos ni guiones bajos/medios del filename.
var reNonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9_\-]`)

// PdfHandler maneja los endpoints de generación de PDF de memoria de cálculo.
type PdfHandler struct {
	generarMemoriaUC *usecase.GenerarMemoriaPdfUseCase
}

// NewPdfHandler crea un nuevo PdfHandler con el use case inyectado.
func NewPdfHandler(generarMemoriaUC *usecase.GenerarMemoriaPdfUseCase) *PdfHandler {
	return &PdfHandler{
		generarMemoriaUC: generarMemoriaUC,
	}
}

// pdfErrorResponse es el body de respuesta para errores de este handler.
type pdfErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// GenerarMemoria POST /api/v1/pdf/memoria
// @Summary Generar memoria de cálculo en PDF
// @Description Genera la memoria de cálculo eléctrica en PDF a partir del resultado de cálculo y datos de presentación
// @Tags PDF
// @Accept json
// @Produce application/pdf
// @Param request body dto.PdfMemoriaRequest true "Datos de cálculo y presentación"
// @Success 200 {file} binary "PDF generado exitosamente"
// @Failure 400 {object} pdfErrorResponse "Error de validación"
// @Failure 500 {object} pdfErrorResponse "Error al generar el PDF"
// @Router /pdf/memoria [post]
func (h *PdfHandler) GenerarMemoria(c *gin.Context) {
	var req dto.PdfMemoriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pdfErrorResponse{
			Success: false,
			Error:   "Error de validación del JSON",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Validar campos requeridos manualmente (no usan binding tags por ser structs anidados)
	if req.Presentacion.EmpresaID == "" {
		c.JSON(http.StatusBadRequest, pdfErrorResponse{
			Success: false,
			Error:   "El campo empresa_id es requerido",
			Code:    "EMPRESA_ID_REQUERIDO",
			Details: "presentacion.empresa_id no puede estar vacío",
		})
		return
	}

	if req.Presentacion.NombreProyecto == "" {
		c.JSON(http.StatusBadRequest, pdfErrorResponse{
			Success: false,
			Error:   "El campo nombre_proyecto es requerido",
			Code:    "NOMBRE_PROYECTO_REQUERIDO",
			Details: "presentacion.nombre_proyecto no puede estar vacío",
		})
		return
	}

	if req.Presentacion.Responsable == "" {
		c.JSON(http.StatusBadRequest, pdfErrorResponse{
			Success: false,
			Error:   "El campo responsable es requerido",
			Code:    "RESPONSABLE_REQUERIDO",
			Details: "presentacion.responsable no puede estar vacío",
		})
		return
	}

	// Validar que el empresa_id existe en el catálogo estático
	if _, ok := domain.BuscarEmpresaPorID(req.Presentacion.EmpresaID); !ok {
		c.JSON(http.StatusBadRequest, pdfErrorResponse{
			Success: false,
			Error:   "La empresa especificada no existe en el catálogo",
			Code:    "EMPRESA_NO_ENCONTRADA",
			Details: fmt.Sprintf("empresa_id=%q no es válido. Valores aceptados: garfex, summa, siemens", req.Presentacion.EmpresaID),
		})
		return
	}

	// Ejecutar el use case de generación de PDF
	pdfBytes, err := h.generarMemoriaUC.Execute(c.Request.Context(), req)
	if err != nil {
		status, resp := h.mapError(err)
		c.JSON(status, resp)
		return
	}

	// Construir nombre de archivo sanitizado
	filename := buildFilename(req.Presentacion.NombreProyecto, req.Memoria.Equipo.Clave)

	// Responder con el PDF
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// mapError convierte errores del dominio/aplicación a respuestas HTTP apropiadas.
func (h *PdfHandler) mapError(err error) (int, pdfErrorResponse) {
	if errors.Is(err, domain.ErrEmpresaNoEncontrada) {
		return http.StatusBadRequest, pdfErrorResponse{
			Success: false,
			Error:   "Empresa no encontrada en el catálogo",
			Code:    "EMPRESA_NO_ENCONTRADA",
			Details: err.Error(),
		}
	}

	if errors.Is(err, domain.ErrRenderizadoHtml) {
		return http.StatusInternalServerError, pdfErrorResponse{
			Success: false,
			Error:   "Error al renderizar la memoria de cálculo",
			Code:    "ERROR_RENDERIZADO_HTML",
			Details: err.Error(),
		}
	}

	if errors.Is(err, domain.ErrGeneracionPdf) {
		return http.StatusInternalServerError, pdfErrorResponse{
			Success: false,
			Error:   "Error al generar el PDF",
			Code:    "ERROR_GENERACION_PDF",
			Details: err.Error(),
		}
	}

	return http.StatusInternalServerError, pdfErrorResponse{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}

// buildFilename construye el nombre del archivo PDF sanitizado.
// Formato: MemoriaCalculo_<proyecto>_<equipo>_<fecha>.pdf
// Los espacios se convierten en guiones bajos y los caracteres especiales se eliminan.
func buildFilename(nombreProyecto, claveEquipo string) string {
	fecha := time.Now().Format("20060102")

	// Sanitizar proyecto: espacios → guión bajo, eliminar caracteres especiales
	proyecto := sanitizeFilenameSegment(nombreProyecto)
	if proyecto == "" {
		proyecto = "Proyecto"
	}

	// Sanitizar clave del equipo
	equipo := sanitizeFilenameSegment(claveEquipo)
	if equipo == "" {
		equipo = "Equipo"
	}

	return fmt.Sprintf("MemoriaCalculo_%s_%s_%s.pdf", proyecto, equipo, fecha)
}

// sanitizeFilenameSegment convierte espacios en guiones bajos y elimina caracteres especiales.
func sanitizeFilenameSegment(s string) string {
	// Reemplazar espacios con guión bajo
	s = strings.ReplaceAll(s, " ", "_")
	// Eliminar caracteres no alfanuméricos (excepto _ y -)
	s = reNonAlphaNum.ReplaceAllString(s, "")
	return s
}
