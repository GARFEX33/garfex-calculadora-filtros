// internal/equipos/application/dto/equipo_filtro_input.go
package dto

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
)

// CreateEquipoInput is the inbound DTO for creating a new equipo filtro.
// All fields are primitives following the same convention as calculos DTOs.
type CreateEquipoInput struct {
	Clave    *string `json:"clave"`
	Tipo     string  `json:"tipo"`
	Voltaje  int     `json:"voltaje"`
	Amperaje int     `json:"amperaje"`
	ITM      int     `json:"itm"`
	Bornes   *int    `json:"bornes"`
	Conexion *string `json:"conexion"` // nullable: "MONOFASICA" | "TRIFASICA"
}

// Validate checks that all required fields are present and valid.
func (i CreateEquipoInput) Validate() error {
	if i.Tipo == "" {
		return fmt.Errorf("%w: tipo es requerido", ErrInputInvalido)
	}
	if _, err := entity.ParseTipoFiltro(i.Tipo); err != nil {
		return fmt.Errorf("%w: %s", ErrInputInvalido, err.Error())
	}
	if i.Voltaje <= 0 {
		return fmt.Errorf("%w: voltaje debe ser mayor que cero", ErrInputInvalido)
	}
	if i.Amperaje <= 0 {
		return fmt.Errorf("%w: amperaje debe ser mayor que cero", ErrInputInvalido)
	}
	if i.ITM <= 0 {
		return fmt.Errorf("%w: itm debe ser mayor que cero", ErrInputInvalido)
	}
	if i.Conexion != nil {
		if _, err := entity.ParseConexion(*i.Conexion); err != nil {
			return fmt.Errorf("%w: %s", ErrInputInvalido, err.Error())
		}
	}
	return nil
}

// ToDomain converts the DTO to a domain entity ready for persistence.
func (i CreateEquipoInput) ToDomain() (*entity.EquipoFiltro, error) {
	tipo, err := entity.ParseTipoFiltro(i.Tipo)
	if err != nil {
		return nil, err
	}

	var conexion *entity.Conexion
	if i.Conexion != nil {
		c, err := entity.ParseConexion(*i.Conexion)
		if err != nil {
			return nil, err
		}
		conexion = &c
	}

	return entity.NewEquipoFiltro(i.Clave, tipo, i.Voltaje, i.Amperaje, i.ITM, i.Bornes, conexion)
}

// UpdateEquipoInput is the inbound DTO for updating an existing equipo filtro.
// The ID comes from the URL path, not the body.
type UpdateEquipoInput struct {
	Clave    *string `json:"clave"`
	Tipo     string  `json:"tipo"`
	Voltaje  int     `json:"voltaje"`
	Amperaje int     `json:"amperaje"`
	ITM      int     `json:"itm"`
	Bornes   *int    `json:"bornes"`
	Conexion *string `json:"conexion"` // nullable: "MONOFASICA" | "TRIFASICA"
}

// Validate checks that all required fields are present and valid.
func (i UpdateEquipoInput) Validate() error {
	if i.Tipo == "" {
		return fmt.Errorf("%w: tipo es requerido", ErrInputInvalido)
	}
	if _, err := entity.ParseTipoFiltro(i.Tipo); err != nil {
		return fmt.Errorf("%w: %s", ErrInputInvalido, err.Error())
	}
	if i.Voltaje <= 0 {
		return fmt.Errorf("%w: voltaje debe ser mayor que cero", ErrInputInvalido)
	}
	if i.Amperaje <= 0 {
		return fmt.Errorf("%w: amperaje debe ser mayor que cero", ErrInputInvalido)
	}
	if i.ITM <= 0 {
		return fmt.Errorf("%w: itm debe ser mayor que cero", ErrInputInvalido)
	}
	if i.Conexion != nil {
		if _, err := entity.ParseConexion(*i.Conexion); err != nil {
			return fmt.Errorf("%w: %s", ErrInputInvalido, err.Error())
		}
	}
	return nil
}

// ListEquiposQuery contains optional query filters and pagination for listing equipos.
type ListEquiposQuery struct {
	Tipo     string // optional: "A" | "KVA" | "KVAR"
	Voltaje  int    // optional: > 0 to filter by voltage
	Page     int    // 1-indexed, default 1
	PageSize int    // default 20, max 100
}

// ApplyDefaults sets sensible defaults for pagination fields.
func (q *ListEquiposQuery) ApplyDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}

// Offset returns the SQL OFFSET value for this page.
func (q ListEquiposQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}
