// internal/equipos/domain/entity/equipo_filtro.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

// EquipoFiltro represents an electrical filter in the catalog.
// Maps to the equipos_filtros table in PostgreSQL.
type EquipoFiltro struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Clave       *string // nullable — unique key for user-readable identification
	Tipo        TipoFiltro
	Voltaje     int          // nominal voltage in volts
	Amperaje    int          // nominal current Qn/In in amps
	ITM         int          // interruptor termomagnético capacity in amps
	Bornes      *int         // nullable — number of terminals
	Conexion    *Conexion    // nullable — electrical connection type
	TipoVoltaje *TipoVoltaje // nullable — voltage reference type (FF or FN); DB default: 'FF'
}

// NewEquipoFiltro creates and validates a new EquipoFiltro entity.
// ID and CreatedAt are set by PostgreSQL on insert; they are zero here.
// tipoVoltaje is optional (nullable); when nil the DB default ('FF') applies on insert.
func NewEquipoFiltro(
	clave *string,
	tipo TipoFiltro,
	voltaje int,
	amperaje int,
	itm int,
	bornes *int,
	conexion *Conexion,
	tipoVoltaje *TipoVoltaje,
) (*EquipoFiltro, error) {
	if voltaje <= 0 {
		return nil, ErrVoltajeInvalido
	}
	if amperaje <= 0 {
		return nil, ErrAmperajeInvalido
	}
	if itm <= 0 {
		return nil, ErrITMInvalido
	}

	return &EquipoFiltro{
		Clave:       clave,
		Tipo:        tipo,
		Voltaje:     voltaje,
		Amperaje:    amperaje,
		ITM:         itm,
		Bornes:      bornes,
		Conexion:    conexion,
		TipoVoltaje: tipoVoltaje,
	}, nil
}
