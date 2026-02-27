// internal/calculos/infrastructure/adapter/driven/mock/calc_equipo_repository.go
package mock

import (
	"context"
	"fmt"

	calcport "github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	calcent "github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
)

// CalcEquipoMockRepository implements calcport.EquipoRepository using in-memory data.
// For use in MOCK_MODE when PostgreSQL is not available.
type CalcEquipoMockRepository struct {
	equipos map[string]calcent.CalculadorCorriente
}

// NewCalcEquipoMockRepository creates a new mock repository preloaded with
// the same 5 equipos as the equipos mock repository.
func NewCalcEquipoMockRepository() *CalcEquipoMockRepository {
	// ── 150 A / ITM 200
	itm200a, _ := calcent.NewITM(200, 3, 6, 220)
	itm200b, _ := calcent.NewITM(200, 3, 6, 480)
	// ── 200 A / ITM 250
	itm250a, _ := calcent.NewITM(250, 3, 6, 440)
	itm250b, _ := calcent.NewITM(250, 3, 6, 220)
	// ── 300 A / ITM 400
	itm400a, _ := calcent.NewITM(400, 3, 6, 440)
	itm400b, _ := calcent.NewITM(400, 3, 6, 480)
	// ── 400 A / ITM 500
	itm500a, _ := calcent.NewITM(500, 3, 6, 480)
	itm500b, _ := calcent.NewITM(500, 3, 6, 220)
	// ── 600 A / ITM 800
	itm800a, _ := calcent.NewITM(800, 3, 6, 440)
	itm800b, _ := calcent.NewITM(800, 3, 6, 440)

	fa1, _ := calcent.NewFiltroActivo("FA-220-150A", 220, 150, itm200a)
	fr1, _ := calcent.NewFiltroRechazo("FKVAR-480-150", 480, 150, itm200b)
	fa2, _ := calcent.NewFiltroActivo("FA-440-200A", 440, 200, itm250a)
	tr1, _ := calcent.NewTransformador("FKVA-220-200", 220, 200, itm250b)
	fr2, _ := calcent.NewFiltroRechazo("FKVAR-440-300", 440, 300, itm400a)
	fa3, _ := calcent.NewFiltroActivo("FA-480-300A", 480, 300, itm400b)
	tr2, _ := calcent.NewTransformador("FKVA-480-400", 480, 400, itm500a)
	fr3, _ := calcent.NewFiltroRechazo("FKVAR-220-400", 220, 400, itm500b)
	fa4, _ := calcent.NewFiltroActivo("FA-440-600A", 440, 600, itm800a)
	tr3, _ := calcent.NewTransformador("FKVA-440-600", 440, 600, itm800b)

	return &CalcEquipoMockRepository{
		equipos: map[string]calcent.CalculadorCorriente{
			"FA-220-150A":   fa1,
			"FKVAR-480-150": fr1,
			"FA-440-200A":   fa2,
			"FKVA-220-200":  tr1,
			"FKVAR-440-300": fr2,
			"FA-480-300A":   fa3,
			"FKVA-480-400":  tr2,
			"FKVAR-220-400": fr3,
			"FA-440-600A":   fa4,
			"FKVA-440-600":  tr3,
		},
	}
}

// Compile-time check: CalcEquipoMockRepository must implement calcport.EquipoRepository.
var _ calcport.EquipoRepository = (*CalcEquipoMockRepository)(nil)

// BuscarPorClave returns the mock equipment with the given clave.
func (r *CalcEquipoMockRepository) BuscarPorClave(_ context.Context, clave string) (calcent.CalculadorCorriente, error) {
	equipo, ok := r.equipos[clave]
	if !ok {
		return nil, fmt.Errorf("equipo no encontrado con clave '%s'", clave)
	}
	return equipo, nil
}
