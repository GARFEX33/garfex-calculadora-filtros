// internal/calculos/domain/service/ajuste_corriente.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ResultadoAjusteCorriente contains the result of the current adjustment.
type ResultadoAjusteCorriente struct {
	CorrienteAjustada valueobject.Corriente
	FactorTotal       float64
}

// AjustarCorriente multiplies the nominal current by all adjustment factors
// (protection, temperature, grouping, etc.). Returns error if any factor <= 0.
func AjustarCorriente(cn valueobject.Corriente, factores map[string]float64) (ResultadoAjusteCorriente, error) {
	if len(factores) == 0 {
		return ResultadoAjusteCorriente{
			CorrienteAjustada: cn,
			FactorTotal:       1.0,
		}, nil
	}

	resultado := cn.Valor()
	factorTotal := 1.0

	for nombre, factor := range factores {
		if factor <= 0 {
			return ResultadoAjusteCorriente{}, fmt.Errorf("factor '%s' inválido: %.4f (debe ser > 0)", nombre, factor)
		}
		resultado *= factor
		factorTotal *= factor
	}

	corrienteAjustada, err := valueobject.NewCorriente(resultado)
	if err != nil {
		return ResultadoAjusteCorriente{}, fmt.Errorf("corriente ajustada inválida: %w", err)
	}

	return ResultadoAjusteCorriente{
		CorrienteAjustada: corrienteAjustada,
		FactorTotal:       factorTotal,
	}, nil
}
