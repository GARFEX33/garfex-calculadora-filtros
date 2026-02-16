// internal/calculos/domain/service/calcular_factor_uso.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
)

// CalcularFactorUso retorna el factor de uso según el tipo de equipo.
//
// Factores según normativa NOM:
//   - FILTRO_ACTIVO, FILTRO_RECHAZO → 1.35
//   - TRANSFORMADOR, CARGA → 1.25
//
// Retorna error si el tipo de equipo no es válido.
func CalcularFactorUso(tipoEquipo entity.TipoEquipo) (float64, error) {
	switch tipoEquipo {
	case entity.TipoEquipoFiltroActivo, entity.TipoEquipoFiltroRechazo:
		return 1.35, nil
	case entity.TipoEquipoTransformador, entity.TipoEquipoCarga:
		return 1.25, nil
	default:
		return 0, fmt.Errorf("CalcularFactorUso: %w: '%s'", entity.ErrTipoEquipoInvalido, tipoEquipo)
	}
}
