package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// EntradaTablaFactorTemperatura representa una fila de la tabla NOM 310-15(b)(2)(a)
type EntradaTablaFactorTemperatura struct {
	RangoTempC string
	Factor60C  float64
	Factor75C  float64
	Factor90C  float64
}

// CalcularFactorTemperatura retorna el factor de corrección según temperatura ambiente y del conductor
func CalcularFactorTemperatura(
	tempAmbiente int,
	tempConductor valueobject.Temperatura,
	tabla []EntradaTablaFactorTemperatura,
) (float64, error) {
	if tempAmbiente < -10 {
		return 0, fmt.Errorf("temperatura ambiente inválida: %d°C", tempAmbiente)
	}

	for _, entrada := range tabla {
		if rangoContiene(entrada.RangoTempC, tempAmbiente) {
			switch tempConductor {
			case valueobject.Temp60:
				return entrada.Factor60C, nil
			case valueobject.Temp75:
				return entrada.Factor75C, nil
			case valueobject.Temp90:
				return entrada.Factor90C, nil
			default:
				return 0, fmt.Errorf("temperatura de conductor no soportada: %v", tempConductor)
			}
		}
	}

	return 0, fmt.Errorf("no se encontró factor para temperatura ambiente %d°C", tempAmbiente)
}

func rangoContiene(rango string, temp int) bool {
	var min, max int
	if _, err := fmt.Sscanf(rango, "%d-%d", &min, &max); err == nil {
		return temp >= min && temp <= max
	}
	if _, err := fmt.Sscanf(rango, "%d+", &min); err == nil {
		return temp >= min
	}
	return false
}
