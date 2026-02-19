// internal/calculos/application/dto/charola_triangular.go
package dto

import "fmt"

// CharolaTriangularInput representa los datos de entrada para calcular charola con configuración triangular.
type CharolaTriangularInput struct {
	HilosPorFase      int      `json:"hilos_por_fase"`
	DiametroFaseMM    float64  `json:"diametro_fase_mm"`
	DiametroTierraMM  float64  `json:"diametro_tierra_mm"`
	DiametroControlMM *float64 `json:"diametro_control_mm,omitempty"`
}

// Validate valida los campos de entrada.
func (i CharolaTriangularInput) Validate() error {
	if i.HilosPorFase < 1 {
		return fmt.Errorf("hilos_por_fase debe ser >= 1")
	}
	if i.DiametroFaseMM <= 0 {
		return fmt.Errorf("diametro_fase_mm debe ser mayor que cero")
	}
	if i.DiametroTierraMM <= 0 {
		return fmt.Errorf("diametro_tierra_mm debe ser mayor que cero")
	}
	if i.DiametroControlMM != nil && *i.DiametroControlMM <= 0 {
		return fmt.Errorf("diametro_control_mm debe ser mayor que cero si se proporciona")
	}
	return nil
}

// CharolaTriangularOutput representa el resultado del cálculo de charola con configuración triangular.
type CharolaTriangularOutput struct {
	Tipo           string  `json:"tipo"`
	Tamano         string  `json:"tamano"`
	TamanoPulgadas string  `json:"tamano_pulgadas"`
	AnchoRequerido float64 `json:"ancho_requerido_mm"`
}
