// internal/calculos/infrastructure/adapter/driver/http/calculo_handler_test.go
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockTablaRepo struct{}

func (m *mockTablaRepo) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return []valueobject.EntradaTablaTierra{
		{ITMHasta: 100, ConductorCu: valueobject.ConductorParams{Calibre: "8 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 8.37}},
	}, nil
}

func (m *mockTablaRepo) ObtenerTablaAmpacidad(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura) ([]valueobject.EntradaTablaConductor, error) {
	return []valueobject.EntradaTablaConductor{
		{Capacidad: 30, Conductor: valueobject.ConductorParams{Calibre: "10 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 5.26}},
		{Capacidad: 55, Conductor: valueobject.ConductorParams{Calibre: "8 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 8.37}},
	}, nil
}

func (m *mockTablaRepo) ObtenerImpedancia(ctx context.Context, calibre string, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor) (valueobject.ResistenciaReactancia, error) {
	return valueobject.ResistenciaReactancia{R: 3.9, X: 0.164}, nil
}

func (m *mockTablaRepo) ObtenerTablaCanalizacion(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	return []valueobject.EntradaTablaCanalizacion{
		{Tamano: "1/2", AreaInteriorMM2: 78},
		{Tamano: "3/4", AreaInteriorMM2: 122},
	}, nil
}

func (m *mockTablaRepo) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	return 25, nil
}

func (m *mockTablaRepo) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	return 1.0, nil
}

func (m *mockTablaRepo) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	return 1.0, nil
}

func (m *mockTablaRepo) ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error) {
	return 3.5, nil
}

func (m *mockTablaRepo) ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error) {
	return valueobject.EntradaTablaCanalizacion{Tamano: "100mm", AreaInteriorMM2: 5000}, nil
}

func (m *mockTablaRepo) ObtenerCapacidadConductor(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura, calibre string) (float64, error) {
	// Simple mock implementation
	return 55, nil
}

type mockEquipoRepo struct{}

func (m *mockEquipoRepo) BuscarPorClave(ctx context.Context, clave string) (entity.CalculadorCorriente, error) {
	return &mockCalculador{corriente: 50}, nil
}

type mockCalculador struct {
	corriente float64
}

func (m *mockCalculador) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	return valueobject.NewCorriente(m.corriente)
}

func TestCalculoHandler_CalcularMemoria_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	tablaRepo := &mockTablaRepo{}
	equipoRepo := &mockEquipoRepo{}
	calcularMemoriaUC := usecase.NewCalcularMemoriaUseCase(tablaRepo, equipoRepo)
	calcularAmperajeUC := usecase.NewCalcularAmperajeNominalUseCase()
	handler := NewCalculoHandler(calcularMemoriaUC, calcularAmperajeUC)

	// Crear request
	reqBody := CalcularMemoriaRequest{
		Modo:               "MANUAL_AMPERAJE",
		AmperajeNominal:    50,
		Tension:            220,
		FactorPotencia:     1.0,
		ITM:                100,
		TipoCanalizacion:   "TUBERIA_PVC",
		LongitudCircuito:   10,
		PorcentajeCaidaMax: 3.0,
		Estado:             "INTERIOR",
		SistemaElectrico:   "MONOFASICO",
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculos/memoria", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Ejecutar
	handler.CalcularMemoria(c)

	// Verificar
	t.Logf("Response body: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	if w.Code == http.StatusOK {
		var response CalcularMemoriaResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, 50.0, response.Data.CorrienteNominal)
	}
}

func TestCalculoHandler_CalcularMemoria_ValidationError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	tablaRepo := &mockTablaRepo{}
	equipoRepo := &mockEquipoRepo{}
	calcularMemoriaUC := usecase.NewCalcularMemoriaUseCase(tablaRepo, equipoRepo)
	calcularAmperajeUC := usecase.NewCalcularAmperajeNominalUseCase()
	handler := NewCalculoHandler(calcularMemoriaUC, calcularAmperajeUC)

	// Crear request inválido (falta modo)
	reqBody := map[string]interface{}{
		"tension": 220,
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculos/memoria", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Ejecutar
	handler.CalcularMemoria(c)

	// Verificar
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculoHandler_mapErrorToResponse(t *testing.T) {
	handler := &CalculoHandler{}

	tests := []struct {
		name         string
		inputErr     error
		wantStatus   int
		wantCode     string
		wantContains string
	}{
		{
			name:         "ErrModoInvalido → 400",
			inputErr:     dto.ErrModoInvalido,
			wantStatus:   http.StatusBadRequest,
			wantCode:     "MODO_INVALIDO",
			wantContains: "Modo de cálculo inválido",
		},
		{
			name:         "ErrConductorNoEncontrado → 422",
			inputErr:     service.ErrConductorNoEncontrado,
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     "CONDUCTOR_NO_ENCONTRADO",
			wantContains: "No se encontró conductor",
		},
		{
			name:         "ErrCanalizacionNoDisponible → 422",
			inputErr:     service.ErrCanalizacionNoDisponible,
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     "CANALIZACION_NO_DISPONIBLE",
			wantContains: "No se encontró canalización",
		},
		{
			name:         "Error genérico → 500",
			inputErr:     errors.New("error desconocido"),
			wantStatus:   http.StatusInternalServerError,
			wantCode:     "INTERNAL_ERROR",
			wantContains: "Error interno del servidor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, response := handler.mapErrorToResponse(tt.inputErr)
			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantCode, response.Code)
			assert.Contains(t, response.Error, tt.wantContains)
			assert.False(t, response.Success)
		})
	}
}
