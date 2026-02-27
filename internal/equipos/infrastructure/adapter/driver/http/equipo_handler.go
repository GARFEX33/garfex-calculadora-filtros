// internal/equipos/infrastructure/adapter/driver/http/equipo_handler.go
package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/usecase"
	"github.com/gin-gonic/gin"
)

// EquipoHandler handles HTTP requests for the equipos feature.
type EquipoHandler struct {
	crearUC      *usecase.CrearEquipoUseCase
	obtenerUC    *usecase.ObtenerEquipoUseCase
	listarUC     *usecase.ListarEquiposUseCase
	actualizarUC *usecase.ActualizarEquipoUseCase
	eliminarUC   *usecase.EliminarEquipoUseCase
}

// NewEquipoHandler creates a new handler with all required use cases.
func NewEquipoHandler(
	crearUC *usecase.CrearEquipoUseCase,
	obtenerUC *usecase.ObtenerEquipoUseCase,
	listarUC *usecase.ListarEquiposUseCase,
	actualizarUC *usecase.ActualizarEquipoUseCase,
	eliminarUC *usecase.EliminarEquipoUseCase,
) *EquipoHandler {
	return &EquipoHandler{
		crearUC:      crearUC,
		obtenerUC:    obtenerUC,
		listarUC:     listarUC,
		actualizarUC: actualizarUC,
		eliminarUC:   eliminarUC,
	}
}

// ─── Response types ──────────────────────────────────────────────────────────

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type errorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
	Details   string `json:"details,omitempty"`
	Timestamp string `json:"timestamp"` // ISO 8601 UTC
}

func ok(data any) successResponse {
	return successResponse{Success: true, Data: data}
}

func errResp(msg, code, details string) errorResponse {
	return errorResponse{
		Success:   false,
		Error:     msg,
		Code:      code,
		Details:   details,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ─── Endpoints ───────────────────────────────────────────────────────────────

// Crear POST /api/v1/equipos
func (h *EquipoHandler) Crear(c *gin.Context) {
	var input dto.CreateEquipoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errResp("Error de validación", "VALIDATION_ERROR", err.Error()))
		return
	}

	output, err := h.crearUC.Execute(c.Request.Context(), input)
	if err != nil {
		status, resp := mapError(err)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusCreated, ok(output))
}

// ObtenerPorID GET /api/v1/equipos/:id
func (h *EquipoHandler) ObtenerPorID(c *gin.Context) {
	id := c.Param("id")

	output, err := h.obtenerUC.Execute(c.Request.Context(), id)
	if err != nil {
		status, resp := mapError(err)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusOK, ok(output))
}

// Listar GET /api/v1/equipos
// Query params: tipo, voltaje, buscar, page (default 1), page_size (default 20, max 100)
func (h *EquipoHandler) Listar(c *gin.Context) {
	query := dto.ListEquiposQuery{
		Tipo:   c.Query("tipo"),
		Buscar: c.Query("buscar"),
	}

	if voltajeStr := c.Query("voltaje"); voltajeStr != "" {
		v, err := strconv.Atoi(voltajeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, errResp("Voltaje inválido", "VOLTAJE_INVALIDO", "debe ser un número entero"))
			return
		}
		query.Voltaje = v
	}

	if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, errResp("Página inválida", "PAGE_INVALIDO", "debe ser un entero mayor que cero"))
			return
		}
		query.Page = p
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err != nil || ps < 1 {
			c.JSON(http.StatusBadRequest, errResp("Tamaño de página inválido", "PAGE_SIZE_INVALIDO", "debe ser un entero mayor que cero"))
			return
		}
		query.PageSize = ps
	}

	output, err := h.listarUC.Execute(c.Request.Context(), query)
	if err != nil {
		status, resp := mapError(err)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusOK, ok(output))
}

// Actualizar PUT /api/v1/equipos/:id
func (h *EquipoHandler) Actualizar(c *gin.Context) {
	id := c.Param("id")

	var input dto.UpdateEquipoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errResp("Error de validación", "VALIDATION_ERROR", err.Error()))
		return
	}

	output, err := h.actualizarUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		status, resp := mapError(err)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusOK, ok(output))
}

// Eliminar DELETE /api/v1/equipos/:id
func (h *EquipoHandler) Eliminar(c *gin.Context) {
	id := c.Param("id")

	if err := h.eliminarUC.Execute(c.Request.Context(), id); err != nil {
		status, resp := mapError(err)
		c.JSON(status, resp)
		return
	}

	c.Status(http.StatusNoContent)
}

// ─── Error mapper ────────────────────────────────────────────────────────────

// mapError converts application/domain errors to HTTP status + response body.
func mapError(err error) (int, errorResponse) {
	switch {
	case errors.Is(err, dto.ErrIDInvalido):
		return http.StatusBadRequest, errResp("ID inválido", "ID_INVALIDO", err.Error())

	case errors.Is(err, dto.ErrInputInvalido):
		return http.StatusBadRequest, errResp("Datos de entrada inválidos", "INPUT_INVALIDO", err.Error())

	case errors.Is(err, dto.ErrEquipoNoEncontrado):
		return http.StatusNotFound, errResp("Equipo no encontrado", "EQUIPO_NO_ENCONTRADO", err.Error())

	case errors.Is(err, dto.ErrClaveYaExiste):
		return http.StatusConflict, errResp("La clave ya existe", "CLAVE_DUPLICADA", err.Error())

	default:
		return http.StatusInternalServerError, errResp("Error interno del servidor", "INTERNAL_ERROR", err.Error())
	}
}
