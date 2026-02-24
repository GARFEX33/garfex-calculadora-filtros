// internal/equipos/domain/entity/tipo_voltaje.go
package entity

import "fmt"

// TipoVoltaje identifies the voltage reference type for an electrical filter.
// Maps to the PostgreSQL enum public.tipo_voltaje.
// Values match exactly the DB enum.
//
// FF (Fase-Fase): line-to-line voltage — e.g., 220V, 480V.
// FN (Fase-Neutro): line-to-neutral voltage — e.g., 127V, 277V.
type TipoVoltaje string

const (
	TipoVoltajeFaseFase   TipoVoltaje = "FF" // Voltaje fase-fase (línea a línea)
	TipoVoltajeFaseNeutro TipoVoltaje = "FN" // Voltaje fase-neutro
)

// ParseTipoVoltaje converts a string (e.g., from the database or HTTP request)
// to a TipoVoltaje. Returns ErrTipoVoltajeInvalido if the value is not recognized.
func ParseTipoVoltaje(s string) (TipoVoltaje, error) {
	switch s {
	case string(TipoVoltajeFaseFase):
		return TipoVoltajeFaseFase, nil
	case string(TipoVoltajeFaseNeutro):
		return TipoVoltajeFaseNeutro, nil
	default:
		return "", fmt.Errorf("%w: '%s' — valores válidos: FF, FN", ErrTipoVoltajeInvalido, s)
	}
}

// String returns the string representation of TipoVoltaje.
func (t TipoVoltaje) String() string {
	return string(t)
}
