package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

var ErrCharolaNoEncontrada = errors.New("no se encontró charola suficiente")

type ConductorConDiametro struct {
	DiametroMM float64
}

func CalcularCharolaEspaciado(
	hilosPorFase int,
	sistema entity.SistemaElectrico,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", errors.New("hilos por fase debe ser >= 1"))
	}

	hilosFaseTotal := hilosPorFase * 3

	var hilosNeutro int
	if sistema == entity.SistemaElectricoEstrella ||
		sistema == entity.SistemaElectricoMonofasico {
		hilosNeutro = 1
	}

	totalHilos := hilosFaseTotal + hilosNeutro
	anchoRequerido := float64(totalHilos-1)*conductorFase.DiametroMM + conductorTierra.DiametroMM

	const alturaCharolaMM float64 = 50.0
	for _, entrada := range tablaCharola {
		anchoCharolaMM := entrada.AreaInteriorMM2 / alturaCharolaMM
		if anchoCharolaMM >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:      string(entity.TipoCanalizacionCharolaCableEspaciado),
				Tamano:    entrada.Tamano,
				AreaTotal: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", ErrCharolaNoEncontrada)
}
