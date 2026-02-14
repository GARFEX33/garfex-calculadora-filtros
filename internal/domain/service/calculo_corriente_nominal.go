// internal/domain/service/calculo_corriente_nominal.go
package service

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// CalcularCorrienteNominal delegates to the equipment's own calculation method.
// Each equipment type implements CalculadorCorriente differently:
//   - FiltroActivo: I = AmperajeNominal (direct)
//   - FiltroRechazo: I = KVAR / (KV × √3)
//   - Transformador: I = KVA / (KV × √3)
//   - Carga: I = KW / (KV × factor × FP)
func CalcularCorrienteNominal(equipo entity.CalculadorCorriente) (valueobject.Corriente, error) {
	return equipo.CalcularCorrienteNominal()
}
