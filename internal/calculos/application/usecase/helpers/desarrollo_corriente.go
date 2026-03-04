package helpers

import (
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
)

// generarDesarrolloCorriente genera el desarrollo paso a paso del cálculo de corriente nominal.
// Esta función replica la lógica del frontend en SeccionCorriente.svelte getInfoCalculo().
// Parámetros:
//   - tipoEquipo: tipo de equipo (FILTRO_ACTIVO, TRANSFORMADOR, FILTRO_RECHAZO, CARGA)
//   - corrienteNominal: corriente nominal calculada
//   - tension: tensión de operación en volts
//   - factorPotencia: factor de potencia del equipo
//   - esTrifasico: indica si el sistema es trifásico (ESTRELLA o DELTA)
//   - amperajeEquipo: amperaje del equipo (para FILTRO_ACTIVO)
func GenerarDesarrolloCorriente(
	tipoEquipo string,
	corrienteNominal float64,
	tension int,
	factorPotencia float64,
	esTrifasico bool,
	amperajeEquipo float64,
) *dto.DatosDesarrolloCorriente {

	// Valor por defecto de factor de potencia
	if factorPotencia == 0 {
		factorPotencia = 1.0
	}

	// Raíz de 3 para cálculos
	sqrt3 := math.Sqrt(3)

	switch tipoEquipo {
	case "FILTRO_ACTIVO":
		return generarFiltroActivo(corrienteNominal, amperajeEquipo)

	case "TRANSFORMADOR":
		return generarTransformador(corrienteNominal, tension, sqrt3)

	case "FILTRO_RECHAZO":
		return generarFiltroRechazo(corrienteNominal, tension, sqrt3)

	case "CARGA":
		if esTrifasico {
			return generarCargaTrifasico(corrienteNominal, tension, factorPotencia, sqrt3)
		}
		return generarCargaMonofasico(corrienteNominal, tension, factorPotencia)

	default:
		// Default a monofásico para equipos desconocidos
		return generarCargaMonofasico(corrienteNominal, tension, factorPotencia)
	}
}

// generarFiltroActivo genera el desarrollo para FILTRO_ACTIVO (amperaje directo).
func generarFiltroActivo(corrienteNominal, amperajeEquipo float64) *dto.DatosDesarrolloCorriente {
	// Usar amperaje del equipo si está disponible, sino usar corriente nominal
	amperaje := amperajeEquipo
	if amperaje == 0 {
		amperaje = corrienteNominal
	}

	return &dto.DatosDesarrolloCorriente{
		TipoCalculo:  "Amperaje directo",
		FormulaUsada: "I = Iₙominal",
		PasosDesarrollo: []dto.PasoDesarrollo{
			{
				Numero:      1,
				Descripcion: fmt.Sprintf("I = %.2f A (dato del equipo)", amperaje),
				Resultado:   fmt.Sprintf("I = %.2f A", corrienteNominal),
			},
		},
		ValoresReferencia: map[string]string{
			"Amperaje": fmt.Sprintf("%.2f A", amperaje),
			"Tipo":     "Filtro Activo (FP = 1.0)",
		},
	}
}

// generarTransformador genera el desarrollo para TRANSFORMADOR (desde KVA).
func generarTransformador(corrienteNominal float64, tension int, sqrt3 float64) *dto.DatosDesarrolloCorriente {
	kva := (corrienteNominal * float64(tension) * sqrt3) / 1000
	kv := float64(tension) / 1000
	divisor := kv * sqrt3

	return &dto.DatosDesarrolloCorriente{
		TipoCalculo:  "Desde KVA (Transformador)",
		FormulaUsada: "I = KVA / (kV × √3)",
		PasosDesarrollo: []dto.PasoDesarrollo{
			{
				Numero:      1,
				Descripcion: fmt.Sprintf("I = %.2f kVA / (%.3f kV × 1.732)", kva, kv),
				Resultado:   "",
			},
			{
				Numero:      2,
				Descripcion: fmt.Sprintf("I = %.2f / %.3f", kva, divisor),
				Resultado:   "",
			},
			{
				Numero:      3,
				Descripcion: "",
				Resultado:   fmt.Sprintf("I = %.2f A", corrienteNominal),
			},
		},
		ValoresReferencia: map[string]string{
			"KVA":     fmt.Sprintf("%.2f kVA", kva),
			"Voltaje": fmt.Sprintf("%d V (%.3f kV)", tension, kv),
			"Fórmula": "I = KVA / (kV × √3)",
		},
	}
}

