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
	tablaRepo port.TablaNOMRepository
}

// NewAjustarCorrienteUseCase creates a new instance.
func NewAjustarCorrienteUseCase(
	tablaRepo port.TablaNOMRepository,
) *AjustarCorrienteUseCase {
	return &AjustarCorrienteUseCase{
		tablaRepo: tablaRepo,
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

	// Select temperature using domain service (pure logic, no I/O)
	// No override for this use case
	temperatura := service.SeleccionarTemperatura(
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

	// Determine if grouping factor applies
	// CHAROLA: no aplica factor de agrupamiento (cables separados o en configuración triangular)
	// TUBERIA: aplica factor de agrupamiento
	esCharola := tipoCanalizacion.EsCharola()

	if !esCharola {
		// Validate that conductors can be evenly distributed (only for tuberia)
		if cantidadTotal%numTuberias != 0 {
			return dto.ResultadoAjusteCorriente{}, fmt.Errorf(
				"cantidad total de conductores (%d) no es divisible por número de tuberías (%d)",
				cantidadTotal, numTuberias,
			)
		}
	}

	conductoresPorTubo := cantidadTotal / numTuberias

	// Get grouping factor - ONLY for tuberia, not for charola
	var factorAgr float64
	if esCharola {
		// Charola: no aplica factor de agrupamiento
		factorAgr = 1.0
	} else {
		// Tubería: aplica factor de agrupamiento basado en conductores por tubo
		factorAgr, err = uc.tablaRepo.ObtenerFactorAgrupamiento(ctx, conductoresPorTubo)
		if err != nil {
			return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor agrupamiento: %w", err)
		}
	}

	// Get usage factor based on equipment type
	factorUso, err := service.CalcularFactorUso(tipoEquipo)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor uso: %w", err)
	}

	// Calculate adjusted current using domain service
	// Fórmula: corrienteAjustada = corrienteNominal * factorUso / (factorTemp * factorAgr)
	// El servicio multiplica todos los factores, entonces pasamos las inversas
	// para temperatura y agrupamiento
	factores := map[string]float64{
		"uso":          factorUso,
		"temperatura":  1.0 / factorTemp,
		"agrupamiento": 1.0 / factorAgr,
	}

	resultado, err := service.AjustarCorriente(corrienteNominal, factores)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("ajustar corriente: %w", err)
	}

	// Return DTO with primitive types (no domain objects exposed)
	return dto.ResultadoAjusteCorriente{
		CorrienteAjustada:        resultado.CorrienteAjustada.Valor(),
		FactorTemperatura:        factorTemp,
		FactorAgrupamiento:       factorAgr,
		FactorUso:                factorUso,
		FactorTotal:              resultado.FactorTotal,
		Temperatura:              temperatura.Valor(),
		ConductoresPorTubo:       conductoresPorTubo,
		CantidadConductoresTotal: cantidadTotal,
		TemperaturaAmbiente:      tempAmbiente,
	}, nil
}
