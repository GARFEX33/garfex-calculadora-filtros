// internal/equipos/application/usecase/obtener_equipo.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	"github.com/google/uuid"
)

// ObtenerEquipoUseCase handles retrieving a single equipo filtro by ID.
type ObtenerEquipoUseCase struct {
	repo port.EquipoFiltroRepository
}

// NewObtenerEquipoUseCase creates a new instance with the required repository.
func NewObtenerEquipoUseCase(repo port.EquipoFiltroRepository) *ObtenerEquipoUseCase {
	return &ObtenerEquipoUseCase{repo: repo}
}

// Execute parses the UUID string, fetches the entity, and converts to output DTO.
func (uc *ObtenerEquipoUseCase) Execute(ctx context.Context, id string) (dto.EquipoOutput, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("%w: %s", dto.ErrIDInvalido, id)
	}

	equipo, err := uc.repo.ObtenerPorID(ctx, parsedID)
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("obtener equipo: %w", err)
	}

	return dto.FromDomain(equipo), nil
}
