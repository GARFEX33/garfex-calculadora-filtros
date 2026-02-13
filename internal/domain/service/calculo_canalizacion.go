package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ErrCanalizacionNoDisponible is returned when no conduit size fits the required area.
var ErrCanalizacionNoDisponible = errors.New("no se encontró canalización con área suficiente")

// factorRellenoTuberia is the NOM fill factor for conduit with 2+ conductors (40%).
const factorRellenoTuberia = 0.40

func determinarFillFactor(cantidad int) float64 {
	switch cantidad {
	case 1:
		return 0.53
	case 2:
		return 0.31
	default:
		return 0.40
	}
}

// ConductorParaCanalizacion holds the quantity and cross-section area
// of a group of identical conductors for conduit sizing calculations.
type ConductorParaCanalizacion struct {
	Cantidad   int
	SeccionMM2 float64
}

// CalcularCanalizacion selects the smallest conduit whose usable area
// (interior area × fill factor) accommodates all conductors.
// tipo should be a TipoCanalizacion string value (e.g., "TUBERIA_CONDUIT").
func CalcularCanalizacion(
	conductores []ConductorParaCanalizacion,
	tipo string,
	tabla []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	if len(conductores) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("lista de conductores vacía")
	}
	if len(tabla) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("%w: tabla vacía", ErrCanalizacionNoDisponible)
	}

	var areaTotal float64
	var cantidadTotal int
	for _, c := range conductores {
		areaTotal += float64(c.Cantidad) * c.SeccionMM2
		cantidadTotal += c.Cantidad
	}

	factorRelleno := determinarFillFactor(cantidadTotal)
	areaRequerida := areaTotal / factorRelleno

	for _, entrada := range tabla {
		if entrada.AreaInteriorMM2 >= areaRequerida {
			return entity.Canalizacion{
				Tipo:           tipo,
				Tamano:         entrada.Tamano,
				AnchoRequerido: areaTotal,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"%w: área requerida %.2f mm² excede máxima disponible %.2f mm²",
		ErrCanalizacionNoDisponible, areaRequerida, tabla[len(tabla)-1].AreaInteriorMM2,
	)
}
