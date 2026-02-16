package integration

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driven/csv"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubEquipoRepository struct{}

func (s *stubEquipoRepository) BuscarPorClave(ctx context.Context, clave string) (entity.CalculadorCorriente, error) {
	return nil, nil
}

var _ port.EquipoRepository = (*stubEquipoRepository)(nil)

func TestFase2_CalculoCompleto(t *testing.T) {
	tablaRepo, err := csv.NewCSVTablaNOMRepository("../../data/tablas_nom")
	require.NoError(t, err)

	equipoRepo := &stubEquipoRepository{}

	// Create micro use cases
	calcularCorrienteUC := usecase.NewCalcularCorrienteUseCase(equipoRepo)
	ajustarCorrienteUC := usecase.NewAjustarCorrienteUseCase(tablaRepo)
	seleccionarConductorUC := usecase.NewSeleccionarConductorUseCase(tablaRepo)
	dimensionarCanalizacionUC := usecase.NewDimensionarCanalizacionUseCase(tablaRepo)
	calcularCaidaTensionUC := usecase.NewCalcularCaidaTensionUseCase(tablaRepo)

	// Create orquestador
	orquestador := usecase.NewOrquestadorMemoriaCalculo(
		calcularCorrienteUC,
		ajustarCorrienteUC,
		seleccionarConductorUC,
		dimensionarCanalizacionUC,
		calcularCaidaTensionUC,
		tablaRepo,
	)

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

	output, err := orquestador.Execute(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "Nuevo Leon", output.Estado)
	assert.Equal(t, 37, output.TemperaturaAmbiente) // 36.8°C → round → 37°C (temp máxima 2022)
	assert.Equal(t, dto.SistemaElectricoDelta, output.SistemaElectrico)
	assert.Equal(t, 3, output.CantidadConductores)
	// Factor temp: 60°C conductor at 37°C ambient → 0.82 (NOM tables)
	// Factor agrupamiento: 1.0 (charola returns 1.0, tuberia with 3 conductors would be 0.70)
	// Note: actual values depend on orchestator logic and temperature selection
	assert.NotEmpty(t, output.Canalizacion.Tamano)
}
