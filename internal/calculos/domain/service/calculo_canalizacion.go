// internal/calculos/domain/service/calculo_canalizacion.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ErrCanalizacionNoDisponible is returned when no conduit size fits the required area.
var ErrCanalizacionNoDisponible = errors.New("no se encontró canalización con área suficiente")

// ErrNumeroDeTubosInvalido is returned when numeroDeTubos is less than 1.
var ErrNumeroDeTubosInvalido = errors.New("numeroDeTubos debe ser mayor a cero")

// ErrListaConductoresVacia is returned when the conductors list is empty.
var ErrListaConductoresVacia = errors.New("lista de conductores vacía")

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
// numeroDeTubos indicates how many parallel conduits to use (must be >= 1).
// When numeroDeTubos > 1, the total conductor area and count are divided evenly
// among the tubes; fill factor is determined per-tube conductor count.
func CalcularCanalizacion(
	conductores []ConductorParaCanalizacion,
	tipo entity.TipoCanalizacion,
	tabla []valueobject.EntradaTablaCanalizacion,
	numeroDeTubos int,
) (entity.Canalizacion, error) {
	if numeroDeTubos < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCanalizacion: %w", ErrNumeroDeTubosInvalido)
	}
	if len(conductores) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCanalizacion: %w", ErrListaConductoresVacia)
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

	conductoresPorTubo := cantidadTotal / numeroDeTubos
	factorRelleno := determinarFillFactor(conductoresPorTubo)
	areaPorTubo := areaTotal / float64(numeroDeTubos)
	areaRequerida := areaPorTubo / factorRelleno

	for _, entrada := range tabla {
		if entrada.AreaInteriorMM2 >= areaRequerida {
			resultado, err := entity.NewCanalizacion(
				tipo,
				entrada.Tamano,
				areaTotal,
				numeroDeTubos,
				factorRelleno,
			)
			if err != nil {
				return entity.Canalizacion{}, fmt.Errorf("crear canalización: %w", err)
			}
			return resultado, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"%w: área requerida %.2f mm² excede máxima disponible %.2f mm²",
		ErrCanalizacionNoDisponible, areaRequerida, tabla[len(tabla)-1].AreaInteriorMM2,
	)
}
