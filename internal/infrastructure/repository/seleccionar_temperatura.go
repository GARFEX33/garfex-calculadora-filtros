// internal/infrastructure/repository/seleccionar_temperatura.go
package repository

import (
	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
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
