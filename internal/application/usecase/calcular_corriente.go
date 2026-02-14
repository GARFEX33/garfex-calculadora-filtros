// internal/application/usecase/calcular_corriente.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ResultadoCorriente contiene el resultado del cálculo de corriente.
type ResultadoCorriente struct {
	Nominal valueobject.Corriente
}

// CalcularCorrienteUseCase ejecuta el Paso 1: Corriente Nominal.
type CalcularCorrienteUseCase struct {
	equipoRepo port.EquipoRepository
}

// NewCalcularCorrienteUseCase crea una nueva instancia.
func NewCalcularCorrienteUseCase(
	equipoRepo port.EquipoRepository,
) *CalcularCorrienteUseCase {
	return &CalcularCorrienteUseCase{
		equipoRepo: equipoRepo,
	}
}

// Execute calcula la corriente nominal según el modo.
func (uc *CalcularCorrienteUseCase) Execute(ctx context.Context, input dto.EquipoInput) (ResultadoCorriente, error) {
	switch input.Modo {
	case dto.ModoListado:
		if input.Clave == "" {
			return ResultadoCorriente{}, dto.ErrEquipoInputInvalido
		}
		equipo, err := uc.equipoRepo.BuscarPorClave(ctx, input.Clave)
		if err != nil {
			return ResultadoCorriente{}, fmt.Errorf("buscar equipo: %w", err)
		}
		// Calcular corriente según tipo de equipo
		corriente, err := equipo.CalcularCorrienteNominal()
		if err != nil {
			return ResultadoCorriente{}, fmt.Errorf("calcular corriente: %w", err)
		}
		return ResultadoCorriente{Nominal: corriente}, nil

	case dto.ModoManualAmperaje:
		if input.AmperajeNominal <= 0 {
			return ResultadoCorriente{}, dto.ErrEquipoInputInvalido
		}
		corriente, err := valueobject.NewCorriente(input.AmperajeNominal)
		if err != nil {
			return ResultadoCorriente{}, err
		}
		return ResultadoCorriente{Nominal: corriente}, nil

	case dto.ModoManualPotencia:
		return ResultadoCorriente{}, fmt.Errorf("modo MANUAL_POTENCIA requiere implementación adicional")

	default:
		return ResultadoCorriente{}, dto.ErrModoInvalido
	}
}
