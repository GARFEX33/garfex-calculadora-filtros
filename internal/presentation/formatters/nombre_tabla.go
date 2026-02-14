// internal/presentation/formatters/nombre_tabla.go
package formatters

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// NombreTablaAmpacidad genera el nombre descriptivo de la tabla NOM usada.
func NombreTablaAmpacidad(
	canalizacion string,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
) string {
	// Mapeo de canalización a tabla NOM
	var tabla string
	switch canalizacion {
	case "TUBERIA_PVC", "TUBERIA_ALUMINIO", "TUBERIA_ACERO_PG", "TUBERIA_ACERO_PD":
		tabla = "NOM-310-15-B-16"
	case "CHAROLA_CABLE_ESPACIADO":
		tabla = "NOM-310-15-B-17"
	case "CHAROLA_CABLE_TRIANGULAR":
		tabla = "NOM-310-15-B-20"
	default:
		tabla = "NOM-310-15-B-16"
	}

	// Material
	mat := "Cu"
	if material == valueobject.MaterialAluminio {
		mat = "Al"
	}

	// Temperatura
	temp := "75°C"
	switch temperatura {
	case valueobject.Temp60:
		temp = "60°C"
	case valueobject.Temp75:
		temp = "75°C"
	case valueobject.Temp90:
		temp = "90°C"
	}

	return fmt.Sprintf("%s (%s, %s)", tabla, mat, temp)
}
