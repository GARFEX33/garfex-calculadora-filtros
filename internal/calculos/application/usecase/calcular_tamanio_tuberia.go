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

	// Total de conductores por tipo (fases × hilos_por_fase = conductores totales de fase)
	// El domain service distribuye entre numTuberias para obtener el área por tubo.
	hilosPorFaseForDomain := input.GetHilosPorFase()
	totalConductoresFase := input.NumFases * hilosPorFaseForDomain
	totalConductoresNeutro := input.NumNeutros * hilosPorFaseForDomain

	// Call domain service to calculate conduit size
	resultado, err := service.CalcularTamanioTuberiaWithMultiplePipes(
		totalConductoresFase,
		totalConductoresNeutro,
		input.NumTierras, // Usar valor calculado según normativa NOM (total = numTuberias)
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

	// Buscar el área de ocupación del tubo seleccionado en la tabla
	var areaOcupacionSeleccionada float64
	var designacionMetrica string
	for _, entrada := range tablaOcupacion {
		if entrada.Tamano == resultado.TuberiaRecomendada() {
			areaOcupacionSeleccionada = entrada.AreaOcupacionMM2
			designacionMetrica = entrada.DesignacionMetrica
			break
		}
	}

	// Preparar puntero a areaNeutro (nil si no hay neutro)
	var areaNeutroPtr *float64
	if input.NumNeutros > 0 {
		areaNeutroPtr = &areaNeutro
	}

	// Conductores por tubo: (fases × hilos_por_fase) / num_tuberias
	// Ejemplo: 3F-3H, 2 hilos/fase, 2 tubos → (3×2)/2 = 3 conductores de fase por tubo
	hilosPorFase := input.GetHilosPorFase()
	numFasesPorTubo := (input.NumFases * hilosPorFase) / input.NumTuberias
	numNeutrosPorTubo := (input.NumNeutros * hilosPorFase) / input.NumTuberias

	// Tierras por tubo: en tubería siempre 1 tierra por tubo (NOM).
	// input.NumTierras es el TOTAL (= num_tuberias) — para el desarrollo mostramos 1 por tubo.
	numTierrasPorTubo := input.GetNumeroTierras() / input.NumTuberias
	if numTierrasPorTubo < 1 {
		numTierrasPorTubo = 1
	}

	// Map domain result to DTO
	return dto.TuberiaOutput{
		AreaPorTuboMM2:     resultado.AreaPorTuboMM2(),
		TuberiaRecomendada: resultado.TuberiaRecomendada(),
		DesignacionMetrica: designacionMetrica,
		TipoCanalizacion:   string(resultado.TipoCanalizacion()),
		NumTuberias:        resultado.NumTuberias(),
		// Nuevos campos
		AreaFaseMM2:          areaFase,
		AreaNeutroMM2:        areaNeutroPtr,
		AreaTierraMM2:        areaTierra,
		NumFasesPorTubo:      numFasesPorTubo,
		NumNeutrosPorTubo:    numNeutrosPorTubo,
		NumTierras:           numTierrasPorTubo, // 1 tierra por tubo (no el total)
		AreaOcupacionTuboMM2: areaOcupacionSeleccionada,
		FillFactor:           0.40,
	}, nil
}
