// internal/application/port/tabla_nom_repository.go
package port

import (
	"context"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
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
	) ([]valueobject.EntradaTablaConductor, error)

	// ObtenerCapacidadConductor returns the ampacity for a specific calibre.
	ObtenerCapacidadConductor(
		ctx context.Context,
		canalizacion entity.TipoCanalizacion,
		material valueobject.MaterialConductor,
		temperatura valueobject.Temperatura,
		calibre string,
	) (float64, error)

	// ObtenerTablaTierra returns the ground conductor table (250-122).
	ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error)

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
	) ([]valueobject.EntradaTablaCanalizacion, error)

	// Tablas de factores
	ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error)
	ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error)
	ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error)

	// Dimensiones para canalización
	ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error)
	ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error)
}
