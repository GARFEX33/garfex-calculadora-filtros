// internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go
package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/gin-gonic/gin"
)

// CalculoHandler maneja los endpoints de cálculo.
type CalculoHandler struct {
	calcularMemoriaUseCase   *usecase.CalcularMemoriaUseCase
	calcularCorrienteUseCase *usecase.CalcularCorrienteUseCase
	ajustarCorrienteUseCase  *usecase.AjustarCorrienteUseCase
}

// NewCalculoHandler crea un nuevo handler de cálculo.
func NewCalculoHandler(
	calcularMemoriaUC *usecase.CalcularMemoriaUseCase,
	calcularCorrienteUC *usecase.CalcularCorrienteUseCase,
	ajustarCorrienteUC *usecase.AjustarCorrienteUseCase,
) *CalculoHandler {
	return &CalculoHandler{
		calcularMemoriaUseCase:   calcularMemoriaUC,
		calcularCorrienteUseCase: calcularCorrienteUC,
		ajustarCorrienteUseCase:  ajustarCorrienteUC,
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
	Material           string  `json:"material,omitempty"` // "Cu" o "Al"; default: Cu
	Estado             string  `json:"estado" binding:"required"`
	SistemaElectrico   string  `json:"sistema_electrico" binding:"required"`
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

	// Tipo de canalización (string para el DTO)
	tipoCanalizacion := req.TipoCanalizacion

	// Tipo de equipo (string para el DTO)
	tipoEquipo := req.TipoEquipo

	// Temperatura override
	var tempOverride *valueobject.Temperatura
	if req.Temperatura != nil {
		temp := valueobject.Temperatura(*req.Temperatura)
		tempOverride = &temp
	}

	// Sistema eléctrico (DTO string)
	sistemaElectrico := dto.SistemaElectrico(req.SistemaElectrico)

	// Validar sistema eléctrico
	if err := entity.ValidarSistemaElectrico(sistemaElectrico.ToEntity()); err != nil {
		return dto.EquipoInput{}, fmt.Errorf("sistema eléctrico inválido: %w", err)
	}

	// Parsear material desde string usando ParseMaterialConductor del dominio
	material := valueobject.MaterialCobre
	if req.Material != "" {
		var err error
		material, err = valueobject.ParseMaterialConductor(req.Material)
		if err != nil {
			return dto.EquipoInput{}, fmt.Errorf("material inválido: %w", err)
		}
	}

	input := dto.EquipoInput{
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
		Material:              material,
		Estado:                req.Estado,
		SistemaElectrico:      sistemaElectrico,
	}

	return input, nil
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
	if errors.Is(err, dto.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "No se encontró conductor adecuado",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, dto.ErrCanalizacionNoDisponible) {
		return http.StatusUnprocessableEntity, CalcularMemoriaResponseError{
			Success: false,
			Error:   "No se encontró canalización adecuada",
			Code:    "CANALIZACION_NO_DISPONIBLE",
			Details: err.Error(),
		}
	}

	// Errores específicos de validación en cálculos
	if errors.Is(err, dto.ErrDistanciaInvalida) ||
		errors.Is(err, dto.ErrHilosPorFaseInvalido) ||
		errors.Is(err, dto.ErrFactorPotenciaInvalido) {
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

	// Crear tensión usando el constructor
	tension, err := valueobject.NewTension(req.Tension)
	if err != nil {
		c.JSON(http.StatusBadRequest, CalcularAmperajeResponseError{
			Success: false,
			Error:   "Tensión inválida",
			Code:    "TENSION_INVALIDA",
			Details: err.Error(),
		})
		return
	}

	// Convertir request a EquipoInput con modo MANUAL_POTENCIA
	input := dto.EquipoInput{
		Modo:             dto.ModoManualPotencia,
		PotenciaNominal:  req.PotenciaWatts,
		Tension:          tension,
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
