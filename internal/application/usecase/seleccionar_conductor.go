// internal/application/usecase/seleccionar_conductor.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/application/usecase/helpers"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ResultadoConductores contiene los conductores seleccionados.
type ResultadoConductores struct {
	Alimentacion valueobject.Conductor
	Tierra       valueobject.Conductor
	TablaUsada   string
	Capacidad    float64
}

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
func (uc *SeleccionarConductorUseCase) Execute(
	ctx context.Context,
	corrienteAjustada valueobject.Corriente,
	hilosPorFase int,
	itm int,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
	tipoCanalizacion entity.TipoCanalizacion,
) (ResultadoConductores, error) {
	// Obtener tabla de ampacidad
	tablaAmpacidad, err := uc.tablaRepo.ObtenerTablaAmpacidad(ctx, tipoCanalizacion, material, temperatura)
	if err != nil {
		return ResultadoConductores{}, fmt.Errorf("obtener tabla ampacidad: %w", err)
	}

	// Seleccionar conductor de alimentación
	conductor, err := service.SeleccionarConductorAlimentacion(corrienteAjustada, hilosPorFase, tablaAmpacidad)
	if err != nil {
		return ResultadoConductores{}, fmt.Errorf("seleccionar conductor alimentación: %w", err)
	}

	// Obtener capacidad del conductor seleccionado
	capacidad, err := uc.tablaRepo.ObtenerCapacidadConductor(ctx, tipoCanalizacion, material, temperatura, conductor.Calibre())
	if err != nil {
		return ResultadoConductores{}, fmt.Errorf("obtener capacidad: %w", err)
	}

	// Seleccionar conductor de tierra
	tablaTierra, err := uc.tablaRepo.ObtenerTablaTierra(ctx)
	if err != nil {
		return ResultadoConductores{}, fmt.Errorf("obtener tabla tierra: %w", err)
	}

	conductorTierra, err := service.SeleccionarConductorTierra(itm, material, tablaTierra)
	if err != nil {
		return ResultadoConductores{}, fmt.Errorf("seleccionar conductor tierra: %w", err)
	}

	// Determinar nombre de tabla usada según canalización
	tablaUsada := helpers.NombreTablaAmpacidad(string(tipoCanalizacion), material, temperatura)

	return ResultadoConductores{
		Alimentacion: conductor,
		Tierra:       conductorTierra,
		Capacidad:    capacidad,
		TablaUsada:   tablaUsada,
	}, nil
}
