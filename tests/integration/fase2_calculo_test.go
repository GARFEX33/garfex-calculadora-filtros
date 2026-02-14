package integration

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/garfex/calculadora-filtros/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubEquipoRepository struct{}

func (s *stubEquipoRepository) BuscarPorClave(ctx context.Context, clave string) (entity.CalculadorCorriente, error) {
	return nil, nil
}

func TestFase2_CalculoCompleto(t *testing.T) {
	tablaRepo, err := repository.NewCSVTablaNOMRepository("../../data/tablas_nom")
	require.NoError(t, err)

	equipoRepo := &stubEquipoRepository{}
	uc := usecase.NewCalcularMemoriaUseCase(tablaRepo, equipoRepo)

	ctx := context.Background()

	tension, err := valueobject.NewTension(480)
	require.NoError(t, err)

	input := dto.EquipoInput{
		Modo:             dto.ModoManualAmperaje,
		TipoEquipo:       "FILTRO_ACTIVO",
		Clave:            "FA-TEST-001",
		AmperajeNominal:  100,
		Tension:          tension,
		FactorPotencia:   0.9,
		Estado:           "Nuevo Leon",
		SistemaElectrico: dto.SistemaElectricoDelta,
		TipoCanalizacion: "TUBERIA_PVC",
		ITM:              125,
		LongitudCircuito: 50,
		HilosPorFase:     1,
	}

	output, err := uc.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Nuevo Leon", output.Estado)
	assert.Equal(t, 37, output.TemperaturaAmbiente) // 36.8°C → round → 37°C (temp máxima 2022)
	assert.Equal(t, dto.SistemaElectricoDelta, output.SistemaElectrico)
	assert.Equal(t, 3, output.CantidadConductores)
	assert.InDelta(t, 0.94, output.FactorTemperaturaCalculado, 0.01)  // rango 36-40°C, conductor 60C (100A) → 0.94
	assert.InDelta(t, 0.70, output.FactorAgrupamientoCalculado, 0.01) // 3 conductores → 0.70
	assert.NotEmpty(t, output.Canalizacion.Tamano)
}

var _ port.EquipoRepository = (*stubEquipoRepository)(nil)
