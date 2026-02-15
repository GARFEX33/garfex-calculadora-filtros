// internal/calculos/application/port/equipo_repository.go
package port

import (
	"context"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
)

// EquipoRepository defines the contract for equipment persistence.
type EquipoRepository interface {
	// BuscarPorClave finds an equipment by its unique key.
	// Returns CalculadorCorriente interface which all equipment types implement.
	BuscarPorClave(ctx context.Context, clave string) (entity.CalculadorCorriente, error)
}
