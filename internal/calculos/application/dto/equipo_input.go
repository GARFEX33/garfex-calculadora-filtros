// internal/calculos/application/dto/equipo_input.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ModoCalculo define cómo se proporcionan los datos del equipo.
type ModoCalculo string

const (
	ModoListado        ModoCalculo = "LISTADO"
	ModoManualAmperaje ModoCalculo = "MANUAL_AMPERAJE"
	ModoManualPotencia ModoCalculo = "MANUAL_POTENCIA"
)

// SistemaElectrico representa el tipo de sistema eléctrico (DTO - sin métodos de negocio).
type SistemaElectrico string

const (
	SistemaElectricoDelta      SistemaElectrico = "DELTA"
	SistemaElectricoEstrella   SistemaElectrico = "ESTRELLA"
	SistemaElectricoBifasico   SistemaElectrico = "BIFASICO"
	SistemaElectricoMonofasico SistemaElectrico = "MONOFASICO"
)

// ToEntity convierte el DTO a la entidad del domain.
func (s SistemaElectrico) ToEntity() entity.SistemaElectrico {
	return entity.SistemaElectrico(s)
}

// EquipoInput contiene todos los datos necesarios para calcular una memoria.
// Es el DTO de entrada para el use case CalcularMemoria.
type EquipoInput struct {
	// Modo indica cómo se proporcionan los datos del equipo
	Modo ModoCalculo

	// Clave del equipo (requerido si Modo = LISTADO)
	Clave string

	// Datos del equipo (requeridos si Modo = MANUAL_*)
	TipoEquipo      string  // Solo para MANUAL_AMPERAJE / MANUAL_POTENCIA
	AmperajeNominal float64 // Solo para MANUAL_AMPERAJE
	PotenciaNominal float64 // Solo para MANUAL_POTENCIA (KVAR o KVA)
	Tension         float64 // Voltaje (ej: 220, 480)
	FactorPotencia  float64 // Solo para CARGA en MANUAL_POTENCIA
	ITM             int

	// Parámetros de instalación
	TipoCanalizacion      string  // "TUBERIA_PVC", "CHAROLA_CABLE_ESPACIADO", etc.
	TemperaturaOverride   *int    // nil = usar lógica por defecto, o valor override (60, 75, 90)
	HilosPorFase          int     // default: 1
	NumTuberias           int     // default: 1, para distribución de conductores
	Material              string  // "Cu" o "Al"; si vacío, default Cu
	LongitudCircuito      float64 // metros, para caída de tensión
	PorcentajeCaidaMaximo float64 // default: 3.0%

	// Sistema eléctrico (DTO con tipos primitivos, no entity)
	SistemaElectrico SistemaElectrico `json:"sistema_electrico" binding:"required"`
	Estado           string           `json:"estado" binding:"required"`
}

// Validate verifica que el input tenga los campos requeridos según el modo.
func (e EquipoInput) Validate() error {
	switch e.Modo {
	case ModoListado:
		if e.Clave == "" {
			return ErrEquipoInputInvalido
		}
	case ModoManualAmperaje:
		if e.AmperajeNominal <= 0 {
			return ErrEquipoInputInvalido
		}
	case ModoManualPotencia:
		if e.PotenciaNominal <= 0 {
			return ErrEquipoInputInvalido
		}
	default:
		return ErrModoInvalido
	}

	if e.Tension <= 0 {
		return ErrEquipoInputInvalido
	}

	if e.ITM <= 0 {
		return ErrEquipoInputInvalido
	}

	if e.Estado == "" {
		return ErrEquipoInputInvalido
	}

	// Validar sistema eléctrico
	entitySistema := e.SistemaElectrico.ToEntity()
	if err := entity.ValidarSistemaElectrico(entitySistema); err != nil {
		return err
	}

	return nil
}

// ToEntityTipoEquipo convierte el DTO string a entity.TipoEquipo.
func (e EquipoInput) ToEntityTipoEquipo() entity.TipoEquipo {
	return entity.TipoEquipo(e.TipoEquipo)
}

// ToEntityTipoCanalizacion convierte el DTO string a entity.TipoCanalizacion.
func (e EquipoInput) ToEntityTipoCanalizacion() entity.TipoCanalizacion {
	return entity.TipoCanalizacion(e.TipoCanalizacion)
}

// ToDomainTension convierte el primitivo float64 a valueobject.Tension.
// Requiere que el valor sea uno de los valores NOM válidos (127, 220, 240, 277, 440, 480, 600).
func (e EquipoInput) ToDomainTension() (valueobject.Tension, error) {
	// Convertir float64 a int para validación NOM
	tensionInt := int(e.Tension)
	return valueobject.NewTension(tensionInt)
}

// ToDomainTemperaturaOverride convierte el primitivo *int a valueobject.Temperatura.
func (e EquipoInput) ToDomainTemperaturaOverride() (valueobject.Temperatura, error) {
	if e.TemperaturaOverride == nil {
		return 0, nil
	}
	temp := valueobject.Temperatura(*e.TemperaturaOverride)
	if err := valueobject.ValidarTemperatura(temp); err != nil {
		return 0, err
	}
	return temp, nil
}

// ToDomainMaterial convierte el string a valueobject.MaterialConductor.
func (e EquipoInput) ToDomainMaterial() (valueobject.MaterialConductor, error) {
	if e.Material == "" {
		// Default a cobre
		return valueobject.ParseMaterialConductor("Cu")
	}
	return valueobject.ParseMaterialConductor(e.Material)
}

// roundToNearest finds the nearest value in the valid list.
func roundToNearest(val float64, valid []int) int {
	nearest := valid[0]
	minDiff := abs(float64(nearest) - val)
	for _, v := range valid[1:] {
		diff := abs(float64(v) - val)
		if diff < minDiff {
			minDiff = diff
			nearest = v
		}
	}
	return nearest
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
