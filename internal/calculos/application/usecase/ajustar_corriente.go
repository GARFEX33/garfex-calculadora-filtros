// internal/calculos/application/usecase/ajustar_corriente.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// AjustarCorrienteUseCase executes Step 2: Current Adjustment.
type AjustarCorrienteUseCase struct {
	tablaRepo           port.TablaNOMRepository
	seleccionarTempPort port.SeleccionarTemperaturaPort
}

// NewAjustarCorrienteUseCase creates a new instance.
func NewAjustarCorrienteUseCase(
	tablaRepo port.TablaNOMRepository,
	seleccionarTempPort port.SeleccionarTemperaturaPort,
) *AjustarCorrienteUseCase {
	return &AjustarCorrienteUseCase{
		tablaRepo:           tablaRepo,
		seleccionarTempPort: seleccionarTempPort,
	}
}

// Execute applies correction factors to the nominal current.
// Accepts tipoEquipo for usage factor, hilosPorFase and numTuberias for grouping calculation.
func (uc *AjustarCorrienteUseCase) Execute(
	ctx context.Context,
	corrienteNominal valueobject.Corriente,
	estado string,
	tipoCanalizacion entity.TipoCanalizacion,
	sistemaElectrico entity.SistemaElectrico,
	tipoEquipo entity.TipoEquipo,
	hilosPorFase int,
	numTuberias int,
) (dto.ResultadoAjusteCorriente, error) {
	// Validate inputs
	if hilosPorFase < 1 {
		hilosPorFase = 1
	}
	if numTuberias < 1 {
		numTuberias = 1
	}

	// Get ambient temperature
	tempAmbiente, err := uc.tablaRepo.ObtenerTemperaturaPorEstado(ctx, estado)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("obtener temperatura: %w", err)
	}

	// Select temperature using the port (delegates to domain service)
	// No override for this use case
	temperatura := uc.seleccionarTempPort.SeleccionarTemperatura(
		corrienteNominal,
		tipoCanalizacion,
		nil,
	)

	// Get temperature factor from repository
	factorTemp, err := uc.tablaRepo.ObtenerFactorTemperatura(ctx, tempAmbiente, temperatura)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor temperatura: %w", err)
	}

	// Calculate conductor distribution
	// fases = number of energized conductors based on electrical system
	fases := sistemaElectrico.CantidadConductores()
	cantidadTotal := fases * hilosPorFase

	// Validate that conductors can be evenly distributed
	if cantidadTotal%numTuberias != 0 {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf(
			"cantidad total de conductores (%d) no es divisible por número de tuberías (%d)",
			cantidadTotal, numTuberias,
		)
	}

	conductoresPorTubo := cantidadTotal / numTuberias

	// Get grouping factor based on conductors per tube (not total)
	factorAgr, err := uc.tablaRepo.ObtenerFactorAgrupamiento(ctx, conductoresPorTubo)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor agrupamiento: %w", err)
	}

	// Get usage factor based on equipment type
	factorUso, err := service.CalcularFactorUso(tipoEquipo)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor uso: %w", err)
	}

	// Calculate total correction factor
	factorTotal := factorTemp * factorAgr * factorUso

	// Calculate adjusted current
	corrienteAjustada := corrienteNominal.Valor() * factorTotal

	// Return DTO with primitive types (no domain objects exposed)
	return dto.ResultadoAjusteCorriente{
		CorrienteAjustada:        corrienteAjustada,
		FactorTemperatura:        factorTemp,
		FactorAgrupamiento:       factorAgr,
		FactorUso:                factorUso,
		FactorTotal:              factorTotal,
		Temperatura:              temperatura.Valor(),
		ConductoresPorTubo:       conductoresPorTubo,
		CantidadConductoresTotal: cantidadTotal,
		TemperaturaAmbiente:      tempAmbiente,
	}, nil
}
