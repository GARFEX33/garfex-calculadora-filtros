// internal/calculos/application/usecase/calcular_corriente.go
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
		return uc.calcularDesdeListado(input)

	case dto.ModoManualAmperaje:
		return uc.calcularManualAmperaje(input)

	case dto.ModoManualPotencia:
		return uc.calcularManualPotencia(input)

	default:
		return dto.ResultadoCorriente{}, dto.ErrModoInvalido
	}
}

// calcularDesdeListado calculates current from equipment data.
func (uc *CalcularCorrienteUseCase) calcularDesdeListado(input dto.EquipoInput) (dto.ResultadoCorriente, error) {
	// Obtener tipo de equipo mapeado desde TipoFiltro
	tipoEquipo, err := input.GetTipoEquipo()
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("mapear tipo de equipo: %w", err)
	}

	// Crear la entidad correcta según el tipo
	var calculador entity.CalculadorCorriente

	switch tipoEquipo {
	case entity.TipoEquipoFiltroActivo:
		// FiltroActivo: el amperaje es la corriente directa
		calculador, err = entity.NewFiltroActivo(
			input.Equipo.Clave,
			input.Equipo.Voltaje,
			input.Equipo.Amperaje,
			entity.ITM{Amperaje: input.Equipo.ITM},
		)

	case entity.TipoEquipoFiltroRechazo:
		// FiltroRechazo: el amperaje es KVAR
		calculador, err = entity.NewFiltroRechazo(
			input.Equipo.Clave,
			input.Equipo.Voltaje,
			input.Equipo.Amperaje,
			entity.ITM{Amperaje: input.Equipo.ITM},
		)

	case entity.TipoEquipoTransformador:
		// Transformador: el amperaje es KVA
		calculador, err = entity.NewTransformador(
			input.Equipo.Clave,
			input.Equipo.Voltaje,
			input.Equipo.Amperaje,
			entity.ITM{Amperaje: input.Equipo.ITM},
		)

	default:
		return dto.ResultadoCorriente{}, fmt.Errorf("tipo de equipo no soportado: %s", tipoEquipo)
	}

	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("crear entidad de equipo: %w", err)
	}

	// Calcular corriente nominal
	corriente, err := service.CalcularCorrienteNominal(calculador)
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
		return dto.ResultadoCorriente{}, fmt.Errorf("crear corriente desde amperaje: %w", err)
	}

	return dto.ResultadoCorriente{CorrienteNominal: corriente.Valor()}, nil
}

// calcularManualPotencia calculates current from manual power.
func (uc *CalcularCorrienteUseCase) calcularManualPotencia(input dto.EquipoInput) (dto.ResultadoCorriente, error) {
	potencia, err := input.ToDomainPotencia()
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("potencia inválida: %w", err)
	}

	sistema := input.SistemaElectrico.ToEntity()

	tension, err := input.ToDomainTension()
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("tensión inválida: %w", err)
	}

	corriente, err := service.CalcularAmperajeNominalCircuito(
		potencia,
		tension,
		sistema,
		input.FactorPotencia,
	)
	if err != nil {
		return dto.ResultadoCorriente{}, fmt.Errorf("calcular amperaje desde potencia: %w", err)
	}

	return dto.ResultadoCorriente{CorrienteNominal: corriente.Valor()}, nil
}
