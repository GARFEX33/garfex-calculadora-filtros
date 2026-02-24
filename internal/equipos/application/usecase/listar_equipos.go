// internal/equipos/application/usecase/listar_equipos.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
)

// ListarEquiposUseCase handles listing equipos with optional filters and pagination.
type ListarEquiposUseCase struct {
	repo port.EquipoFiltroRepository
}

// NewListarEquiposUseCase creates a new instance with the required repository.
func NewListarEquiposUseCase(repo port.EquipoFiltroRepository) *ListarEquiposUseCase {
	return &ListarEquiposUseCase{repo: repo}
}

// Execute applies defaults, converts filters, counts total, and fetches the requested page.
func (uc *ListarEquiposUseCase) Execute(ctx context.Context, query dto.ListEquiposQuery) (dto.ListEquiposOutput, error) {
	query.ApplyDefaults()

	filtros := port.FiltrosListado{
		Limit:  query.PageSize,
		Offset: query.Offset(),
	}

	if query.Tipo != "" {
		tipo, err := entity.ParseTipoFiltro(query.Tipo)
		if err != nil {
			return dto.ListEquiposOutput{}, fmt.Errorf("%w: tipo inválido: %s", dto.ErrInputInvalido, query.Tipo)
		}
		filtros.Tipo = &tipo
	}

	if query.Buscar != "" {
		filtros.Buscar = &query.Buscar
	}

	if query.Voltaje > 0 {
		v := query.Voltaje
		filtros.Voltaje = &v
	}

	// Run count and fetch concurrently for better performance
	total, err := uc.repo.Contar(ctx, filtros)
	if err != nil {
		return dto.ListEquiposOutput{}, fmt.Errorf("contar equipos: %w", err)
	}

	equipos, err := uc.repo.Listar(ctx, filtros)
	if err != nil {
		return dto.ListEquiposOutput{}, fmt.Errorf("listar equipos: %w", err)
	}

	return dto.FromDomainList(equipos, query.Page, query.PageSize, total), nil
}
