// internal/equipos/application/port/equipo_filtro_repository.go
package port

import (
	"context"

	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/google/uuid"
)

// FiltrosListado contains optional filters and pagination for listing equipment.
type FiltrosListado struct {
	Tipo    *entity.TipoFiltro // nil = all types
	Buscar  *string            // nil = no search, otherwise ILIKE on clave
	Voltaje *int               // nil = all voltages
	Limit   int                // page size (> 0)
	Offset  int                // number of rows to skip
}

// EquipoFiltroRepository defines the persistence contract for filter equipment.
// Infrastructure must implement this interface.
type EquipoFiltroRepository interface {
	// Crear persists a new equipo and returns it with the DB-generated ID and CreatedAt.
	Crear(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error)

	// ObtenerPorID finds an equipo by its UUID. Returns ErrEquipoNoEncontrado if missing.
	ObtenerPorID(ctx context.Context, id uuid.UUID) (*entity.EquipoFiltro, error)

	// Listar returns a paginated page of equipos matching the optional filters.
	Listar(ctx context.Context, filtros FiltrosListado) ([]*entity.EquipoFiltro, error)

	// Contar returns the total count of equipos matching the filters (ignoring pagination).
	Contar(ctx context.Context, filtros FiltrosListado) (int, error)

	// Actualizar updates an existing equipo and returns the updated record.
	// Returns ErrEquipoNoEncontrado if the ID does not exist.
	Actualizar(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error)

	// Eliminar deletes an equipo by UUID. Idempotent — no error if not found.
	Eliminar(ctx context.Context, id uuid.UUID) error
}
