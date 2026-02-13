package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

var ErrCharolaTriangularNoEncontrada = errors.New("no se encontró charola triangular suficiente")

func CalcularCharolaTriangular(
	hilosPorFase int,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", errors.New("hilos por fase debe ser >= 1"))
	}
	if len(tablaCharola) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", errors.New("tabla vacía"))
	}

	const factorTriangular = 2.15

	anchoRequerido := (float64(hilosPorFase-1) * factorTriangular * conductorFase.DiametroMM) +
		conductorTierra.DiametroMM

	for _, entrada := range tablaCharola {
		if entrada.AreaInteriorMM2 >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:           string(entity.TipoCanalizacionCharolaCableTriangular),
				Tamano:         entrada.Tamano,
				AnchoRequerido: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", ErrCharolaTriangularNoEncontrada)
}
