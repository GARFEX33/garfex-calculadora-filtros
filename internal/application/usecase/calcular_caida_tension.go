// internal/application/usecase/calcular_caida_tension.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ResultadoCaidaTensionUseCase contiene el resultado del cálculo de caída.
type ResultadoCaidaTensionUseCase struct {
	Porcentaje          float64
	CaidaVolts          float64
	Cumple              bool
	LimitePorcentaje    float64
	ResistenciaEfectiva float64
}

// CalcularCaidaTensionUseCase ejecuta el Paso 7: Calcular Caída de Tensión.
type CalcularCaidaTensionUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewCalcularCaidaTensionUseCase crea una nueva instancia.
func NewCalcularCaidaTensionUseCase(
	tablaRepo port.TablaNOMRepository,
) *CalcularCaidaTensionUseCase {
	return &CalcularCaidaTensionUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute calcula la caída de tensión.
func (uc *CalcularCaidaTensionUseCase) Execute(
	ctx context.Context,
	conductor valueobject.Conductor,
	corrienteAjustada valueobject.Corriente,
	longitudCircuito float64,
	tension valueobject.Tension,
	limiteCaida float64,
	tipoCanalizacion entity.TipoCanalizacion,
	factorPotencia float64,
	hilosPorFase int,
) (ResultadoCaidaTensionUseCase, error) {
	// Convertir material string a MaterialConductor
	material := valueobject.MaterialCobre
	if conductor.Material() == "Al" || conductor.Material() == "AL" {
		material = valueobject.MaterialAluminio
	}

	// Obtener impedancia
	impedancia, err := uc.tablaRepo.ObtenerImpedancia(ctx, conductor.Calibre(), tipoCanalizacion, material)
	if err != nil {
		return ResultadoCaidaTensionUseCase{}, fmt.Errorf("obtener impedancia: %w", err)
	}

	entradaCaida := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: impedancia.R,
		ReactanciaOhmPorKm:  impedancia.X,
		TipoCanalizacion:    tipoCanalizacion,
		HilosPorFase:        hilosPorFase,
		FactorPotencia:      factorPotencia,
	}

	resultadoCaida, err := service.CalcularCaidaTension(
		entradaCaida,
		corrienteAjustada,
		longitudCircuito,
		tension,
		limiteCaida,
	)
	if err != nil {
		return ResultadoCaidaTensionUseCase{}, fmt.Errorf("calcular caída de tensión: %w", err)
	}

	return ResultadoCaidaTensionUseCase{
		Porcentaje:          resultadoCaida.Porcentaje,
		CaidaVolts:          resultadoCaida.CaidaVolts,
		Cumple:              resultadoCaida.Cumple,
		LimitePorcentaje:    limiteCaida,
		ResistenciaEfectiva: resultadoCaida.Impedancia,
	}, nil
}
