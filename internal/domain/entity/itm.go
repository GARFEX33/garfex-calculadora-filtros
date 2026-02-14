// internal/domain/entity/itm.go
package entity

import "fmt"

// ITM represents a thermal-magnetic circuit breaker (Interruptor Termomagnético).
// Every electrical installation requires a protection device.
type ITM struct {
	Amperaje int // rated current [A]
	Polos    int // number of poles (3 for three-phase installations)
	Bornes   int // number of conductor terminals
	Voltaje  int // rated voltage [V]
}

// NewITM constructs a validated ITM.
// For three-phase filter installations, Polos=3 and Voltaje=equipment voltage.
func NewITM(amperaje, polos, bornes, voltaje int) (ITM, error) {
	if amperaje <= 0 {
		return ITM{}, fmt.Errorf("ITM: amperaje debe ser mayor que cero: %d", amperaje)
	}
	if polos <= 0 {
		return ITM{}, fmt.Errorf("ITM: polos debe ser mayor que cero: %d", polos)
	}
	if bornes <= 0 {
		return ITM{}, fmt.Errorf("ITM: bornes debe ser mayor que cero: %d", bornes)
	}
	if voltaje <= 0 {
		return ITM{}, fmt.Errorf("ITM: voltaje debe ser mayor que cero: %d", voltaje)
	}
	return ITM{
		Amperaje: amperaje,
		Polos:    polos,
		Bornes:   bornes,
		Voltaje:  voltaje,
	}, nil
}
