package helpers

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

func NombreTablaAmpacidad(
	canalizacion string,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
) string {
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

	mat := "Cu"
	if material == valueobject.MaterialAluminio {
		mat = "Al"
	}

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
