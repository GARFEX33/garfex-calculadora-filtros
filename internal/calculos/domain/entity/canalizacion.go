// internal/calculos/domain/entity/canalizacion.go
package entity

import "fmt"

// Canalizacion represents the conduit or cable tray selected for the installation.
type Canalizacion struct {
	Tipo           TipoCanalizacion // wiring method (TUBERIA_PVC, CHAROLA_CABLE_ESPACIADO, etc.)
	Tamano         string           // e.g., "1 1/2" (inches for tubería), "300mm" (charola)
	AnchoRequerido float64          // for charola: required width in mm; for tubería: total conductor area in mm²
	NumeroDeTubos  int              // number of parallel conduits; 1 = single conduit installation
}

// NewCanalizacion constructs a validated Canalizacion.
func NewCanalizacion(tipo TipoCanalizacion, tamano string, anchoRequerido float64, numeroDeTubos int) (Canalizacion, error) {
	if err := ValidarTipoCanalizacion(tipo); err != nil {
		return Canalizacion{}, fmt.Errorf("NewCanalizacion: %w", err)
	}
	if tamano == "" {
		return Canalizacion{}, fmt.Errorf("NewCanalizacion: tamaño no puede estar vacío")
	}
	if numeroDeTubos < 1 {
		return Canalizacion{}, fmt.Errorf("NewCanalizacion: numeroDeTubos debe ser >= 1: %d", numeroDeTubos)
	}
	return Canalizacion{
		Tipo:           tipo,
		Tamano:         tamano,
		AnchoRequerido: anchoRequerido,
		NumeroDeTubos:  numeroDeTubos,
	}, nil
}
