// internal/calculos/application/usecase/calcular_corriente.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// CalcularCorrienteUseCase executes Step 1: Nominal Current.
type CalcularCorrienteUseCase struct {
	equipoRepo port.EquipoRepository
}

// NewCalcularCorrienteUseCase creates a new instance.
func NewCalcularCorrienteUseCase(
	equipoRepo port.EquipoRepository,
) *CalcularCorrienteUseCase {
	return &CalcularCorrienteUseCase{
		equipoRepo: equipoRepo,
	}
}

// Execute calculates the nominal current according to the mode.
func (uc *CalcularCorrienteUseCase) Execute(ctx context.Context, input dto.EquipoInput) (dto.ResultadoCorriente, error) {
	switch input.Modo {
	case dto.ModoListado:
		return uc.calcularDesdeListado(ctx, input)

	case dto.ModoManualAmperaje:
		return uc.calcularManualAmperaje(input)

	case dto.ModoManualPotencia:
		return uc.calcularManualPotencia(input)

	default:
		return dto.ResultadoCorriente{}, dto.ErrModoInvalido
	}
}

// calcularDesdeListado calculates current from equipment listing.
func (uc *CalcularCorrienteUseCase) calcularDesdeListado(ctx context.Context, input dto.EquipoInput) (dto.ResultadoCorriente, error) {
	if input.Clave == "" {
		return dto.ResultadoCorriente{}, dto.ErrEquipoInputInvalido
	}

	equipo, err := uc.equipoRepo.BuscarPorClave(ctx, input.Clave)
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("buscar equipo: %w", err)
	}

	corriente, err := service.CalcularCorrienteNominal(equipo)
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("calcular corriente: %w", err)
	}

	return dto.ResultadoCorriente{CorrienteNominal: corriente.Valor()}, nil
}

// calcularManualAmperaje calculates current from manual amperage.
func (uc *CalcularCorrienteUseCase) calcularManualAmperaje(input dto.EquipoInput) (dto.ResultadoCorriente, error) {
	if input.AmperajeNominal <= 0 {
		return dto.ResultadoCorriente{}, dto.ErrEquipoInputInvalido
	}

	corriente, err := valueobject.NewCorriente(input.AmperajeNominal)
	if err != nil {
		return dto.ResultadoCorriente{}, err
	}

	return dto.ResultadoCorriente{CorrienteNominal: corriente.Valor()}, nil
}

// calcularManualPotencia calculates current from manual power using entity.SistemaElectrico.
// The DTO SistemaElectrico converts directly to entity via ToEntity().
func (uc *CalcularCorrienteUseCase) calcularManualPotencia(input dto.EquipoInput) (dto.ResultadoCorriente, error) {
	if input.PotenciaNominal <= 0 {
		return dto.ResultadoCorriente{}, dto.ErrEquipoInputInvalido
	}

	// Convertir DTO a entity — los valores son idénticos (DELTA, ESTRELLA, etc.)
	sistema := input.SistemaElectrico.ToEntity()

	corriente, err := service.CalcularAmperajeNominalCircuito(
		input.PotenciaNominal,
		input.Tension,
		sistema,
		input.FactorPotencia,
	)
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("calcular amperaje desde potencia: %w", err)
	}

	return dto.ResultadoCorriente{CorrienteNominal: corriente.Valor()}, nil
}
