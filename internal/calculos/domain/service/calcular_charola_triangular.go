// internal/calculos/domain/service/calcular_charola_triangular.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

var ErrCharolaTriangularNoEncontrada = errors.New("no se encontró charola triangular suficiente")

func CalcularCharolaTriangular(
	hilosPorFase int,
	conductorFase valueobject.ConductorCharola,
	conductorTierra valueobject.ConductorCharola,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
	cablesControl []valueobject.CableControl,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", errors.New("hilos por fase debe ser >= 1"))
	}
	if len(tablaCharola) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", errors.New("tabla vacía"))
	}

	// factorTriangular: factor de espaciado NOM-001-SEDE para disposición triangular de cables en charola.
	const factorTriangular = 2.15
	// Calcular ancho requerido para charola triangular
	anchoPotencia := 2.0 * conductorFase.DiametroMM()
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

	for _, entrada := range tablaCharola {
		if entrada.AreaInteriorMM2 >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:           string(entity.TipoCanalizacionCharolaCableTriangular),
				Tamano:         entrada.Tamano,
				AnchoRequerido: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaTriangular: %w", ErrCharolaTriangularNoEncontrada)
}
