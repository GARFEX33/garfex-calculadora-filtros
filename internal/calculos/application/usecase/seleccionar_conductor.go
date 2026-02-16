// internal/calculos/application/usecase/seleccionar_conductor.go
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

// SeleccionarConductorUseCase ejecuta los Pasos 4 y 5: Selección de conductores.
type SeleccionarConductorUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewSeleccionarConductorUseCase crea una nueva instancia.
func NewSeleccionarConductorUseCase(
	tablaRepo port.TablaNOMRepository,
) *SeleccionarConductorUseCase {
	return &SeleccionarConductorUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute selecciona conductor de alimentación y tierra.
// Retorna un DTO plano para evitar domain bleeding.
func (uc *SeleccionarConductorUseCase) Execute(
	ctx context.Context,
	corrienteAjustada valueobject.Corriente,
	hilosPorFase int,
	itm int,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
	tipoCanalizacion entity.TipoCanalizacion,
) (dto.ResultadoConductores, error) {
	// Obtener tabla de ampacidad
	tablaAmpacidad, err := uc.tablaRepo.ObtenerTablaAmpacidad(ctx, tipoCanalizacion, material, temperatura)
	if err != nil {
		return dto.ResultadoConductores{}, fmt.Errorf("obtener tabla ampacidad: %w", err)
	}

	// Seleccionar conductor de alimentación
	conductor, err := service.SeleccionarConductorAlimentacion(corrienteAjustada, hilosPorFase, tablaAmpacidad)
	if err != nil {
		return dto.ResultadoConductores{}, fmt.Errorf("seleccionar conductor alimentación: %w", err)
	}

	// Obtener capacidad del conductor seleccionado
	capacidad, err := uc.tablaRepo.ObtenerCapacidadConductor(ctx, tipoCanalizacion, material, temperatura, conductor.Calibre())
	if err != nil {
		return dto.ResultadoConductores{}, fmt.Errorf("obtener capacidad: %w", err)
	}

	// Seleccionar conductor de tierra
	tablaTierra, err := uc.tablaRepo.ObtenerTablaTierra(ctx)
	if err != nil {
		return dto.ResultadoConductores{}, fmt.Errorf("obtener tabla tierra: %w", err)
	}

	conductorTierra, err := service.SeleccionarConductorTierra(itm, material, tablaTierra)
	if err != nil {
		return dto.ResultadoConductores{}, fmt.Errorf("seleccionar conductor tierra: %w", err)
	}

	// Determinar nombre de tabla usada según canalización
	tablaUsada := helpers.NombreTablaAmpacidad(string(tipoCanalizacion), material, temperatura)

	return dto.ResultadoConductores{
		Alimentacion: dto.ResultadoConductor{
			Calibre:         conductor.Calibre(),
			Material:        conductor.Material().String(),
			SeccionMM2:      conductor.SeccionMM2(),
			TipoAislamiento: conductor.TipoAislamiento(),
			Capacidad:       capacidad,
		},
		Tierra: dto.ResultadoConductor{
			Calibre:    conductorTierra.Calibre(),
			Material:   conductorTierra.Material().String(),
			SeccionMM2: conductorTierra.SeccionMM2(),
		},
		TablaUsada: tablaUsada,
	}, nil
}
