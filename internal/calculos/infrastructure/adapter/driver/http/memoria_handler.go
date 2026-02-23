// internal/calculos/infrastructure/adapter/driver/http/memoria_handler.go
package http

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/gin-gonic/gin"
)

// MemoriaHandler handles the complete memory calculation endpoint.
type MemoriaHandler struct {
	orquestadorUC *usecase.OrquestadorMemoriaCalculoUseCase
}

// NewMemoriaHandler creates a new memoria handler.
func NewMemoriaHandler(orquestadorUC *usecase.OrquestadorMemoriaCalculoUseCase) *MemoriaHandler {
	return &MemoriaHandler{
		orquestadorUC: orquestadorUC,
	}
}

// CalcularMemoriaRequest represents the request body for the memoria endpoint.
type CalcularMemoriaRequest struct {
	// Mode indicates how equipment data is provided
	Modo dto.ModoCalculo `json:"modo" binding:"required"`

	// Equipment key (required if Modo = LISTADO)
	Clave string `json:"clave"`

	// Equipment data (required if Modo = MANUAL_*)
	TipoEquipo      string  `json:"tipo_equipo"`
	AmperajeNominal float64 `json:"amperaje_nominal"`
	PotenciaNominal float64 `json:"potencia_nominal"`
	PotenciaUnidad  string  `json:"potencia_unidad"` // "W", "KW", "KVA", "KVAR"
	Tension         float64 `json:"tension" binding:"required,gt=0"`
	TensionUnidad   string  `json:"tension_unidad"` // "V" o "kV" (default: "V")
	FactorPotencia  float64 `json:"factor_potencia"`
	ITM             int     `json:"itm" binding:"required,gt=0"`

	// Installation parameters
	TipoCanalizacion      string   `json:"tipo_canalizacion" binding:"required"`
	TemperaturaOverride   *int     `json:"temperatura_override,omitempty"`
	HilosPorFase          int      `json:"hilos_por_fase"`
	NumTuberias           int      `json:"num_tuberias"`
	Material              string   `json:"material"`
	LongitudCircuito      float64  `json:"longitud_circuito" binding:"required,gt=0"`
	PorcentajeCaidaMaximo float64  `json:"porcentaje_caida_maximo"`
	DiametroControlMM     *float64 `json:"diametro_control_mm,omitempty"`

	// Electrical system
	SistemaElectrico dto.SistemaElectrico `json:"sistema_electrico" binding:"required"`
	Estado           string               `json:"estado" binding:"required"`

	// Voltage type (FASE_NEUTRO or FASE_FASE)
	TipoVoltaje string `json:"tipo_voltaje" binding:"required"`
}

// CalcularMemoriaResponse represents the response for the memoria endpoint.
type CalcularMemoriaResponse struct {
	Success bool              `json:"success"`
	Data    dto.MemoriaOutput `json:"data"`
}

// CalcularMemoriaResponseError represents an error response.
type CalcularMemoriaResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// CalcularMemoria POST /api/v1/calculos/memoria
// Executes the complete memory calculation pipeline.
func (h *MemoriaHandler) CalcularMemoria(c *gin.Context) {
	var req CalcularMemoriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Build EquipoInput from request
	input := dto.EquipoInput{
		Modo:                  req.Modo,
		Clave:                 req.Clave,
		TipoEquipo:            req.TipoEquipo,
		AmperajeNominal:       req.AmperajeNominal,
		PotenciaNominal:       req.PotenciaNominal,
		PotenciaUnidad:        req.PotenciaUnidad,
		Tension:               req.Tension,
		TensionUnidad:         req.TensionUnidad,
		FactorPotencia:        req.FactorPotencia,
		ITM:                   req.ITM,
		TipoCanalizacion:      req.TipoCanalizacion,
		TemperaturaOverride:   req.TemperaturaOverride,
		HilosPorFase:          req.HilosPorFase,
		NumTuberias:           req.NumTuberias,
		Material:              req.Material,
		LongitudCircuito:      req.LongitudCircuito,
		PorcentajeCaidaMaximo: req.PorcentajeCaidaMaximo,
		DiametroControlMM:     req.DiametroControlMM,
		SistemaElectrico:      req.SistemaElectrico,
		Estado:                req.Estado,
		TipoVoltaje:           req.TipoVoltaje,
	}

	// Execute orchestrator
	result, err := h.orquestadorUC.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CalcularMemoriaResponse{
		Success: true,
		Data:    result,
	})
}

// mapErrorToResponse maps domain/application errors to HTTP responses.
func (h *MemoriaHandler) mapErrorToResponse(err error) (int, CalcularMemoriaResponseError) {
	// Validation errors (400)
	if errors.Is(err, dto.ErrEquipoInputInvalido) ||
		errors.Is(err, dto.ErrModoInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Datos de entrada inválidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	// Entity validation errors (400)
	if errors.Is(err, entity.ErrTipoCanalizacionInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Tipo de canalización inválido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrSistemaElectricoInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Sistema eléctrico inválido",
			Code:    "SISTEMA_ELECTRICO_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoVoltajeInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Tipo de voltaje inválido",
			Code:    "TIPO_VOLTAJE_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrVoltajeInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Tensión fuera de rango permitido",
			Code:    "TENSION_INVALIDA",
			Details: err.Error(),
		}
	}

	// Business logic errors (422)
	if errors.Is(err, dto.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "No se encontró conductor con las características solicitadas",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, dto.ErrCanalizacionNoDisponible) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Canalización no disponible para los parámetros dados",
			Code:    "CANALIZACION_NO_DISPONIBLE",
			Details: err.Error(),
		}
	}

	// Default: internal server error (500)
	return http.StatusInternalServerError, CalcularMemoriaResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