// generarFiltroRechazo genera el desarrollo para FILTRO_RECHAZO (desde KVAR).
func generarFiltroRechazo(corrienteNominal float64, tension int, sqrt3 float64) *dto.DatosDesarrolloCorriente {
	kvar := (corrienteNominal * float64(tension) * sqrt3) / 1000
	kv := float64(tension) / 1000
	divisor := kv * sqrt3

	return &dto.DatosDesarrolloCorriente{
		TipoCalculo:  "Desde KVAR (Filtro de Rechazo)",
		FormulaUsada: "I = KVAR / (kV × √3)",
		PasosDesarrollo: []dto.PasoDesarrollo{
			{
				Numero:      1,
				Descripcion: fmt.Sprintf("I = %.2f kVAR / (%.3f kV × 1.732)", kvar, kv),
				Resultado:   "",
			},
			{
				Numero:      2,
				Descripcion: fmt.Sprintf("I = %.2f / %.3f", kvar, divisor),
				Resultado:   "",
			},
			{
				Numero:      3,
				Descripcion: "",
				Resultado:   fmt.Sprintf("I = %.2f A", corrienteNominal),
			},
		},
		ValoresReferencia: map[string]string{
			"KVAR":    fmt.Sprintf("%.2f kVAR", kvar),
			"Voltaje": fmt.Sprintf("%d V (%.3f kV)", tension, kv),
			"Fórmula": "I = KVAR / (kV × √3)",
		},
	}
}

// generarCargaTrifasico genera el desarrollo para CARGA trifásica.
func generarCargaTrifasico(corrienteNominal float64, tension int, factorPotencia float64, sqrt3 float64) *dto.DatosDesarrolloCorriente {
	potenciaKW := (corrienteNominal * float64(tension) * sqrt3 * factorPotencia) / 1000
	potenciaW := potenciaKW * 1000
	divisor := float64(tension) * sqrt3 * factorPotencia

	return &dto.DatosDesarrolloCorriente{
		TipoCalculo:  "Desde Potencia (Sistema Trifásico)",
		FormulaUsada: "I = P / (V × √3 × cosθ)",
		PasosDesarrollo: []dto.PasoDesarrollo{
			{
				Numero:      1,
				Descripcion: fmt.Sprintf("P = %.2f kW = %.0f W", potenciaKW, potenciaW),
				Resultado:   "",
			},
			{
				Numero:      2,
				Descripcion: fmt.Sprintf("I = %.0f / (%d × 1.732 × %.2f)", potenciaW, tension, factorPotencia),
				Resultado:   "",
			},
			{
				Numero:      3,
				Descripcion: fmt.Sprintf("I = %.0f / %.2f", potenciaW, divisor),
				Resultado:   "",
			},
			{
				Numero:      4,
				Descripcion: "",
				Resultado:   fmt.Sprintf("I = %.2f A", corrienteNominal),
			},
		},
		ValoresReferencia: map[string]string{
			"Potencia":           fmt.Sprintf("%.2f kW", potenciaKW),
			"Voltaje":            fmt.Sprintf("%d V", tension),
			"Factor de Potencia": fmt.Sprintf("%.2f", factorPotencia),
			"Sistema":            "Trifásico",
			"Fórmula":            "I = P / (V × √3 × cosθ)",
		},
	}
}

// generarCargaMonofasico genera el desarrollo para CARGA monofásica o bifásica.
func generarCargaMonofasico(corrienteNominal float64, tension int, factorPotencia float64) *dto.DatosDesarrolloCorriente {
	potenciaKW := (corrienteNominal * float64(tension) * factorPotencia) / 1000
	potenciaW := potenciaKW * 1000
	divisor := float64(tension) * factorPotencia

	return &dto.DatosDesarrolloCorriente{
		TipoCalculo:  "Desde Potencia (Sistema Monofásico)",
		FormulaUsada: "I = P / (V × cosθ)",
		PasosDesarrollo: []dto.PasoDesarrollo{
			{
				Numero:      1,
				Descripcion: fmt.Sprintf("P = %.2f kW = %.0f W", potenciaKW, potenciaW),
				Resultado:   "",
			},
			{
				Numero:      2,
				Descripcion: fmt.Sprintf("I = %.0f / (%d × %.2f)", potenciaW, tension, factorPotencia),
				Resultado:   "",
			},
			{
				Numero:      3,
				Descripcion: fmt.Sprintf("I = %.0f / %.2f", potenciaW, divisor),
				Resultado:   "",
			},
			{
				Numero:      4,
				Descripcion: "",
				Resultado:   fmt.Sprintf("I = %.2f A", corrienteNominal),
			},
		},
		ValoresReferencia: map[string]string{
			"Potencia":           fmt.Sprintf("%.2f kW", potenciaKW),
			"Voltaje":            fmt.Sprintf("%d V", tension),
			"Factor de Potencia": fmt.Sprintf("%.2f", factorPotencia),
			"Sistema":            "Monofásico",
			"Fórmula":            "I = P / (V × cosθ)",
		},
	}
}
