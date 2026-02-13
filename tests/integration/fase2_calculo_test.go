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
		TipoEquipo:       entity.TipoEquipoFiltroActivo,
		Clave:            "FA-TEST-001",
		AmperajeNominal:  100,
		Tension:          tension,
		FactorPotencia:   0.9,
		Estado:           "Nuevo Leon",
		SistemaElectrico: entity.SistemaElectricoDelta,
		TipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
		ITM:              125,
		LongitudCircuito: 50,
		HilosPorFase:     1,
	}

	output, err := uc.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Nuevo Leon", output.Estado)
	assert.Equal(t, 21, output.TemperaturaAmbiente)
	assert.Equal(t, entity.SistemaElectricoDelta, output.SistemaElectrico)
	assert.Equal(t, 3, output.CantidadConductores)
	assert.InDelta(t, 1.1, output.FactorTemperaturaCalculado, 0.01)
	assert.InDelta(t, 0.3, output.FactorAgrupamientoCalculado, 0.01)
	assert.NotEmpty(t, output.Canalizacion.Tamano)
}

var _ port.EquipoRepository = (*stubEquipoRepository)(nil)
