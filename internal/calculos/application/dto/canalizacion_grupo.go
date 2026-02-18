// internal/calculos/application/dto/canalizacion_grupo.go
package dto

import (
	"errors"
	"fmt"
)

// ErrConductoresVacios es retornado cuando no hay conductores.
var ErrConductoresVacios = errors.New("la lista de conductores no puede estar vacía")

// ConductorGrupoInput representa un grupo de conductores idénticos (fases)
type ConductorGrupoInput struct {
	Cantidad   int     `json:"cantidad" binding:"required,min=1"`
	SeccionMM2 float64 `json:"seccion_mm2" binding:"required,gt=0"`
}

// CanalizacionGrupoInput es el DTO de entrada para calcular canalización de grupo
type CanalizacionGrupoInput struct {
	Conductores      []ConductorGrupoInput `json:"conductores" binding:"required,min=1,dive"`
	SeccionTierraMM2 float64               `json:"seccion_tierra_mm2" binding:"required,gt=0"`
	TipoCanalizacion string                `json:"tipo_canalizacion" binding:"required"`
	NumeroDeTubos    int                   `json:"numero_de_tubos" binding:"required,min=1"`
}

// Validate valida los campos de entrada.
func (c CanalizacionGrupoInput) Validate() error {
	if len(c.Conductores) == 0 {
		return fmt.Errorf("%w", ErrConductoresVacios)
	}
	if c.SeccionTierraMM2 <= 0 {
		return errors.New("seccion_tierra_mm2 debe ser mayor a cero")
	}
	if c.TipoCanalizacion == "" {
		return errors.New("tipo_canalizacion es requerido")
	}
	if c.NumeroDeTubos < 1 {
		return errors.New("numero_de_tubos debe ser mayor o igual a 1")
	}
	for i, conductor := range c.Conductores {
		if conductor.Cantidad < 1 {
			return fmt.Errorf("conductor[%d].cantidad debe ser mayor a cero", i)
		}
		if conductor.SeccionMM2 <= 0 {
			return fmt.Errorf("conductor[%d].seccion_mm2 debe ser mayor a cero", i)
		}
	}
	return nil
}

// CanalizacionGrupoOutput es el resultado del cálculo
type CanalizacionGrupoOutput struct {
	Tamano         string  `json:"tamano"`
	AreaTotalMM2   float64 `json:"area_total_mm2"`
	AreaPorTuboMM2 float64 `json:"area_por_tubo_mm2"`
	NumeroDeTubos  int     `json:"numero_de_tubos"`
	FactorRelleno  float64 `json:"factor_relleno"`
}
