// internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go
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

// CalculoHandler maneja los endpoints de cálculo.
type CalculoHandler struct {
	calcularCorrienteUseCase *usecase.CalcularCorrienteUseCase
	ajustarCorrienteUseCase  *usecase.AjustarCorrienteUseCase
}

// NewCalculoHandler crea un nuevo handler de cálculo.
func NewCalculoHandler(
	calcularCorrienteUC *usecase.CalcularCorrienteUseCase,
	ajustarCorrienteUC *usecase.AjustarCorrienteUseCase,
) *CalculoHandler {
	return &CalculoHandler{
		calcularCorrienteUseCase: calcularCorrienteUC,
		ajustarCorrienteUseCase:  ajustarCorrienteUC,
	}
}

// ============================================
// Endpoint: Calcular Amperaje Nominal
// ============================================

// CalcularAmperajeRequest representa el body de la petición POST /amperaje.
type CalcularAmperajeRequest struct {
	// PotenciaWatts es la potencia activa en Watts (requerido, > 0)
	PotenciaWatts float64 `json:"potencia_watts" binding:"required,gt=0"`

	// Tension es la tensión del circuito en volts (requerido)
	Tension int `json:"tension" binding:"required"`

	// TipoCarga indica el tipo de carga eléctrica (requerido)
	// Valores: "MONOFASICA" | "TRIFASICA"
	TipoCarga string `json:"tipo_carga" binding:"required"`

	// SistemaElectrico indica el tipo de sistema eléctrico (requerido)
	// Valores: "ESTRELLA" | "DELTA"
	SistemaElectrico string `json:"sistema_electrico" binding:"required"`

	// FactorPotencia es el factor de potencia (requerido, > 0 y <= 1)
	FactorPotencia float64 `json:"factor_potencia" binding:"required,gt=0,lte=1"`
}

// CalcularAmperajeResponse representa la respuesta exitosa.
type CalcularAmperajeResponse struct {
	Success  bool    `json:"success"`
	Amperaje float64 `json:"amperaje"`
	Unidad   string  `json:"unidad"`
}

// CalcularAmperajeResponseError representa la respuesta de error.
type CalcularAmperajeResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// CalcularAmperaje POST /api/v1/calculos/amperaje
// Usa el modo MANUAL_POTENCIA del CalcularCorrienteUseCase existente.
func (h *CalculoHandler) CalcularAmperaje(c *gin.Context) {
	var req CalcularAmperajeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CalcularAmperajeResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Validar tensión - debe ser valor NOM válido
	if req.Tension <= 0 {
		c.JSON(http.StatusBadRequest, CalcularAmperajeResponseError{
			Success: false,
			Error:   "Tensión inválida",
			Code:    "TENSION_INVALIDA",
			Details: "La tensión debe ser un valor positivo",
		})
		return
	}

	// Convertir request a EquipoInput con modo MANUAL_POTENCIA
	input := dto.EquipoInput{
		Modo:             dto.ModoManualPotencia,
		PotenciaNominal:  req.PotenciaWatts,
		Tension:          float64(req.Tension),
		SistemaElectrico: dto.SistemaElectrico(req.SistemaElectrico),
		FactorPotencia:   req.FactorPotencia,
		ITM:              0, // No requerido para este cálculo
		Estado:           "default",
	}

	// Ejecutar use case de corriente en modo MANUAL_POTENCIA
	output, err := h.calcularCorrienteUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapAmperajeErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CalcularAmperajeResponse{
		Success:  true,
		Amperaje: output.CorrienteNominal,
		Unidad:   "A",
	})
}

// mapAmperajeErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *CalculoHandler) mapAmperajeErrorToResponse(err error) (int, CalcularAmperajeResponseError) {
	// Errores 400 - Bad Request
	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, CalcularAmperajeResponseError{
			Success: false,
			Error:   "Datos de entrada inválidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrVoltajeInvalido) {
		return http.StatusBadRequest, CalcularAmperajeResponseError{
			Success: false,
			Error:   "Tensión fuera de rango permitido",
			Code:    "TENSION_INVALIDA",
			Details: err.Error(),
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, CalcularAmperajeResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}

// ============================================
// Endpoint: Calcular Corriente Ajustada
// ============================================

// CorrienteAjustadaRequest representa el body de la petición POST /corriente-ajustada.
type CorrienteAjustadaRequest struct {
	CorrienteNominal float64 `json:"corriente_nominal" binding:"required,gt=0"`
	Estado           string  `json:"estado" binding:"required"`
	TipoCanalizacion string  `json:"tipo_canalizacion" binding:"required"`
	SistemaElectrico string  `json:"sistema_electrico" binding:"required"`
	TipoEquipo       string  `json:"tipo_equipo" binding:"required"`
	HilosPorFase     int     `json:"hilos_por_fase" binding:"gte=1"`
	NumTuberias      int     `json:"num_tuberias" binding:"gte=1"`
}

// CorrienteAjustadaResponse representa la respuesta exitosa.
type CorrienteAjustadaResponse struct {
	Success bool                         `json:"success"`
	Data    dto.ResultadoAjusteCorriente `json:"data"`
}

// CorrienteAjustadaResponseError representa la respuesta de error.
type CorrienteAjustadaResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// CalcularCorrienteAjustada POST /api/v1/calculos/corriente-ajustada
func (h *CalculoHandler) CalcularCorrienteAjustada(c *gin.Context) {
	var req CorrienteAjustadaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Parsear tipo de canalización
	tipoCanalizacion, err := entity.ParseTipoCanalizacion(req.TipoCanalizacion)
	if err != nil {
		c.JSON(http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Tipo de canalización inválido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Parsear sistema eléctrico
	sistemaElectrico, err := entity.ParseSistemaElectrico(req.SistemaElectrico)
	if err != nil {
		c.JSON(http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Sistema eléctrico inválido",
			Code:    "SISTEMA_ELECTRICO_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Parsear tipo de equipo
	tipoEquipo, err := entity.ParseTipoEquipo(req.TipoEquipo)
	if err != nil {
		c.JSON(http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Tipo de equipo inválido",
			Code:    "TIPO_EQUIPO_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Crear corriente nominal
	corrienteNominal, err := valueobject.NewCorriente(req.CorrienteNominal)
	if err != nil {
		c.JSON(http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Corriente nominal inválida",
			Code:    "CORRIENTE_INVALIDA",
			Details: err.Error(),
		})
		return
	}

	// Valores por defecto
	hilosPorFase := req.HilosPorFase
	if hilosPorFase == 0 {
		hilosPorFase = 1
	}
	numTuberias := req.NumTuberias
	if numTuberias == 0 {
		numTuberias = 1
	}

	// Ejecutar use case
	resultado, err := h.ajustarCorrienteUseCase.Execute(
		c.Request.Context(),
		corrienteNominal,
		req.Estado,
		tipoCanalizacion,
		sistemaElectrico,
		tipoEquipo,
		hilosPorFase,
		numTuberias,
	)
	if err != nil {
		status, response := h.mapCorrienteAjustadaErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CorrienteAjustadaResponse{
		Success: true,
		Data:    resultado,
	})
}

// mapCorrienteAjustadaErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *CalculoHandler) mapCorrienteAjustadaErrorToResponse(err error) (int, CorrienteAjustadaResponseError) {
	// Errores 400 - Bad Request
	if errors.Is(err, entity.ErrTipoEquipoInvalido) {
		return http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Tipo de equipo inválido",
			Code:    "TIPO_EQUIPO_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoCanalizacionInvalido) {
		return http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Tipo de canalización inválido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrSistemaElectricoInvalido) {
		return http.StatusBadRequest, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "Sistema eléctrico inválido",
			Code:    "SISTEMA_ELECTRICO_INVALIDO",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity
	if errors.Is(err, dto.ErrConductorNoEncontrado) ||
		errors.Is(err, dto.ErrCanalizacionNoDisponible) {
		return http.StatusUnprocessableEntity, CorrienteAjustadaResponseError{
			Success: false,
			Error:   "No se pudo calcular el ajuste de corriente",
			Code:    "CALCULO_NO_POSIBLE",
			Details: err.Error(),
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, CorrienteAjustadaResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
