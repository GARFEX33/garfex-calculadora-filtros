// internal/presentation/formatters/observaciones.go
package formatters

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
)

// GenerarObservaciones genera observaciones sobre el cálculo.
func GenerarObservaciones(output dto.MemoriaOutput) []string {
	var obs []string

	if !output.CaidaTension.Cumple {
		obs = append(obs, fmt.Sprintf(
			"Caída de tensión %.2f%% excede el límite de %.2f%%",
			output.CaidaTension.Porcentaje,
			output.CaidaTension.LimitePorcentaje,
		))
	}

	if output.HilosPorFase > 1 {
		obs = append(obs, fmt.Sprintf(
			"Se usan %d hilos por fase en paralelo",
			output.HilosPorFase,
		))
	}

	return obs
}
