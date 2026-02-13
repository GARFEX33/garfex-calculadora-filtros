package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
)

type ConductorConDiametro struct {
	DiametroMM float64
}

func CalcularCharolaEspaciado(
	hilosPorFase int,
	sistema entity.SistemaElectrico,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []struct {
		Tamano  string
		AnchoMM float64
	},
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: hilos por fase debe ser >= 1")
	}

	hilosFaseTotal := hilosPorFase * 3

	var hilosNeutro int
	if sistema == entity.SistemaElectricoEstrella ||
		sistema == entity.SistemaElectricoBifasico ||
		sistema == entity.SistemaElectricoMonofasico {
		hilosNeutro = 1
	}

	totalHilos := hilosFaseTotal + hilosNeutro
	anchoRequerido := float64(totalHilos-1)*conductorFase.DiametroMM + conductorTierra.DiametroMM

	for _, entrada := range tablaCharola {
		if entrada.AnchoMM >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:      "CHAROLA_ESPACIADO",
				Tamano:    entrada.Tamano,
				AreaTotal: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"CalcularCharolaEspaciado: no se encontró charola suficiente: ancho requerido %.2f mm",
		anchoRequerido,
	)
}
