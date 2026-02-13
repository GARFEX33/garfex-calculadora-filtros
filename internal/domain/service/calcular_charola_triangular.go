package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

func CalcularCharolaTriangular(
	hilosPorFase int,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: hilos por fase debe ser >= 1: %d", hilosPorFase)
	}
	if len(tablaCharola) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: tabla vacía")
	}

	const factorTriangular = 2.15

	anchoRequerido := (float64(hilosPorFase-1) * factorTriangular * conductorFase.DiametroMM) +
		conductorTierra.DiametroMM

	for _, entrada := range tablaCharola {
		if entrada.AreaInteriorMM2 >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:      string(entity.TipoCanalizacionCharolaCableTriangular),
				Tamano:    entrada.Tamano,
				AreaTotal: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"CalcularCharolaTriangular: no se encontró charola triangular suficiente: ancho requerido %.2f mm",
		anchoRequerido,
	)
}
