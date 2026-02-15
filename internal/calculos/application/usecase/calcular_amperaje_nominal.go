// internal/calculos/application/usecase/calcular_amperaje_nominal.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// CalcularAmperajeNominalUseCase orquesta el cálculo de amperaje nominal desde potencia.
// No tiene dependencias externas (usa servicios de dominio puros).
type CalcularAmperajeNominalUseCase struct{}

// NewCalcularAmperajeNominalUseCase crea una nueva instancia.
func NewCalcularAmperajeNominalUseCase() *CalcularAmperajeNominalUseCase {
	return &CalcularAmperajeNominalUseCase{}
}

// Execute calcula el amperaje nominal de un circuito eléctrico.
//
// Parámetros:
//   - ctx: contexto de la operación
//   - input: DTO con potencia, tensión, tipo de carga, sistema eléctrico y factor de potencia
//
// Retorna:
//   - AmperajeNominalOutput con el amperaje calculado
//   - Error si la validación falla o el cálculo no puede completarse
func (uc *CalcularAmperajeNominalUseCase) Execute(
	ctx context.Context,
	input dto.AmperajeNominalInput,
) (dto.AmperajeNominalOutput, error) {
	// Validar input
	if err := input.Validate(); err != nil {
		return dto.AmperajeNominalOutput{}, fmt.Errorf("validar input: %w", err)
	}

	// Crear value object Tension
	tension, err := valueobject.NewTension(input.Tension)
	if err != nil {
		return dto.AmperajeNominalOutput{}, fmt.Errorf("crear tensión: %w", err)
	}

	// Convertir DTOs a tipos del dominio
	tipoCarga := input.TipoCarga.ToEntity()
	sistemaElectrico := input.SistemaElectrico.ToEntity()

	// Llamar al servicio de dominio
	corriente, err := service.CalcularAmperajeNominalCircuito(
		input.PotenciaWatts,
		tension,
		tipoCarga,
		sistemaElectrico,
		input.FactorPotencia,
	)
	if err != nil {
		return dto.AmperajeNominalOutput{}, fmt.Errorf("calcular amperaje nominal: %w", err)
	}

	// Mappear resultado a DTO de salida
	return dto.AmperajeNominalOutput{
		Amperaje: corriente.Valor(),
		Unidad:   "A",
	}, nil
}
