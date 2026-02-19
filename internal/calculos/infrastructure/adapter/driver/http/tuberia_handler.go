// internal/calculos/infrastructure/adapter/driver/http/tuberia_handler.go
package http

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/gin-gonic/gin"
)

// TuberiaHandler maneja el endpoint de cálculo de tamaño de tubería.
type TuberiaHandler struct {
	calcularTamanioTuberiaUC *usecase.CalcularTamanioTuberiaUseCase
}

// NewTuberiaHandler crea un nuevo handler de tubería.
func NewTuberiaHandler(calcularTamanioTuberiaUC *usecase.CalcularTamanioTuberiaUseCase) *TuberiaHandler {
	return &TuberiaHandler{
		calcularTamanioTuberiaUC: calcularTamanioTuberiaUC,
	}
}

// CalcularTuberiaRequest representa el body de la petición POST /tuberia.
type CalcularTuberiaRequest struct {
	NumFases         int    `json:"num_fases" binding:"required,gt=0"`
	CalibreFase      string `json:"calibre_fase" binding:"required"`
	NumNeutros       int    `json:"num_neutros" binding:"gte=0"`
	CalibreNeutro    string `json:"calibre_neutral"`
	CalibreTierra    string `json:"calibre_tierra" binding:"required"`
	TipoCanalizacion string `json:"tipo_canalizacion" binding:"required"`
	NumTuberias      int    `json:"num_tuberias" binding:"required,gt=0"`
}

// CalcularTuberiaResponse representa la respuesta exitosa.
type CalcularTuberiaResponse struct {
	Success bool              `json:"success"`
	Data    dto.TuberiaOutput `json:"data"`
}

// CalcularTuberiaResponseError representa la respuesta de error.
type CalcularTuberiaResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// CalcularTuberia POST /api/v1/calculos/tuberia
func (h *TuberiaHandler) CalcularTuberia(c *gin.Context) {
	var req CalcularTuberiaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CalcularTuberiaResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Map request to DTO
	input := dto.TuberiaInput{
		NumFases:         req.NumFases,
		CalibreFase:      req.CalibreFase,
		NumNeutros:       req.NumNeutros,
		CalibreNeutro:    req.CalibreNeutro,
		CalibreTierra:    req.CalibreTierra,
		TipoCanalizacion: req.TipoCanalizacion,
		NumTuberias:      req.NumTuberias,
	}

	// Execute use case
	output, err := h.calcularTamanioTuberiaUC.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapTuberiaErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CalcularTuberiaResponse{
		Success: true,
		Data:    output,
	})
}

// mapTuberiaErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *TuberiaHandler) mapTuberiaErrorToResponse(err error) (int, CalcularTuberiaResponseError) {
	// Errores 400 - Bad Request
	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, CalcularTuberiaResponseError{
			Success: false,
			Error:   "Datos de entrada inválidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity
	if errors.Is(err, dto.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, CalcularTuberiaResponseError{
			Success: false,
			Error:   "Conductor no encontrado en tablas NOM",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, CalcularTuberiaResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
