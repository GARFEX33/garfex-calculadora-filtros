// internal/calculos/infrastructure/adapter/driver/http/charola_handler.go
package http

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/gin-gonic/gin"
)

// CharolaHandler maneja los endpoints de cálculo de charolas.
type CharolaHandler struct {
	calcularEspaciadoUseCase  *usecase.CalcularCharolaEspaciadoUseCase
	calcularTriangularUseCase *usecase.CalcularCharolaTriangularUseCase
}

// NewCharolaHandler crea un nuevo handler de charolas.
func NewCharolaHandler(
	calcularEspaciadoUC *usecase.CalcularCharolaEspaciadoUseCase,
	calcularTriangularUC *usecase.CalcularCharolaTriangularUseCase,
) *CharolaHandler {
	return &CharolaHandler{
		calcularEspaciadoUseCase:  calcularEspaciadoUC,
		calcularTriangularUseCase: calcularTriangularUC,
	}
}

// ============================================
// Endpoint: Charola Espaciado
// ============================================

// CharolaEspaciadoRequest representa el body de la petición POST /charola/espaciado.
type CharolaEspaciadoRequest struct {
	HilosPorFase      int      `json:"hilos_por_fase" binding:"required,gt=0"`
	SistemaElectrico  string   `json:"sistema_electrico" binding:"required"`
	DiametroFaseMM    float64  `json:"diametro_fase_mm" binding:"required,gt=0"`
	DiametroTierraMM  float64  `json:"diametro_tierra_mm" binding:"required,gt=0"`
	DiametroControlMM *float64 `json:"diametro_control_mm,omitempty"`
}

// CharolaEspaciadoResponse representa la respuesta exitosa.
type CharolaEspaciadoResponse struct {
	Success bool                       `json:"success"`
	Data    dto.CharolaEspaciadoOutput `json:"data"`
}

// CharolaEspaciadoResponseError representa la respuesta de error.
type CharolaEspaciadoResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// PostCharolaEspaciado POST /api/v1/calculos/charola/espaciado
func (h *CharolaHandler) PostCharolaEspaciado(c *gin.Context) {
	var req CharolaEspaciadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CharolaEspaciadoResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Convertir request a DTO
	input := dto.CharolaEspaciadoInput{
		HilosPorFase:      req.HilosPorFase,
		SistemaElectrico:  req.SistemaElectrico,
		DiametroFaseMM:    req.DiametroFaseMM,
		DiametroTierraMM:  req.DiametroTierraMM,
		DiametroControlMM: req.DiametroControlMM,
	}

	// Ejecutar use case
	output, err := h.calcularEspaciadoUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapCharolaErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CharolaEspaciadoResponse{
		Success: true,
		Data:    output,
	})
}

// ============================================
// Endpoint: Charola Triangular
// ============================================

// CharolaTriangularRequest representa el body de la petición POST /charola/triangular.
type CharolaTriangularRequest struct {
	HilosPorFase      int      `json:"hilos_por_fase" binding:"required,gt=0"`
	DiametroFaseMM    float64  `json:"diametro_fase_mm" binding:"required,gt=0"`
	DiametroTierraMM  float64  `json:"diametro_tierra_mm" binding:"required,gt=0"`
	DiametroControlMM *float64 `json:"diametro_control_mm,omitempty"`
}

// CharolaTriangularResponse representa la respuesta exitosa.
type CharolaTriangularResponse struct {
	Success bool                        `json:"success"`
	Data    dto.CharolaTriangularOutput `json:"data"`
}

// CharolaTriangularResponseError representa la respuesta de error.
type CharolaTriangularResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// PostCharolaTriangular POST /api/v1/calculos/charola/triangular
func (h *CharolaHandler) PostCharolaTriangular(c *gin.Context) {
	var req CharolaTriangularRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CharolaTriangularResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Convertir request a DTO
	input := dto.CharolaTriangularInput{
		HilosPorFase:      req.HilosPorFase,
		DiametroFaseMM:    req.DiametroFaseMM,
		DiametroTierraMM:  req.DiametroTierraMM,
		DiametroControlMM: req.DiametroControlMM,
	}

	// Ejecutar use case
	output, err := h.calcularTriangularUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapCharolaErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, CharolaTriangularResponse{
		Success: true,
		Data:    output,
	})
}

// mapCharolaErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *CharolaHandler) mapCharolaErrorToResponse(err error) (int, interface{}) {
	// Errores 400 - Bad Request (validación)
	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, CharolaEspaciadoResponseError{
			Success: false,
			Error:   "Datos de entrada inválidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	// Verificar mensajes de error de validación del DTO
	errMsg := err.Error()
	if contains(errMsg, "hilos_por_fase") ||
		contains(errMsg, "diametro_fase_mm") ||
		contains(errMsg, "diametro_tierra_mm") ||
		contains(errMsg, "sistema_electrico") {
		return http.StatusBadRequest, CharolaEspaciadoResponseError{
			Success: false,
			Error:   "Error de validación",
			Code:    "VALIDATION_ERROR",
			Details: errMsg,
		}
	}

	// Errores 422 - Unprocessable Entity
	if contains(errMsg, "no se encontró") || contains(errMsg, "no disponible") {
		return http.StatusUnprocessableEntity, CharolaEspaciadoResponseError{
			Success: false,
			Error:   "No se pudo calcular el tamaño de charola",
			Code:    "CALCULO_NO_POSIBLE",
			Details: errMsg,
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, CharolaEspaciadoResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}

// contains verifica si una cadena contiene el substring (case insensitive).
func contains(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	return len(sLower) >= len(substrLower) && (sLower == substrLower || len(sLower) > 0 && containsHelper(sLower, substrLower))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}
