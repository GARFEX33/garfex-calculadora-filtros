package entity

import "fmt"

type SistemaElectrico string

const (
	SistemaElectricoDelta      SistemaElectrico = "DELTA"
	SistemaElectricoEstrella   SistemaElectrico = "ESTRELLA"
	SistemaElectricoBifasico   SistemaElectrico = "BIFASICO"
	SistemaElectricoMonofasico SistemaElectrico = "MONOFASICO"
)

var ErrSistemaElectricoInvalido = fmt.Errorf("sistema eléctrico no válido")

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
		return "", fmt.Errorf("%w: '%s'", ErrSistemaElectricoInvalido, s)
	}
}

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
