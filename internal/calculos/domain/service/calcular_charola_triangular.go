// internal/calculos/domain/service/calcular_charola_triangular.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

var ErrCharolaTriangularNoEncontrada = errors.New("no se encontró charola triangular suficiente")

// ErrTablaCharolaVacia is returned when the charola sizing table is empty.
var ErrTablaCharolaVacia = errors.New("tabla de charola vacía")

// obtenerAnchoCharola retorna el ancho de la charola en mm directamente del valor en la tabla.
// El archivo CSV charola_dimensiones.csv tiene los valores en mm (ej: 152.4 para 6 pulgadas).
func obtenerAnchoCharola(entrada valueobject.EntradaTablaCanalizacion) float64 {
	return entrada.AreaInteriorMM2
}

func CalcularCharolaTriangular(
	hilosPorFase int,
	conductorFase valueobject.ConductorCharola,
	conductorTierra valueobject.ConductorCharola,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
	cablesControl []valueobject.CableControl,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", ErrHilosPorFaseInvalido)
	}
	if len(tablaCharola) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", ErrTablaCharolaVacia)
	}

	// factorTriangular: factor de espaciado NOM-001-SEDE para disposición triangular de cables en charola.
	const factorTriangular = 2.15
	// Calcular ancho requerido para charola triangular
	// AP = 2 * Ø_fase * hilosPorFase
	anchoPotencia := 2.0 * conductorFase.DiametroMM() * float64(hilosPorFase)
	espacioFuerza := float64(hilosPorFase-1) * factorTriangular * conductorFase.DiametroMM()

	// Espacio y ancho de control (a ambos lados)
	var espacioControl float64
	var anchoControl float64
	for _, cable := range cablesControl {
		if cable.Cantidad() > 0 && cable.DiametroMM() > 0 {
			espacioControl += 2.15 * cable.DiametroMM() // espacio a cada lado
			anchoControl += cable.DiametroMM()          // diametro del cable
		}
	}

	// Ancho total = potencia + espacio fuerza + control + tierra
	anchoRequerido := anchoPotencia + espacioFuerza + espacioControl + anchoControl + conductorTierra.DiametroMM()

	// Seleccionar charola por ancho
	for _, entrada := range tablaCharola {
		anchoCharolaMM := obtenerAnchoCharola(entrada)
		if anchoCharolaMM >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:           entity.TipoCanalizacionCharolaCableTriangular,
				Tamano:         entrada.Tamano,
				AnchoRequerido: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", ErrCharolaTriangularNoEncontrada)
}
