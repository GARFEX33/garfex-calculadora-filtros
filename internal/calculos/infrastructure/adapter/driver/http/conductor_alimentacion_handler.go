// internal/calculos/infrastructure/adapter/driver/http/conductor_alimentacion_handler.go
package http

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/gin-gonic/gin"
)

// ConductorAlimentacionHandler maneja el endpoint de conductor de alimentacion.
type ConductorAlimentacionHandler struct {
	useCase *usecase.SeleccionarConductorAlimentacionUseCase
}

// NewConductorAlimentacionHandler crea un nuevo handler.
func NewConductorAlimentacionHandler(uc *usecase.SeleccionarConductorAlimentacionUseCase) *ConductorAlimentacionHandler {
	return &ConductorAlimentacionHandler{
		useCase: uc,
	}
}

// ConductorAlimentacionRequest representa el body de la peticion POST.
type ConductorAlimentacionRequest struct {
	CorrienteAjustada float64 `json:"corriente_ajustada" binding:"required,gt=0"`
	TipoCanalizacion  string  `json:"tipo_canalizacion" binding:"required"`
	Material          string  `json:"material"`
	Temperatura       *int    `json:"temperatura"`
	HilosPorFase      int     `json:"hilos_por_fase"`
}

// ConductorAlimentacionResponse representa la respuesta exitosa.
type ConductorAlimentacionResponse struct {
	Success bool                            `json:"success"`
	Data    dto.ConductorAlimentacionOutput `json:"data"`
}

// ConductorAlimentacionResponseError representa la respuesta de error.
type ConductorAlimentacionResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SeleccionarConductorAlimentacion POST /api/v1/calculos/conductor-alimentacion
func (h *ConductorAlimentacionHandler) SeleccionarConductorAlimentacion(c *gin.Context) {
	var req ConductorAlimentacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Error de validacion",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Convertir request a DTO
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: req.CorrienteAjustada,
		TipoCanalizacion:  req.TipoCanalizacion,
		Material:          req.Material,
		Temperatura:       req.Temperatura,
		HilosPorFase:      req.HilosPorFase,
	}

	// Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, ConductorAlimentacionResponse{
		Success: true,
		Data:    output,
	})
}

// mapErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *ConductorAlimentacionHandler) mapErrorToResponse(err error) (int, ConductorAlimentacionResponseError) {
	// Errores 400 - Bad Request
	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Datos de entrada invalidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoCanalizacionInvalido) {
		return http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Tipo de canalizacion invalido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrCorrienteInvalida) {
		return http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Corriente invalida",
			Code:    "CORRIENTE_INVALIDA",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity
	if errors.Is(err, service.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "No se encontro conductor adecuado",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, ConductorAlimentacionResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
