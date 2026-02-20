// internal/calculos/infrastructure/adapter/driver/http/caida_tension_handler.go
package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/gin-gonic/gin"
)

// CaidaTensionHandler maneja el endpoint de cálculo de caída de tensión.
type CaidaTensionHandler struct {
	calcularCaidaTensionUC *usecase.CalcularCaidaTensionUseCase
}

// NewCaidaTensionHandler crea un nuevo handler de caída de tensión.
func NewCaidaTensionHandler(
	calcularCaidaTensionUC *usecase.CalcularCaidaTensionUseCase,
) *CaidaTensionHandler {
	return &CaidaTensionHandler{
		calcularCaidaTensionUC: calcularCaidaTensionUC,
	}
}

// ============================================
// Endpoint: Calcular Caída de Tensión
// ============================================

// CaidaTensionRequest representa el body de la petición POST /caida-tension.
type CaidaTensionRequest struct {
	// Calibre es el calibre del conductor (ej: "2 AWG", "250 MCM")
	Calibre string `json:"calibre" binding:"required"`

	// Material es el material del conductor ("Cu" = cobre, "Al" = aluminio)
	Material string `json:"material" binding:"required"`

	// TipoCanalizacion es el tipo de canalización
	// Valores: "TUBERIA_PVC", "TUBERIA_METALICA", "CHAROLA"
	TipoCanalizacion string `json:"tipo_canalizacion" binding:"required"`

	// CorrienteAjustada es la corriente ajustada en amperes
	CorrienteAjustada float64 `json:"corriente_ajustada" binding:"required,gt=0"`

	// LongitudCircuito es la longitud del circuito en metros
	LongitudCircuito float64 `json:"longitud_circuito" binding:"required,gt=0"`

	// Tension es la tensión del sistema en volts
	Tension int `json:"tension" binding:"required,gt=0"`

	// SistemaElectrico es el tipo de sistema eléctrico
	// Valores: "MONOFASICO", "BIFASICO", "DELTA", "ESTRELLA"
	SistemaElectrico string `json:"sistema_electrico" binding:"required"`

	// TipoVoltaje indica si el voltaje ingresado es fase-neutro o fase-fase
	// Valores: "FASE_NEUTRO" (Vfn), "FASE_FASE" (Vff), también acepta "FN" o "FF"
	// Ejemplos: 127V es típicamente Vfn, 220V es típicamente Vff
	TipoVoltaje string `json:"tipo_voltaje" binding:"required"`

	// HilosPorFase es el número de hilos por fase
	HilosPorFase int `json:"hilos_por_fase" binding:"required,min=1"`

	// LimiteCaida es el límite de caída de tensión en porcentaje
	LimiteCaida float64 `json:"limite_caida" binding:"required,gt=0"`
}

// CaidaTensionResponse representa la respuesta exitosa.
type CaidaTensionResponse struct {
	Success bool                      `json:"success"`
	Data    dto.ResultadoCaidaTension `json:"data"`
}

// CaidaTensionResponseError representa la respuesta de error.
type CaidaTensionResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// CalcularCaidaTension POST /api/v1/calculos/caida-tension
// Calcula la caída de tensión en un circuito según IEEE-141.
func (h *CaidaTensionHandler) CalcularCaidaTension(c *gin.Context) {
	var req CaidaTensionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
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
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Tipo de canalización inválido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Parsear material del conductor
	material, err := valueobject.ParseMaterialConductor(req.Material)
	if err != nil {
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Material del conductor inválido",
			Code:    "MATERIAL_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Parsear sistema eléctrico
	sistemaElectrico, err := entity.ParseSistemaElectrico(req.SistemaElectrico)
	if err != nil {
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Sistema eléctrico inválido",
			Code:    "SISTEMA_ELECTRICO_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Parsear tipo de voltaje
	tipoVoltaje, err := entity.ParseTipoVoltaje(req.TipoVoltaje)
	if err != nil {
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Tipo de voltaje inválido",
			Code:    "TIPO_VOLTAJE_INVALIDO",
			Details: err.Error(),
		})
		return
	}

	// Crear valor de tensión
	tension, err := valueobject.NewTension(req.Tension)
	if err != nil {
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Tensión inválida",
			Code:    "TENSION_INVALIDA",
			Details: err.Error(),
		})
		return
	}

	// Crear valor de corriente ajustada
	corrienteAjustada, err := valueobject.NewCorriente(req.CorrienteAjustada)
	if err != nil {
		c.JSON(http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Corriente ajustada inválida",
			Code:    "CORRIENTE_INVALIDA",
			Details: err.Error(),
		})
		return
	}

	// Valor por defecto de hilos por fase
	hilosPorFase := req.HilosPorFase
	if hilosPorFase == 0 {
		hilosPorFase = 1
	}

	// Ejecutar use case - el use case se encarga de obtener los datos del conductor
	resultado, err := h.calcularCaidaTensionUC.Execute(
		c.Request.Context(),
		req.Calibre,
		material,
		corrienteAjustada,
		req.LongitudCircuito,
		tension,
		req.LimiteCaida,
		tipoCanalizacion,
		sistemaElectrico,
		tipoVoltaje,
		hilosPorFase,
	)
	if err != nil {
		status, response := h.mapCaidaTensionErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CaidaTensionResponse{
		Success: true,
		Data:    resultado,
	})
}

// mapCaidaTensionErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *CaidaTensionHandler) mapCaidaTensionErrorToResponse(err error) (int, CaidaTensionResponseError) {
	// Errores 400 - Bad Request
	if errors.Is(err, entity.ErrTipoCanalizacionInvalido) {
		return http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Tipo de canalización inválido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrMaterialConductorInvalido) {
		return http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Material del conductor inválido",
			Code:    "MATERIAL_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrVoltajeInvalido) {
		return http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Tensión fuera de rango permitido",
			Code:    "TENSION_INVALIDA",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoVoltajeInvalido) {
		return http.StatusBadRequest, CaidaTensionResponseError{
			Success: false,
			Error:   "Tipo de voltaje inválido",
			Code:    "TIPO_VOLTAJE_INVALIDO",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity
	// El error de impedancia no encontrada viene del CSV repository como fmt.Errorf
	errStr := err.Error()
	if strings.Contains(errStr, "not found") || strings.Contains(errStr, "no encontrado") {
		return http.StatusUnprocessableEntity, CaidaTensionResponseError{
			Success: false,
			Error:   "No se encontró la impedancia para el calibre y canalización",
			Code:    "IMPEDANCIA_NO_ENCONTRADA",
			Details: errStr,
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, CaidaTensionResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: errStr,
	}
}
