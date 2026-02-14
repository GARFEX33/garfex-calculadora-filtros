// internal/domain/entity/transformador.go
package entity

import (
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// Transformador represents a power transformer.
// Nominal current: I = KVA / (KV × √3), same formula as FiltroRechazo.
type Transformador struct {
	Equipo
	KVA int
}

func NewTransformador(clave string, voltaje, kva int, itm ITM) (*Transformador, error) {
	if kva <= 0 {
		return nil, fmt.Errorf("KVA debe ser mayor que cero: %d", kva)
	}
	if voltaje <= 0 {
		return nil, fmt.Errorf("%w: voltaje es %d", ErrDivisionPorCero, voltaje)
	}
	return &Transformador{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoEquipoTransformador,
			Voltaje: voltaje,
			ITM:     itm,
		},
		KVA: kva,
	}, nil
}

func (tr *Transformador) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	kv := float64(tr.Voltaje) / 1000.0
	denominador := kv * math.Sqrt(3)
	corriente := float64(tr.KVA) / denominador
	return valueobject.NewCorriente(corriente)
}

// PotenciaKVA returns apparent power directly from KVA rating.
func (tr *Transformador) PotenciaKVA() float64 {
	return float64(tr.KVA)
}

// PotenciaKW returns 0 — transformer is reported as apparent power only.
func (tr *Transformador) PotenciaKW() float64 {
	return 0
}

// PotenciaKVAR returns 0 — transformer is reported as apparent power only.
func (tr *Transformador) PotenciaKVAR() float64 {
	return 0
}
