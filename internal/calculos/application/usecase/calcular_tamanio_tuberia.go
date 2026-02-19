// internal/calculos/application/usecase/calcular_tamanio_tuberia.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
)

// CalcularTamanioTuberiaUseCase executes the conduit sizing calculation
// based on conductor areas and NOM Chapter 9 fill requirements.
type CalcularTamanioTuberiaUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewCalcularTamanioTuberiaUseCase creates a new instance.
func NewCalcularTamanioTuberiaUseCase(
	tablaRepo port.TablaNOMRepository,
) *CalcularTamanioTuberiaUseCase {
	return &CalcularTamanioTuberiaUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute calculates the appropriate conduit size based on conductor areas.
func (uc *CalcularTamanioTuberiaUseCase) Execute(
	ctx context.Context,
	input dto.TuberiaInput,
) (dto.TuberiaOutput, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		return dto.TuberiaOutput{}, fmt.Errorf("validar input: %w", err)
	}

	// Parse tipo canalización
	tipoCanalizacion := entity.TipoCanalizacion(input.TipoCanalizacion)
	if err := entity.ValidarTipoCanalizacion(tipoCanalizacion); err != nil {
		return dto.TuberiaOutput{}, fmt.Errorf("validar tipo canalización: %w", err)
	}

	// Get areas for each conductor type
	areaFase, err := uc.tablaRepo.ObtenerAreaConductor(ctx, input.CalibreFase)
	if err != nil {
		return dto.TuberiaOutput{}, fmt.Errorf("obtener área fase %s: %w", input.CalibreFase, err)
	}

	// Get area for neutral only if there are neutrals
	var areaNeutro float64
	if input.NumNeutros > 0 {
		areaNeutro, err = uc.tablaRepo.ObtenerAreaConductor(ctx, input.CalibreNeutro)
		if err != nil {
			return dto.TuberiaOutput{}, fmt.Errorf("obtener área neutro %s: %w", input.CalibreNeutro, err)
		}
	}

	// El conductor de tierra es desnudo - usar tabla 8 (conductor desnudo)
	areaTierra, err := uc.tablaRepo.ObtenerAreaConductorDesnudo(ctx, input.CalibreTierra)
	if err != nil {
		return dto.TuberiaOutput{}, fmt.Errorf("obtener área tierra %s: %w", input.CalibreTierra, err)
	}

	// Get conduit occupation table
	tablaOcupacion, err := uc.tablaRepo.ObtenerTablaOcupacionTuberia(ctx, tipoCanalizacion)
	if err != nil {
		return dto.TuberiaOutput{}, fmt.Errorf("obtener tabla ocupación: %w", err)
	}

	// Call domain service to calculate conduit size
	resultado, err := service.CalcularTamanioTuberiaWithMultiplePipes(
		input.NumFases,
		input.NumNeutros,
		1, // 1 tierra por tubo
		areaFase,
		areaNeutro,
		areaTierra,
		input.NumTuberias,
		tipoCanalizacion,
		tablaOcupacion,
	)
	if err != nil {
		return dto.TuberiaOutput{}, fmt.Errorf("calcular tamaño tubería: %w", err)
	}

	// Map domain result to DTO
	return dto.TuberiaOutput{
		AreaPorTuboMM2:     resultado.AreaPorTuboMM2(),
		TuberiaRecomendada: resultado.TuberiaRecomendada(),
		DesignacionMetrica: resultado.DesignacionMetrica(),
		TipoCanalizacion:   string(resultado.TipoCanalizacion()),
		NumTuberias:        resultado.NumTuberias(),
	}, nil
}
