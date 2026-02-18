// internal/calculos/infrastructure/adapter/driver/http/conductor_tierra_handler.go
package http

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/gin-gonic/gin"
)

// ConductorTierraHandler maneja el endpoint de conductor de tierra.
type ConductorTierraHandler struct {
	useCase *usecase.SeleccionarConductorTierraUseCase
}

// NewConductorTierraHandler crea un nuevo handler de conductor de tierra.
func NewConductorTierraHandler(uc *usecase.SeleccionarConductorTierraUseCase) *ConductorTierraHandler {
	return &ConductorTierraHandler{
		useCase: uc,
	}
}

// ConductorTierraRequest representa el body de la petición POST.
type ConductorTierraRequest struct {
	// ITM es el Interruptor Termomagnético en amperes (requerido, > 0).
	ITM int `json:"itm" binding:"required,gt=0"`
	// Material es el material del conductor ("Cu" o "Al").
	// Si está vacío, se usa "Cu" por defecto.
	Material string `json:"material"`
}

// ConductorTierraResponse representa la respuesta exitosa.
type ConductorTierraResponse struct {
	Success bool                      `json:"success"`
	Data    dto.ConductorTierraOutput `json:"data"`
}

// ConductorTierraResponseError representa la respuesta de error.
type ConductorTierraResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SeleccionarConductorTierra POST /api/v1/calculos/conductor-tierra
func (h *ConductorTierraHandler) SeleccionarConductorTierra(c *gin.Context) {
	var req ConductorTierraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ConductorTierraResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), req.ITM, req.Material)
	if err != nil {
		status, response := h.mapErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, ConductorTierraResponse{
		Success: true,
		Data:    output,
	})
}

// mapErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *ConductorTierraHandler) mapErrorToResponse(err error) (int, ConductorTierraResponseError) {
	// Errores 400 - Bad Request (validación de input)
	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, ConductorTierraResponseError{
			Success: false,
			Error:   "Datos de entrada inválidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity
	if errors.Is(err, dto.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, ConductorTierraResponseError{
			Success: false,
			Error:   "No se encontró conductor de tierra adecuado",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, ConductorTierraResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
