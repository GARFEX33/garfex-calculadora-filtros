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
func (uc *AjustarCorrienteUseCase) Execute(
	ctx context.Context,
	corrienteNominal valueobject.Corriente,
	estado string,
	tipoCanalizacion entity.TipoCanalizacion,
	sistemaElectrico entity.SistemaElectrico,
	temperaturaOverride *valueobject.Temperatura,
) (dto.ResultadoAjusteCorriente, error) {
	// Get ambient temperature
	tempAmbiente, err := uc.tablaRepo.ObtenerTemperaturaPorEstado(ctx, estado)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("obtener temperatura: %w", err)
	}

	// Select temperature using the port (delegates to domain service)
	temperatura := uc.seleccionarTempPort.SeleccionarTemperatura(
		corrienteNominal,
		tipoCanalizacion,
		temperaturaOverride,
	)

	// Get temperature factor from repository
	factorTemp, err := uc.tablaRepo.ObtenerFactorTemperatura(ctx, tempAmbiente, temperatura)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor temperatura: %w", err)
	}

	// Get conductors count from electrical system
	cantidadConductores := sistemaElectrico.CantidadConductores()

	// Get grouping factor from repository
	factorAgr, err := uc.tablaRepo.ObtenerFactorAgrupamiento(ctx, cantidadConductores)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor agrupamiento: %w", err)
	}

	// Apply adjustment using domain service (application CAN depend on domain)
	factores := map[string]float64{
		"agrupamiento": factorAgr,
		"temperatura":  factorTemp,
	}

	resultadoAjuste, err := service.AjustarCorriente(corrienteNominal, factores)
	if err != nil {
		return dto.ResultadoAjusteCorriente{}, fmt.Errorf("ajustar corriente: %w", err)
	}

	// Return DTO with primitive types (no domain objects exposed)
	return dto.ResultadoAjusteCorriente{
		CorrienteAjustada:  resultadoAjuste.CorrienteAjustada.Valor(),
		FactorTemperatura:  factorTemp,
		FactorAgrupamiento: factorAgr,
		FactorTotal:        resultadoAjuste.FactorTotal,
		Temperatura:        temperatura.Valor(),
	}, nil
}
