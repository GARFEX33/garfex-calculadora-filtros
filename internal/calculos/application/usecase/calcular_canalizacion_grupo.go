// internal/calculos/application/usecase/calcular_canalizacion_grupo.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
)

// CalcularCanalizacionGrupoUseCase ejecuta el cálculo de canalización para grupos de conductores.
type CalcularCanalizacionGrupoUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewCalcularCanalizacionGrupoUseCase crea una nueva instancia.
func NewCalcularCanalizacionGrupoUseCase(
	tablaRepo port.TablaNOMRepository,
) *CalcularCanalizacionGrupoUseCase {
	return &CalcularCanalizacionGrupoUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute calcula la canalización para un grupo de conductores.
// Agrega automáticamente conductores de tierra (1 por tubo) según la regla NOM.
func (uc *CalcularCanalizacionGrupoUseCase) Execute(
	ctx context.Context,
	input dto.CanalizacionGrupoInput,
) (dto.CanalizacionGrupoOutput, error) {
	// 1. Validar input
	if err := input.Validate(); err != nil {
		return dto.CanalizacionGrupoOutput{}, fmt.Errorf("validar input: %w", err)
	}

	// 2. Convertir conductores DTO a domain
	conductores := make([]service.ConductorParaCanalizacion, 0, len(input.Conductores)+input.NumeroDeTubos)
	for _, c := range input.Conductores {
		conductores = append(conductores, service.ConductorParaCanalizacion{
			Cantidad:   c.Cantidad,
			SeccionMM2: c.SeccionMM2,
		})
	}

	// 3. Agregar conductores de tierra (1 por tubo, según norma)
	for i := 0; i < input.NumeroDeTubos; i++ {
		conductores = append(conductores, service.ConductorParaCanalizacion{
			Cantidad:   1,
			SeccionMM2: input.SeccionTierraMM2,
		})
	}

	// 4. Parsear tipo de canalización
	tipoCanalizacion, err := entity.ParseTipoCanalizacion(input.TipoCanalizacion)
	if err != nil {
		return dto.CanalizacionGrupoOutput{}, fmt.Errorf("parsear tipo canalización: %w", err)
	}

	// 5. Obtener tabla NOM
	tabla, err := uc.tablaRepo.ObtenerTablaCanalizacion(ctx, tipoCanalizacion)
	if err != nil {
		return dto.CanalizacionGrupoOutput{}, fmt.Errorf("obtener tabla canalización: %w", err)
	}

	// 6. Calcular canalización
	resultado, err := service.CalcularCanalizacion(conductores, tipoCanalizacion, tabla, input.NumeroDeTubos)
	if err != nil {
		return dto.CanalizacionGrupoOutput{}, fmt.Errorf("calcular canalización: %w", err)
	}

	// 7. Mapear resultado a DTO
	return dto.CanalizacionGrupoOutput{
		Tamano:         resultado.Tamano,
		AreaTotalMM2:   resultado.AnchoRequerido,
		AreaPorTuboMM2: resultado.AnchoRequerido / float64(resultado.NumeroDeTubos),
		NumeroDeTubos:  resultado.NumeroDeTubos,
		FactorRelleno:  resultado.FactorRelleno,
	}, nil
}
