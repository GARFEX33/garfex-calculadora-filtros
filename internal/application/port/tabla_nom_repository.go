// internal/application/port/tabla_nom_repository.go
package port

import (
	"context"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// TablaNOMRepository defines the contract for reading NOM tables.
type TablaNOMRepository interface {
	// ObtenerTablaAmpacidad returns ampacity table entries for the given conduit type, material, and temperature.
	ObtenerTablaAmpacidad(
		ctx context.Context,
		canalizacion entity.TipoCanalizacion,
		material valueobject.MaterialConductor,
		temperatura valueobject.Temperatura,
	) ([]service.EntradaTablaConductor, error)

	// ObtenerTablaTierra returns the ground conductor table (250-122).
	ObtenerTablaTierra(ctx context.Context) ([]service.EntradaTablaTierra, error)

	// ObtenerImpedancia returns R and X values for the given calibre and conduit type.
	ObtenerImpedancia(
		ctx context.Context,
		calibre string,
		canalizacion entity.TipoCanalizacion,
		material valueobject.MaterialConductor,
	) (valueobject.ResistenciaReactancia, error)

	// ObtenerTablaCanalizacion returns conduit sizing table entries.
	ObtenerTablaCanalizacion(
		ctx context.Context,
		canalizacion entity.TipoCanalizacion,
	) ([]service.EntradaTablaCanalizacion, error)
}
