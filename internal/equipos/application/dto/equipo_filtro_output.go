// internal/equipos/application/dto/equipo_filtro_output.go
package dto

import (
	"time"

	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
)

// EquipoOutput is the outbound DTO for a single equipo filtro.
// All fields are primitives — no domain types exposed.
type EquipoOutput struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"` // ISO 8601
	Clave     *string `json:"clave"`
	Tipo      string  `json:"tipo"`
	Voltaje   int     `json:"voltaje"`
	Amperaje  int     `json:"amperaje"`
	ITM       int     `json:"itm"`
	Bornes    *int    `json:"bornes"`
	Conexion  *string `json:"conexion"` // nullable: "MONOFASICA" | "TRIFASICA"
}

// FromDomain converts a domain entity to an output DTO.
func FromDomain(e *entity.EquipoFiltro) EquipoOutput {
	var conexion *string
	if e.Conexion != nil {
		s := e.Conexion.String()
		conexion = &s
	}

	return EquipoOutput{
		ID:        e.ID.String(),
		CreatedAt: e.CreatedAt.UTC().Format(time.RFC3339),
		Clave:     e.Clave,
		Tipo:      e.Tipo.String(),
		Voltaje:   e.Voltaje,
		Amperaje:  e.Amperaje,
		ITM:       e.ITM,
		Bornes:    e.Bornes,
		Conexion:  conexion,
	}
}

// PaginationMeta contains pagination metadata for collection responses.
type PaginationMeta struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ListEquiposOutput is the outbound DTO for a paginated list of equipos.
type ListEquiposOutput struct {
	Equipos    []EquipoOutput `json:"equipos"`
	Pagination PaginationMeta `json:"pagination"`
}

// FromDomainList converts a slice of domain entities to a paginated list output DTO.
// total is the FULL count (all matching rows, not just this page).
func FromDomainList(entities []*entity.EquipoFiltro, page, pageSize, total int) ListEquiposOutput {
	out := make([]EquipoOutput, len(entities))
	for i, e := range entities {
		out[i] = FromDomain(e)
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return ListEquiposOutput{
		Equipos: out,
		Pagination: PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}
}
