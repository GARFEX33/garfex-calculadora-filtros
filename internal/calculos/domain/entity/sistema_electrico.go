// internal/calculos/domain/entity/sistema_electrico.go
package entity

import "errors"

// SistemaElectrico represents the type of electrical system configuration.
// It determines the number of conductors needed for the installation:
//   - Delta: 3-phase, 3-wire system (3 conductors)
//   - Estrella (Wye): 3-phase, 4-wire system (4 conductors)
//   - Bifasico: 2-phase system (3 conductors)
//   - Monofasico: single-phase system (2 conductors)
type SistemaElectrico string

const (
	SistemaElectricoDelta      SistemaElectrico = "DELTA"
	SistemaElectricoEstrella   SistemaElectrico = "ESTRELLA"
	SistemaElectricoBifasico   SistemaElectrico = "BIFASICO"
	SistemaElectricoMonofasico SistemaElectrico = "MONOFASICO"
)

var ErrSistemaElectricoInvalido = errors.New("sistema eléctrico no válido")

// ParseSistemaElectrico converts a string to SistemaElectrico.
func ParseSistemaElectrico(s string) (SistemaElectrico, error) {
	switch s {
	case string(SistemaElectricoDelta):
		return SistemaElectricoDelta, nil
	case string(SistemaElectricoEstrella):
		return SistemaElectricoEstrella, nil
	case string(SistemaElectricoBifasico):
		return SistemaElectricoBifasico, nil
	case string(SistemaElectricoMonofasico):
		return SistemaElectricoMonofasico, nil
	default:
		return "", errors.New("sistema eléctrico no válido: " + s)
	}
}

// ValidarSistemaElectrico returns an error if se is not a recognized electrical system type.
func ValidarSistemaElectrico(se SistemaElectrico) error {
	switch se {
	case SistemaElectricoDelta,
		SistemaElectricoEstrella,
		SistemaElectricoBifasico,
		SistemaElectricoMonofasico:
		return nil
	default:
		return ErrSistemaElectricoInvalido
	}
}

// CantidadConductores returns the number of conductors required for the electrical system.
func (s SistemaElectrico) CantidadConductores() int {
	switch s {
	case SistemaElectricoDelta:
		return 3
	case SistemaElectricoEstrella:
		return 4
	case SistemaElectricoBifasico:
		return 3
	case SistemaElectricoMonofasico:
		return 2
	default:
		return 0
	}
}
