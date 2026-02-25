// internal/calculos/application/dto/charola_espaciado.go
package dto

import "fmt"

// CharolaEspaciadoInput representa los datos de entrada para calcular charola con espaciado.
type CharolaEspaciadoInput struct {
	HilosPorFase      int      `json:"hilos_por_fase"`
	SistemaElectrico  string   `json:"sistema_electrico"`
	DiametroFaseMM    float64  `json:"diametro_fase_mm"`
	DiametroTierraMM  float64  `json:"diametro_tierra_mm"`
	DiametroControlMM *float64 `json:"diametro_control_mm,omitempty"`
}

// Validate valida los campos de entrada.
func (i CharolaEspaciadoInput) Validate() error {
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

// CharolaEspaciadoOutput representa el resultado del cálculo de charola con espaciado.
type CharolaEspaciadoOutput struct {
	Tipo             string  `json:"tipo"`
	Tamano           string  `json:"tamano"`
	TamanoPulgadas   string  `json:"tamano_pulgadas"`
	AnchoRequerido   float64 `json:"ancho_requerido_mm"`
	AnchoComercialMM float64 `json:"ancho_comercial_mm"`

	// Valores intermedios del desarrollo — para mostrar en memoria de cálculo
	DiametroFaseMM    float64  `json:"diametro_fase_mm"`
	DiametroTierraMM  float64  `json:"diametro_tierra_mm"`
	DiametroControlMM *float64 `json:"diametro_control_mm,omitempty"`
	NumHilosTotal     int      `json:"num_hilos_total"`
	EspacioFuerzaMM   float64  `json:"espacio_fuerza_mm"`
	AnchoFuerzaMM     float64  `json:"ancho_fuerza_mm"`
	EspacioControlMM  float64  `json:"espacio_control_mm,omitempty"`
	AnchoControlMM    float64  `json:"ancho_control_mm,omitempty"`
	AnchoTierraMM     float64  `json:"ancho_tierra_mm"`
}
