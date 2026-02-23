// internal/equipos/application/usecase/eliminar_equipo.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	"github.com/google/uuid"
)

// EliminarEquipoUseCase handles deleting an equipo filtro by ID.
// The operation is idempotent — no error if the equipo does not exist.
type EliminarEquipoUseCase struct {
	repo port.EquipoFiltroRepository
}

// NewEliminarEquipoUseCase creates a new instance with the required repository.
func NewEliminarEquipoUseCase(repo port.EquipoFiltroRepository) *EliminarEquipoUseCase {
	return &EliminarEquipoUseCase{repo: repo}
}

// Execute parses the UUID and delegates deletion to the repository.
func (uc *EliminarEquipoUseCase) Execute(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: %s", dto.ErrIDInvalido, id)
	}

	if err := uc.repo.Eliminar(ctx, parsedID); err != nil {
		return fmt.Errorf("eliminar equipo: %w", err)
	}

	return nil
}
