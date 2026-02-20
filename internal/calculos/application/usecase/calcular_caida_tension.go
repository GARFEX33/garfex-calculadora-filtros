// internal/calculos/application/usecase/calcular_caida_tension.go
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
	calibre string,
	material valueobject.MaterialConductor,
	corrienteAjustada valueobject.Corriente,
	longitudCircuito float64,
	tension valueobject.Tension,
	limiteCaida float64,
	tipoCanalizacion entity.TipoCanalizacion,
	sistemaElectrico entity.SistemaElectrico,
	hilosPorFase int,
) (dto.ResultadoCaidaTension, error) {
	// Obtener impedancia
	impedancia, err := uc.tablaRepo.ObtenerImpedancia(ctx, calibre, tipoCanalizacion, material)
	if err != nil {
		return dto.ResultadoCaidaTension{}, fmt.Errorf("obtener impedancia: %w", err)
	}

	entradaCaida := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: impedancia.R(),
		ReactanciaOhmPorKm:  impedancia.X(),
		TipoCanalizacion:    tipoCanalizacion,
		SistemaElectrico:    sistemaElectrico,
		HilosPorFase:        hilosPorFase,
	}

	resultadoCaida, err := service.CalcularCaidaTension(
		entradaCaida,
		corrienteAjustada,
		longitudCircuito,
		tension,
		limiteCaida,
	)
	if err != nil {
		return dto.ResultadoCaidaTension{}, fmt.Errorf("calcular caída de tensión: %w", err)
	}

	return dto.ResultadoCaidaTension{
		Porcentaje:          resultadoCaida.Porcentaje,
		CaidaVolts:          resultadoCaida.CaidaVolts,
		Cumple:              resultadoCaida.Cumple,
		LimitePorcentaje:    limiteCaida,
		ResistenciaEfectiva: resultadoCaida.Impedancia,
	}, nil
}
