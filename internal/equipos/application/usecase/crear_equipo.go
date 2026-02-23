// internal/equipos/application/usecase/crear_equipo.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
)

// CrearEquipoUseCase handles the creation of a new equipo filtro.
type CrearEquipoUseCase struct {
	repo port.EquipoFiltroRepository
}

// NewCrearEquipoUseCase creates a new instance with the required repository.
func NewCrearEquipoUseCase(repo port.EquipoFiltroRepository) *CrearEquipoUseCase {
	return &CrearEquipoUseCase{repo: repo}
}

// Execute validates the input, creates the domain entity, and persists it.
// Returns the persisted equipo with DB-generated ID and CreatedAt.
func (uc *CrearEquipoUseCase) Execute(ctx context.Context, input dto.CreateEquipoInput) (dto.EquipoOutput, error) {
	if err := input.Validate(); err != nil {
		return dto.EquipoOutput{}, err
	}

	equipo, err := input.ToDomain()
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("convertir input a domain: %w", err)
	}

	created, err := uc.repo.Crear(ctx, equipo)
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("crear equipo: %w", err)
	}

	return dto.FromDomain(created), nil
}
