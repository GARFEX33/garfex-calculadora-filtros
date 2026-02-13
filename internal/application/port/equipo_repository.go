// internal/application/port/equipo_repository.go
package port

import (
	"context"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
)

// EquipoRepository defines the contract for equipment persistence.
type EquipoRepository interface {
	// BuscarPorClave finds an equipment by its unique key.
	BuscarPorClave(ctx context.Context, clave string) (entity.Equipo, error)
}
