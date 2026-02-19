// internal/calculos/application/dto/tuberia_input.go
package dto

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
)

// TuberiaInput contiene los datos necesarios para calcular el tamaño de tubería.
// Es el DTO de entrada para el use case CalcularTamanioTuberia.
type TuberiaInput struct {
	NumFases         int    `json:"num_fases" binding:"required,gt=0"`
	CalibreFase      string `json:"calibre_fase" binding:"required"`
	NumNeutros       int    `json:"num_neutros" binding:"gte=0"`
	CalibreNeutro    string `json:"calibre_neutral"`
	CalibreTierra    string `json:"calibre_tierra" binding:"required"`
	TipoCanalizacion string `json:"tipo_canalizacion" binding:"required"`
	NumTuberias      int    `json:"num_tuberias" binding:"required,gt=0"`
}

// Validate verifica que el input tenga los campos requeridos.
func (t TuberiaInput) Validate() error {
	if t.NumFases <= 0 {
		return fmt.Errorf("num_fases debe ser mayor que cero")
	}
	if t.CalibreFase == "" {
		return fmt.Errorf("calibre_fase es requerido")
	}
	if t.NumNeutros < 0 {
		return fmt.Errorf("num_neutros no puede ser negativo")
	}
	// Si hay neutros, el calibre es requerido
	if t.NumNeutros > 0 && t.CalibreNeutro == "" {
		return fmt.Errorf("calibre_neutral es requerido cuando num_neutros > 0")
	}
	if t.CalibreTierra == "" {
		return fmt.Errorf("calibre_tierra es requerido")
	}
	if t.TipoCanalizacion == "" {
		return fmt.Errorf("tipo_canalizacion es requerido")
	}
	if t.NumTuberias <= 0 {
		return fmt.Errorf("num_tuberias debe ser mayor que cero")
	}
	return nil
}

// ToDomainTipoCanalizacion convierte el string a entity.TipoCanalizacion.
func (t TuberiaInput) ToDomainTipoCanalizacion() (entity.TipoCanalizacion, error) {
	return entity.ParseTipoCanalizacion(t.TipoCanalizacion)
}
