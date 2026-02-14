// internal/application/usecase/dimensionar_canalizacion.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ResultadoCanalizacionUseCase contiene el resultado del dimensionamiento.
type ResultadoCanalizacionUseCase struct {
	Tamano        string
	NumeroDeTubos int
	AreaTotalMM2  float64
}

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
	conductorAlimentacion valueobject.Conductor,
	conductorTierra valueobject.Conductor,
	hilosPorFase int,
	tipoCanalizacion entity.TipoCanalizacion,
) (ResultadoCanalizacionUseCase, error) {
	tablaCanalizacion, err := uc.tablaRepo.ObtenerTablaCanalizacion(ctx, tipoCanalizacion)
	if err != nil {
		return ResultadoCanalizacionUseCase{}, fmt.Errorf("obtener tabla canalización: %w", err)
	}

	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: hilosPorFase * 3, SeccionMM2: conductorAlimentacion.SeccionMM2()}, // Fases
		{Cantidad: 1, SeccionMM2: conductorTierra.SeccionMM2()},                      // Tierra
	}

	resultado, err := service.CalcularCanalizacion(conductores, string(tipoCanalizacion), tablaCanalizacion, hilosPorFase)
	if err != nil {
		return ResultadoCanalizacionUseCase{}, fmt.Errorf("calcular canalización: %w", err)
	}

	return ResultadoCanalizacionUseCase{
		Tamano:        resultado.Tamano,
		NumeroDeTubos: resultado.NumeroDeTubos,
		AreaTotalMM2:  resultado.AnchoRequerido,
	}, nil
}
