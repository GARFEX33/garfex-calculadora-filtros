// internal/calculos/application/usecase/seleccionar_conductor_alimentacion.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase/helpers"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// SeleccionarConductorAlimentacionUseCase ejecuta la seleccion de conductor
// de alimentacion segun tablas NOM 310-15.
type SeleccionarConductorAlimentacionUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewSeleccionarConductorAlimentacionUseCase crea una nueva instancia.
func NewSeleccionarConductorAlimentacionUseCase(
	tablaRepo port.TablaNOMRepository,
) *SeleccionarConductorAlimentacionUseCase {
	return &SeleccionarConductorAlimentacionUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute selecciona el conductor de alimentacion apropiado.
func (uc *SeleccionarConductorAlimentacionUseCase) Execute(
	ctx context.Context,
	input dto.ConductorAlimentacionInput,
) (dto.ConductorAlimentacionOutput, error) {
	// 1. Validar DTO
	if err := input.Validate(); err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("validar input: %w", err)
	}

	// 2. Convertir primitivos a value objects
	corrienteAjustada, err := valueobject.NewCorriente(input.CorrienteAjustada)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("corriente invalida: %w", err)
	}

	tipoCanalizacion, err := entity.ParseTipoCanalizacion(input.TipoCanalizacion)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("tipo canalizacion invalido: %w", err)
	}

	material := input.ToDomainMaterial()

	hilosPorFase := input.HilosPorFase
	if hilosPorFase < 1 {
		hilosPorFase = 1
	}

	// 3. Determinar temperatura (input o regla NOM)
	var temperatura valueobject.Temperatura
	if input.Temperatura != nil {
		temperatura = valueobject.Temperatura(*input.Temperatura)
	} else {
		temperatura = service.SeleccionarTemperatura(corrienteAjustada, tipoCanalizacion, nil)
	}

	// 4. Obtener tabla de ampacidad
	tablaAmpacidad, err := uc.tablaRepo.ObtenerTablaAmpacidad(ctx, tipoCanalizacion, material, temperatura)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("obtener tabla ampacidad: %w", err)
	}

	// 5. Llamar servicio de dominio
	conductor, err := service.SeleccionarConductorAlimentacion(corrienteAjustada, hilosPorFase, tablaAmpacidad)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("seleccionar conductor: %w", err)
	}

	// 6. Obtener capacidad del conductor seleccionado
	capacidad, err := uc.tablaRepo.ObtenerCapacidadConductor(ctx, tipoCanalizacion, material, temperatura, conductor.Calibre())
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("obtener capacidad: %w", err)
	}

	// 7. Generar nombre de tabla usada
	tablaUsada := helpers.NombreTablaAmpacidad(string(tipoCanalizacion), material, temperatura)

	// 8. Retornar DTO output
	return dto.ConductorAlimentacionOutput{
		Calibre:          conductor.Calibre(),
		Material:         conductor.Material().String(),
		SeccionMM2:       conductor.SeccionMM2(),
		TipoAislamiento:  conductor.TipoAislamiento(),
		CapacidadNominal: capacidad,
		TablaUsada:       tablaUsada,
	}, nil
}
