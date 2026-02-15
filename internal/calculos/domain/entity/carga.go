// internal/calculos/domain/entity/carga.go
package entity

import (
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// Carga represents a generic electrical load.
// Nominal current depends on the number of phases:
//
//	3-phase: I = KW / (KV × √3 × FP)
//	2-phase: I = KW / (KV × 2 × FP)
//	1-phase: I = KW / (KV × FP)
type Carga struct {
	Equipo
	KW             int
	FactorPotencia float64
	Fases          int // 1, 2, or 3
}

func NewCarga(clave string, voltaje, kw int, fp float64, fases int, itm ITM) (*Carga, error) {
	if kw <= 0 {
		return nil, fmt.Errorf("KW debe ser mayor que cero: %d", kw)
	}
	if fp <= 0 || fp > 1 {
		return nil, fmt.Errorf("factor de potencia debe estar entre 0 (exclusivo) y 1: %f", fp)
	}
	if fases < 1 || fases > 3 {
		return nil, fmt.Errorf("fases debe ser 1, 2 o 3: %d", fases)
	}
	if voltaje <= 0 {
		return nil, fmt.Errorf("%w: voltaje es %d", ErrDivisionPorCero, voltaje)
	}
	return &Carga{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoEquipoCarga,
			Voltaje: voltaje,
			ITM:     itm,
		},
		KW:             kw,
		FactorPotencia: fp,
		Fases:          fases,
	}, nil
}

// factorFases returns the phase multiplier for current calculation.
func (c *Carga) factorFases() float64 {
	switch c.Fases {
	case 1:
		return 1
	case 2:
		return 2
	default:
		return math.Sqrt(3)
	}
}

func (c *Carga) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	kv := float64(c.Voltaje) / 1000.0
	denominador := kv * c.factorFases() * c.FactorPotencia
	corriente := float64(c.KW) / denominador
	return valueobject.NewCorriente(corriente)
}

// PotenciaKW returns active power directly from KW rating.
func (c *Carga) PotenciaKW() float64 {
	return float64(c.KW)
}

// PotenciaKVA returns apparent power: KW / FP.
func (c *Carga) PotenciaKVA() float64 {
	return float64(c.KW) / c.FactorPotencia
}

// PotenciaKVAR returns reactive power: √(KVA² - KW²).
func (c *Carga) PotenciaKVAR() float64 {
	kva := c.PotenciaKVA()
	kw := float64(c.KW)
	return math.Sqrt(kva*kva - kw*kw)
}
