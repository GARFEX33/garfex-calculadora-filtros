// internal/presentation/handler/calculo_handler_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock implementations
type mockTablaRepo struct{}

func (m *mockTablaRepo) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return []valueobject.EntradaTablaTierra{
		{ITMHasta: 100, Conductor: valueobject.ConductorParams{Calibre: "8 AWG", Material: "Cu", SeccionMM2: 8.37}},
	}, nil
}

func (m *mockTablaRepo) ObtenerTablaAmpacidad(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura) ([]valueobject.EntradaTablaConductor, error) {
	return []valueobject.EntradaTablaConductor{
		{Capacidad: 30, Conductor: valueobject.ConductorParams{Calibre: "10 AWG", Material: "Cu", SeccionMM2: 5.26}},
		{Capacidad: 55, Conductor: valueobject.ConductorParams{Calibre: "8 AWG", Material: "Cu", SeccionMM2: 8.37}},
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
	handler := NewCalculoHandler(calcularMemoriaUC)

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
		assert.Equal(t, 50.0, response.Data.CorrienteNominal.Valor())
	}
}

func TestCalculoHandler_CalcularMemoria_ValidationError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	tablaRepo := &mockTablaRepo{}
	equipoRepo := &mockEquipoRepo{}
	calcularMemoriaUC := usecase.NewCalcularMemoriaUseCase(tablaRepo, equipoRepo)
	handler := NewCalculoHandler(calcularMemoriaUC)

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
