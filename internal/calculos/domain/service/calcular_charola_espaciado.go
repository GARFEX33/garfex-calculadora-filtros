// internal/calculos/domain/service/calcular_charola_espaciado.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

var ErrCharolaNoEncontrada = errors.New("no se encontró charola suficiente")

// CalcularCharolaEspaciado calcula el ancho requerido de charola para cables espaciados.
// Recibe value objects del dominio para representar conductores y cables de control.

func CalcularCharolaEspaciado(
	hilosPorFase int,
	sistema entity.SistemaElectrico,
	conductorFase valueobject.ConductorCharola,
	conductorTierra valueobject.ConductorCharola,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
	cablesControl []valueobject.CableControl,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", ErrHilosPorFaseInvalido)
	}

	// Determinar numero de fases segun el tipo de sistema
	var numFases int
	var tieneNeutro bool
	switch sistema {
	case entity.SistemaElectricoMonofasico:
		numFases = 1
		tieneNeutro = true
	case entity.SistemaElectricoBifasico:
		numFases = 2
		tieneNeutro = true
	case entity.SistemaElectricoDelta:
		numFases = 3
		tieneNeutro = false
	case entity.SistemaElectricoEstrella:
		numFases = 3
		tieneNeutro = true
	default:
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", fmt.Errorf("sistema eléctrico no válido: %v", sistema))
	}

	// Calcular hilos de fase y neutro multiplicando por hilosPorFase (conductores en paralelo)
	hilosFaseTotal := numFases * hilosPorFase

	hilosNeutro := 0
	if tieneNeutro {
		hilosNeutro = hilosPorFase // El neutro se multiplica igual que las fases
	}

	totalHilos := hilosFaseTotal + hilosNeutro

	// Calcular ancho para cables en charola con espaciado
	// Formula: EF + ancho_fuerza + EC + ancho_control + tierra
	// EF (espacio fuerza) = total_hilos * diametro_fase
	// EC (espacio control) = 2 * diametro_control (uno a cada lado)
	// ancho_fuerza = total_hilos * diametro_fase
	// ancho_control = diametro_control

	// Espacio fuerza = cantidad de fases * diametro fase
	espacioFuerza := float64(totalHilos) * conductorFase.DiametroMM()

	// Espacio control = 2 * diametro control (uno a cada lado)
	var espacioControl float64
	var anchoControl float64
	for _, cable := range cablesControl {
		if cable.Cantidad() > 0 && cable.DiametroMM() > 0 {
			espacioControl += 2.0 * cable.DiametroMM() // espacio a cada lado
			anchoControl += cable.DiametroMM()         // diametro del cable
		}
	}

	// Ancho fuerza = total_hilos * diametro fase
	anchoFuerza := float64(totalHilos) * conductorFase.DiametroMM()

	// Ancho total = EF + ancho_fuerza + EC + ancho_control + tierra
	anchoRequerido := espacioFuerza + anchoFuerza + espacioControl + anchoControl + conductorTierra.DiametroMM()

	// El valor de la tabla es el ancho directo en mm
	for _, entrada := range tablaCharola {
		anchoCharolaMM := entrada.AreaInteriorMM2
		if anchoCharolaMM >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:           entity.TipoCanalizacionCharolaCableEspaciado,
				Tamano:         entrada.Tamano,
				AnchoRequerido: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", ErrCharolaNoEncontrada)
}
