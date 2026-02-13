package service

import (
	"errors"
	"fmt"
)

var ErrCantidadConductoresInvalida = errors.New("cantidad de conductores debe ser mayor que cero")

type EntradaTablaFactorAgrupamiento struct {
	CantidadMin int
	CantidadMax int
	Factor      float64
}

func CalcularFactorAgrupamiento(
	cantidad int,
	tabla []EntradaTablaFactorAgrupamiento,
) (float64, error) {
	if cantidad <= 0 {
		return 0, fmt.Errorf("CalcularFactorAgrupamiento: %w: %d", ErrCantidadConductoresInvalida, cantidad)
	}

	for _, entrada := range tabla {
		if entrada.CantidadMax == -1 {
			if cantidad >= entrada.CantidadMin {
				return entrada.Factor, nil
			}
		} else {
			if cantidad >= entrada.CantidadMin && cantidad <= entrada.CantidadMax {
				return entrada.Factor, nil
			}
		}
	}

	return 0.30, nil
}
