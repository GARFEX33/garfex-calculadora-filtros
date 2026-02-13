// internal/presentation/handler/calculo_handler.go
package handler

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/gin-gonic/gin"
)

// CalculoHandler maneja los endpoints de cálculo.
type CalculoHandler struct {
	calcularMemoriaUseCase *usecase.CalcularMemoriaUseCase
}

// NewCalculoHandler crea un nuevo handler de cálculo.
func NewCalculoHandler(calcularMemoriaUseCase *usecase.CalcularMemoriaUseCase) *CalculoHandler {
	return &CalculoHandler{
		calcularMemoriaUseCase: calcularMemoriaUseCase,
	}
}

// CalcularMemoriaRequest representa el body de la petición POST.
type CalcularMemoriaRequest struct {
	Modo               string  `json:"modo" binding:"required,oneof=LISTADO MANUAL_AMPERAJE MANUAL_POTENCIA"`
	Clave              string  `json:"clave,omitempty"`
	TipoEquipo         string  `json:"tipo_equipo,omitempty"`
	AmperajeNominal    float64 `json:"amperaje_nominal,omitempty"`
	PotenciaNominal    float64 `json:"potencia_nominal,omitempty"`
	Tension            float64 `json:"tension" binding:"required,gt=0"`
	FactorPotencia     float64 `json:"factor_potencia,omitempty"`
	ITM                int     `json:"itm" binding:"required,gt=0"`
	TipoCanalizacion   string  `json:"tipo_canalizacion" binding:"required"`
	Temperatura        *int    `json:"temperatura,omitempty"`
	HilosPorFase       int     `json:"hilos_por_fase,omitempty"`
	LongitudCircuito   float64 `json:"longitud_circuito" binding:"required,gt=0"`
	PorcentajeCaidaMax float64 `json:"porcentaje_caida_max,omitempty"`
	FactorAgrupamiento float64 `json:"factor_agrupamiento,omitempty"`
	FactorTemperatura  float64 `json:"factor_temperatura,omitempty"`
}

// CalcularMemoriaResponse representa la respuesta exitosa.
type CalcularMemoriaResponse struct {
	Success bool              `json:"success"`
	Data    dto.MemoriaOutput `json:"data"`
}

// CalcularMemoriaResponseError representa la respuesta de error.
type CalcularMemoriaResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// CalcularMemoria POST /api/v1/calculos/memoria
func (h *CalculoHandler) CalcularMemoria(c *gin.Context) {
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

	// Convertir request a DTO
	input, err := h.mapRequestToDTO(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Datos inválidos",
			Code:    "INVALID_DATA",
			Details: err.Error(),
		})
		return
	}

	// Ejecutar use case
	output, err := h.calcularMemoriaUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CalcularMemoriaResponse{
		Success: true,
		Data:    output,
	})
}

// mapRequestToDTO convierte el request HTTP al DTO de application.
func (h *CalculoHandler) mapRequestToDTO(req CalcularMemoriaRequest) (dto.EquipoInput, error) {
	// Parsear modo
	var modo dto.ModoCalculo
	switch req.Modo {
	case "LISTADO":
		modo = dto.ModoListado
	case "MANUAL_AMPERAJE":
		modo = dto.ModoManualAmperaje
	case "MANUAL_POTENCIA":
		modo = dto.ModoManualPotencia
	default:
		return dto.EquipoInput{}, dto.ErrModoInvalido
	}

	// Crear tensión
	tension, err := valueobject.NewTension(int(req.Tension))
	if err != nil {
		return dto.EquipoInput{}, err
	}

	// Parsear tipo de canalización
	tipoCanalizacion := entity.TipoCanalizacion(req.TipoCanalizacion)

	// Parsear tipo de equipo si aplica
	var tipoEquipo entity.TipoEquipo
	if req.TipoEquipo != "" {
		tipoEquipo = entity.TipoEquipo(req.TipoEquipo)
	}

	// Temperatura override
	var tempOverride *valueobject.Temperatura
	if req.Temperatura != nil {
		temp := valueobject.Temperatura(*req.Temperatura)
		tempOverride = &temp
	}

	return dto.EquipoInput{
		Modo:                  modo,
		Clave:                 req.Clave,
		TipoEquipo:            tipoEquipo,
		AmperajeNominal:       req.AmperajeNominal,
		PotenciaNominal:       req.PotenciaNominal,
		Tension:               tension,
		FactorPotencia:        req.FactorPotencia,
		ITM:                   req.ITM,
		TipoCanalizacion:      tipoCanalizacion,
		TemperaturaOverride:   tempOverride,
		HilosPorFase:          req.HilosPorFase,
		LongitudCircuito:      req.LongitudCircuito,
		PorcentajeCaidaMaximo: req.PorcentajeCaidaMax,
		FactorAgrupamiento:    req.FactorAgrupamiento,
		FactorTemperatura:     req.FactorTemperatura,
	}, nil
}

// mapErrorToResponse mapea errores del dominio a respuestas HTTP.
// Siguiendo la tabla del AGENTS.md:
// - ErrEquipoNoEncontrado → 404
// - ErrModoInvalido → 400
// - ErrCanalizacionNoSoportada → 400
// - Validación de input → 400
// - Error de cálculo (datos insuficientes) → 422
// - ErrConductorNoEncontrado → 422
// - Error interno → 500
func (h *CalculoHandler) mapErrorToResponse(err error) (int, CalcularMemoriaResponseError) {
	// Errores 400 - Bad Request (datos inválidos)
	if errors.Is(err, dto.ErrModoInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Modo de cálculo inválido",
			Code:    "MODO_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Datos de entrada inválidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoCanalizacionInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Tipo de canalización no soportado",
			Code:    "CANALIZACION_NO_SOPORTADA",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoEquipoInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Tipo de equipo inválido",
			Code:    "TIPO_EQUIPO_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrVoltajeInvalido) ||
		errors.Is(err, valueobject.ErrCorrienteInvalida) ||
		errors.Is(err, valueobject.ErrConductorInvalido) {
		return http.StatusBadRequest, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Valor fuera de rango permitido",
			Code:    "VALOR_INVALIDO",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity (errores de cálculo)
	if errors.Is(err, service.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "No se encontró conductor adecuado",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, service.ErrCanalizacionNoDisponible) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "No se encontró canalización adecuada",
			Code:    "CANALIZACION_NO_DISPONIBLE",
			Details: err.Error(),
		}
	}

	// Errores específicos de validación en cálculos
	if errors.Is(err, service.ErrDistanciaInvalida) ||
		errors.Is(err, service.ErrHilosPorFaseInvalido) ||
		errors.Is(err, service.ErrFactorPotenciaInvalido) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "Parámetros de cálculo inválidos",
			Code:    "PARAMETROS_INVALIDOS",
			Details: err.Error(),
		}
	}

	// Error 404 - Not Found
	// Nota: Esto normalmente vendría del repositorio de equipos
	// if errors.Is(err, repository.ErrEquipoNoEncontrado) {
	//     return http.StatusNotFound, CalcularMemoriaResponseError{...}
	// }

	// Por defecto: error interno 500
	return http.StatusInternalServerError, CalcularMemoriaResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
