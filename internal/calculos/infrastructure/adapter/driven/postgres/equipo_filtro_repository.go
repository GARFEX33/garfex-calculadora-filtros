// internal/calculos/infrastructure/adapter/driven/postgres/equipo_filtro_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	calcport "github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	calcent "github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	equipoent "github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CalcEquipoFiltroRepository implements calculos/port.EquipoRepository.
// It queries the equipos_filtros table and maps each TipoFiltro to the
// correct calculos domain entity (FiltroActivo, FiltroRechazo, Transformador).
type CalcEquipoFiltroRepository struct {
	pool *pgxpool.Pool
}

// NewCalcEquipoFiltroRepository creates a new instance sharing the given pool.
func NewCalcEquipoFiltroRepository(pool *pgxpool.Pool) *CalcEquipoFiltroRepository {
	return &CalcEquipoFiltroRepository{pool: pool}
}

// Compile-time check: CalcEquipoFiltroRepository must implement calcport.EquipoRepository.
var _ calcport.EquipoRepository = (*CalcEquipoFiltroRepository)(nil)

// BuscarPorClave finds an equipment by its unique clave and returns the
// correct calculos domain entity based on TipoFiltro:
//
//   - TipoFiltroA    → FiltroActivo  (I = Amperaje directo)
//   - TipoFiltroKVA  → Transformador (I = KVA / (kV × √3))
//   - TipoFiltroKVAR → FiltroRechazo (I = KVAR / (kV × √3))
func (r *CalcEquipoFiltroRepository) BuscarPorClave(ctx context.Context, clave string) (calcent.CalculadorCorriente, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT clave, tipo, voltaje, "qn/In", itm, bornes
		FROM equipos_filtros
		WHERE clave = $1
		LIMIT 1
	`

	row := r.pool.QueryRow(ctx, query, clave)

	var (
		claveDB  *string
		tipoStr  string
		voltaje  int
		amperaje int
		itmVal   int
		bornes   *int
	)

	if err := row.Scan(&claveDB, &tipoStr, &voltaje, &amperaje, &itmVal, &bornes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("equipo no encontrado con clave '%s'", clave)
		}
		return nil, fmt.Errorf("buscar equipo por clave: %w", err)
	}

	tipo, err := equipoent.ParseTipoFiltro(tipoStr)
	if err != nil {
		return nil, fmt.Errorf("tipo de filtro inválido en BD '%s': %w", tipoStr, err)
	}

	claveStr := clave
	if claveDB != nil {
		claveStr = *claveDB
	}

	// Construir ITM con defaults para polos y bornes según NOM
	// (la BD solo guarda amperaje del ITM; polos=3 es estándar trifásico)
	bornesITM := 3
	if bornes != nil && *bornes > 0 {
		bornesITM = *bornes
	}
	itm, err := calcent.NewITM(itmVal, 3, bornesITM, voltaje)
	if err != nil {
		return nil, fmt.Errorf("construir ITM para equipo '%s': %w", claveStr, err)
	}

	return mapToCalculadorCorriente(claveStr, tipo, voltaje, amperaje, itm)
}

// mapToCalculadorCorriente converts a BD row to the correct calculos domain entity.
//
//   - A    → FiltroActivo:  campo amperaje = corriente nominal directa (A)
//   - KVA  → Transformador: campo amperaje = potencia aparente (KVA)
//   - KVAR → FiltroRechazo: campo amperaje = potencia reactiva (KVAR)
func mapToCalculadorCorriente(
	clave string,
	tipo equipoent.TipoFiltro,
	voltaje int,
	amperaje int,
	itm calcent.ITM,
) (calcent.CalculadorCorriente, error) {
	switch tipo {
	case equipoent.TipoFiltroA:
		// Amperaje = corriente nominal directa → FiltroActivo
		return calcent.NewFiltroActivo(clave, voltaje, amperaje, itm)

	case equipoent.TipoFiltroKVA:
		// Amperaje = KVA → Transformador: I = KVA / (kV × √3)
		return calcent.NewTransformador(clave, voltaje, amperaje, itm)

	case equipoent.TipoFiltroKVAR:
		// Amperaje = KVAR → FiltroRechazo: I = KVAR / (kV × √3)
		return calcent.NewFiltroRechazo(clave, voltaje, amperaje, itm)

	default:
		return nil, fmt.Errorf("tipo de filtro no soportado para cálculo: '%s'", tipo)
	}
}
