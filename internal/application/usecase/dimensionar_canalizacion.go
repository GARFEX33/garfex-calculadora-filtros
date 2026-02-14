// internal/application/usecase/dimensionar_canalizacion.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
)

// DimensionarCanalizacionUseCase ejecuta el Paso 6: Dimensionar Canalización.
type DimensionarCanalizacionUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewDimensionarCanalizacionUseCase crea una nueva instancia.
func NewDimensionarCanalizacionUseCase(
	tablaRepo port.TablaNOMRepository,
) *DimensionarCanalizacionUseCase {
	return &DimensionarCanalizacionUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute dimensiona la canalización.
func (uc *DimensionarCanalizacionUseCase) Execute(
	ctx context.Context,
	conductorAlimentacionSeccionMM2 float64,
	conductorTierraSeccionMM2 float64,
	hilosPorFase int,
	tipoCanalizacion entity.TipoCanalizacion,
) (dto.ResultadoCanalizacion, error) {
	tablaCanalizacion, err := uc.tablaRepo.ObtenerTablaCanalizacion(ctx, tipoCanalizacion)
	if err != nil {
		return dto.ResultadoCanalizacion{}, fmt.Errorf("obtener tabla canalización: %w", err)
	}

	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: hilosPorFase * 3, SeccionMM2: conductorAlimentacionSeccionMM2}, // Fases
		{Cantidad: 1, SeccionMM2: conductorTierraSeccionMM2},                      // Tierra
	}

	resultado, err := service.CalcularCanalizacion(conductores, string(tipoCanalizacion), tablaCanalizacion, hilosPorFase)
	if err != nil {
		return dto.ResultadoCanalizacion{}, fmt.Errorf("calcular canalización: %w", err)
	}

	// Map domain entity to DTO with primitive types (no domain objects exposed)
	return dto.ResultadoCanalizacion{
		Tamano:           resultado.Tamano,
		AreaTotalMM2:     resultado.AnchoRequerido,
		AreaRequeridaMM2: resultado.AnchoRequerido,
		NumeroDeTubos:    resultado.NumeroDeTubos,
	}, nil
}
