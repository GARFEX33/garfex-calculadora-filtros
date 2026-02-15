// internal/calculos/infrastructure/adapter/driven/csv/seleccionar_temperatura.go
package csv

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// SeleccionarTemperaturaRepository implements the SeleccionarTemperaturaPort interface
// by delegating to the domain service.
type SeleccionarTemperaturaRepository struct{}

// NewSeleccionarTemperaturaRepository creates a new SeleccionarTemperaturaRepository.
func NewSeleccionarTemperaturaRepository() *SeleccionarTemperaturaRepository {
	return &SeleccionarTemperaturaRepository{}
}

// SeleccionarTemperatura delegates to the domain service.
func (r *SeleccionarTemperaturaRepository) SeleccionarTemperatura(
	corriente valueobject.Corriente,
	tipoCanalizacion entity.TipoCanalizacion,
	override *valueobject.Temperatura,
) valueobject.Temperatura {
	return service.SeleccionarTemperatura(corriente, tipoCanalizacion, override)
}

// Ensure SeleccionarTemperaturaRepository implements the port interface.
var _ port.SeleccionarTemperaturaPort = (*SeleccionarTemperaturaRepository)(nil)
