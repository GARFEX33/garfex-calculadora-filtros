// internal/calculos/domain/service/calculo_amperaje_nominal.go
package service

import (
	"errors"
	"math"
	"math/cmplx"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// Tipos de carga eléctrica
type TipoCarga int

const (
	TipoCargaMonofasica TipoCarga = iota
	TipoCargaTrifasica
)

// Sistemas eléctricos
type SistemaElectrico int

const (
	SistemaElectricoEstrella SistemaElectrico = iota
	SistemaElectricoDelta
)

// Errores del cálculo de amperaje nominal
var (
	ErrPotenciaInvalida = errors.New("potencia debe ser mayor que cero")
	ErrTensionInvalida  = errors.New("tensión no válida")
	// ErrFactorPotenciaInvalido está definido en calculo_caida_tension.go
)

// CalcularAmperajeNominalCircuito calcula el amperaje nominal de un circuito
// eléctrico según fórmulas NOM-001-SEDE.
//
// Parámetros:
//   - potenciaWatts: potencia activa en Watts
//   - tension: tensión del circuito (V)
//   - tipoCarga: TipoCargaMonofasica o TipoCargaTrifasica
//   - sistemaElectrico: SistemaElectricoEstrella o SistemaElectricoDelta
//   - factorPotencia: factor de potencia (0 a 1)
//
// Fórmulas según NOM:
//   - Monofásico: I = W / (V × FP)
//   - Trifásico: I = W / (V × √3 × FP)
func CalcularAmperajeNominalCircuito(
	potenciaWatts float64,
	tension valueobject.Tension,
	tipoCarga TipoCarga,
	sistemaElectrico SistemaElectrico,
	factorPotencia float64,
) (valueobject.Corriente, error) {
	// Validaciones
	if potenciaWatts <= 0 {
		return valueobject.Corriente{}, ErrPotenciaInvalida
	}
	if factorPotencia <= 0 || factorPotencia > 1 {
		return valueobject.Corriente{}, ErrFactorPotenciaInvalido
	}

	// Obtener valor de tensión
	tensionV := float64(tension.Valor())
	if tensionV <= 0 {
		return valueobject.Corriente{}, ErrTensionInvalida
	}

	// Calcular según tipo de carga
	var amperaje float64

	// Raíz de 3 para cálculos trifásicos
	raiz3 := math.Sqrt(3)

	switch tipoCarga {
	case TipoCargaMonofasica:
		// I = W / (V × FP)
		amperaje = potenciaWatts / (tensionV * factorPotencia)

	case TipoCargaTrifasica:
		// I = W / (V × √3 × FP)
		amperaje = potenciaWatts / (tensionV * raiz3 * factorPotencia)

	default:
		// Por defecto, asumir trifásico (caso más común en instalaciones industriales)
		amperaje = potenciaWatts / (tensionV * raiz3 * factorPotencia)
	}

	// Crear value object Corriente
	corriente, err := valueobject.NewCorriente(amperaje)
	if err != nil {
		return valueobject.Corriente{}, err
	}

	return corriente, nil
}

// CalcularAmperajeNominalComplejo calcula la corriente considerando potencia aparente,
// potencia reactiva y factor de potencia. Útil para sistemas con carga mixta.
//
// Potencia aparente (VA) = W / FP
// Potencia reactiva (VAR) = W × tan(arccos(FP))
// I = √(I_activa² + I_reactiva²)
func CalcularAmperajeNominalComplejo(
	potenciaWatts float64,
	tension valueobject.Tension,
	tipoCarga TipoCarga,
	factorPotencia float64,
) (valueobject.Corriente, error) {
	// Validaciones
	if potenciaWatts <= 0 {
		return valueobject.Corriente{}, ErrPotenciaInvalida
	}
	if factorPotencia <= 0 || factorPotencia > 1 {
		return valueobject.Corriente{}, ErrFactorPotenciaInvalido
	}

	tensionV := float64(tension.Valor())
	if tensionV <= 0 {
		return valueobject.Corriente{}, ErrTensionInvalida
	}

	// Calcular ángulo de fase
	anguloFase := math.Acos(factorPotencia)

	// Potencia aparente
	potenciaAparente := potenciaWatts / factorPotencia

	var amperaje float64

	// Raíz de 3 para cálculos trifásicos
	raiz3 := math.Sqrt(3)

	switch tipoCarga {
	case TipoCargaMonofasica:
		// Corriente aparente = VA / V
		corrienteAparente := potenciaAparente / tensionV
		// Componente activa
		corrienteActiva := corrienteAparente * factorPotencia
		// Componente reactiva = corrienteAparente × sin(ángulo)
		corrienteReactiva := corrienteAparente * math.Sin(anguloFase)
		// Corriente total = √(I² + Ir²)
		amperaje = cmplx.Abs(complex(corrienteActiva, corrienteReactiva))

	case TipoCargaTrifasica:
		// Corriente aparente = VA / (V × √3)
		corrienteAparente := potenciaAparente / (tensionV * raiz3)
		// Componente activa
		corrienteActiva := corrienteAparente * factorPotencia
		// Componente reactiva
		corrienteReactiva := corrienteAparente * math.Sin(anguloFase)
		// Corriente total
		amperaje = cmplx.Abs(complex(corrienteActiva, corrienteReactiva))

	default:
		amperaje = potenciaWatts / (tensionV * raiz3 * factorPotencia)
	}

	corriente, err := valueobject.NewCorriente(amperaje)
	if err != nil {
		return valueobject.Corriente{}, err
	}

	return corriente, nil
}
