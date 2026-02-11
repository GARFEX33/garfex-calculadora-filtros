// internal/domain/service/ajuste_corriente.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// AjustarCorriente multiplies the nominal current by all adjustment factors
// (protection, temperature, grouping, etc.). Returns error if any factor <= 0.
func AjustarCorriente(cn valueobject.Corriente, factores map[string]float64) (valueobject.Corriente, error) {
	if len(factores) == 0 {
		return cn, nil
	}

	resultado := cn.Valor()
	for nombre, factor := range factores {
		if factor <= 0 {
			return valueobject.Corriente{}, fmt.Errorf("factor '%s' inválido: %.4f (debe ser > 0)", nombre, factor)
		}
		resultado *= factor
	}

	return valueobject.NewCorriente(resultado)
}
