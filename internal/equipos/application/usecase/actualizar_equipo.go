// internal/equipos/application/usecase/actualizar_equipo.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/google/uuid"
)

// ActualizarEquipoUseCase handles updating an existing equipo filtro.
type ActualizarEquipoUseCase struct {
	repo port.EquipoFiltroRepository
}

// NewActualizarEquipoUseCase creates a new instance with the required repository.
func NewActualizarEquipoUseCase(repo port.EquipoFiltroRepository) *ActualizarEquipoUseCase {
	return &ActualizarEquipoUseCase{repo: repo}
}

// Execute validates input, builds the updated entity, and persists changes.
func (uc *ActualizarEquipoUseCase) Execute(ctx context.Context, id string, input dto.UpdateEquipoInput) (dto.EquipoOutput, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("%w: %s", dto.ErrIDInvalido, id)
	}

	if err := input.Validate(); err != nil {
		return dto.EquipoOutput{}, err
	}

	tipo, err := entity.ParseTipoFiltro(input.Tipo)
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("%w: %s", dto.ErrInputInvalido, err.Error())
	}

	var conexion *entity.Conexion
	if input.Conexion != nil {
		c, err := entity.ParseConexion(*input.Conexion)
		if err != nil {
			return dto.EquipoOutput{}, fmt.Errorf("%w: %s", dto.ErrInputInvalido, err.Error())
		}
		conexion = &c
	}

	equipo := &entity.EquipoFiltro{
		ID:       parsedID,
		Clave:    input.Clave,
		Tipo:     tipo,
		Voltaje:  input.Voltaje,
		Amperaje: input.Amperaje,
		ITM:      input.ITM,
		Bornes:   input.Bornes,
		Conexion: conexion,
	}

	updated, err := uc.repo.Actualizar(ctx, equipo)
	if err != nil {
		return dto.EquipoOutput{}, fmt.Errorf("actualizar equipo: %w", err)
	}

	return dto.FromDomain(updated), nil
}
