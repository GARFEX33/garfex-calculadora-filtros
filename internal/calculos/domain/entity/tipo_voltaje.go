// internal/calculos/domain/entity/tipo_voltaje.go
package entity

import (
	"errors"
	"strings"
)

// TipoVoltaje representa el tipo de voltaje de referencia ingresado por el usuario.
//
// En sistemas eléctricos existen dos tipos de voltaje:
//   - Vfn (Voltaje Fase-Neutro): voltaje entre una fase y el neutro
//   - Vff (Voltaje Fase-Fase): voltaje entre dos fases (también llamado "línea a línea")
//
// Relación: Vff = √3 × Vfn (para sistemas trifásicos balanceados)
//
// Según NOM-001-SEDE-2012, el cálculo de caída de tensión requiere:
//   - Sistemas MONOFASICO, BIFASICO, ESTRELLA → usar Vfn como referencia
//   - Sistema DELTA → usar Vff como referencia
type TipoVoltaje string

const (
	// TipoVoltajeFaseNeutro indica que el voltaje ingresado es entre fase y neutro.
	// Ejemplos: 127V, 277V
	TipoVoltajeFaseNeutro TipoVoltaje = "FASE_NEUTRO"

	// TipoVoltajeFaseFase indica que el voltaje ingresado es entre fases.
	// Ejemplos: 220V, 480V
	// También conocido como "voltaje línea a línea" o "voltaje de línea".
	TipoVoltajeFaseFase TipoVoltaje = "FASE_FASE"
)

// ErrTipoVoltajeInvalido se retorna cuando el tipo de voltaje no es reconocido.
var ErrTipoVoltajeInvalido = errors.New("tipo de voltaje inválido: debe ser 'FASE_NEUTRO' o 'FASE_FASE'")

// ParseTipoVoltaje convierte un string a TipoVoltaje.
//
// Acepta (case-insensitive):
//   - "FASE_NEUTRO", "fase_neutro", "FN"
//   - "FASE_FASE", "fase_fase", "FF"
func ParseTipoVoltaje(s string) (TipoVoltaje, error) {
	upper := strings.ToUpper(strings.TrimSpace(s))

	switch upper {
	case "FASE_NEUTRO", "FN":
		return TipoVoltajeFaseNeutro, nil
	case "FASE_FASE", "FF":
		return TipoVoltajeFaseFase, nil
	default:
		return "", ErrTipoVoltajeInvalido
	}
}

// String retorna la representación en string del tipo de voltaje.
func (t TipoVoltaje) String() string {
	return string(t)
}

// EsFaseNeutro retorna true si el tipo es fase-neutro.
func (t TipoVoltaje) EsFaseNeutro() bool {
	return t == TipoVoltajeFaseNeutro
}

// EsFaseFase retorna true si el tipo es fase-fase.
func (t TipoVoltaje) EsFaseFase() bool {
	return t == TipoVoltajeFaseFase
}
