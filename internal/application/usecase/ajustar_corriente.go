// internal/application/usecase/ajustar_corriente.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ResultadoAjusteCorriente contiene el resultado del ajuste de corriente.
type ResultadoAjusteCorriente struct {
	CorrienteAjustada  valueobject.Corriente
	FactorTemperatura  float64
	FactorAgrupamiento float64
	FactorTotal        float64
	Temperatura        valueobject.Temperatura
}

// AjustarCorrienteUseCase ejecuta el Paso 2: Ajuste de Corriente.
type AjustarCorrienteUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewAjustarCorrienteUseCase crea una nueva instancia.
func NewAjustarCorrienteUseCase(
	tablaRepo port.TablaNOMRepository,
) *AjustarCorrienteUseCase {
	return &AjustarCorrienteUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute aplica los factores de corrección a la corriente nominal.
func (uc *AjustarCorrienteUseCase) Execute(
	ctx context.Context,
	corrienteNominal valueobject.Corriente,
	estado string,
	tipoCanalizacion entity.TipoCanalizacion,
	sistemaElectrico entity.SistemaElectrico,
	temperaturaOverride *valueobject.Temperatura,
) (ResultadoAjusteCorriente, error) {
	// Obtener temperatura ambiente
	tempAmbiente, err := uc.tablaRepo.ObtenerTemperaturaPorEstado(ctx, estado)
	if err != nil {
		return ResultadoAjusteCorriente{}, fmt.Errorf("obtener temperatura: %w", err)
	}

	// Seleccionar temperatura según reglas NOM
	temperatura := uc.seleccionarTemperatura(corrienteNominal, tipoCanalizacion, temperaturaOverride)

	// Calcular factor temperatura usando el puerto
	factorTemp, err := uc.tablaRepo.ObtenerFactorTemperatura(ctx, tempAmbiente, temperatura)
	if err != nil {
		return ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor temperatura: %w", err)
	}

	// Obtener cantidad de conductors según sistema eléctrico
	cantidadConductores := sistemaElectrico.CantidadConductores()

	// Calcular factor agrupamiento usando el puerto
	factorAgr, err := uc.tablaRepo.ObtenerFactorAgrupamiento(ctx, cantidadConductores)
	if err != nil {
		return ResultadoAjusteCorriente{}, fmt.Errorf("calcular factor agrupamiento: %w", err)
	}

	// Aplicar ajuste
	factores := map[string]float64{
		"agrupamiento": factorAgr,
		"temperatura":  factorTemp,
	}

	corrienteAjustada, err := service.AjustarCorriente(corrienteNominal, factores)
	if err != nil {
		return ResultadoAjusteCorriente{}, fmt.Errorf("ajustar corriente: %w", err)
	}

	return ResultadoAjusteCorriente{
		CorrienteAjustada:  corrienteAjustada,
		FactorTemperatura:  factorTemp,
		FactorAgrupamiento: factorAgr,
		FactorTotal:        factorAgr * factorTemp,
		Temperatura:        temperatura,
	}, nil
}

// seleccionarTemperatura determina la temperatura según las reglas del AGENTS.md.
func (uc *AjustarCorrienteUseCase) seleccionarTemperatura(
	corriente valueobject.Corriente,
	tipoCanalizacion entity.TipoCanalizacion,
	override *valueobject.Temperatura,
) valueobject.Temperatura {
	// Si hay override explícito, usarlo
	if override != nil {
		return *override
	}

	// Reglas según AGENTS.md de application
	if corriente.Valor() <= 100 {
		// <= 100A -> 60C (o 75C si charola triangular sin columna 60C)
		if tipoCanalizacion == entity.TipoCanalizacionCharolaCableTriangular {
			return valueobject.Temp75
		}
		return valueobject.Temp60
	}

	// > 100A -> 75C
	return valueobject.Temp75
}
