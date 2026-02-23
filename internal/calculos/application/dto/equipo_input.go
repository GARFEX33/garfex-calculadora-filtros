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
	PotenciaNominal float64 // Solo para MANUAL_POTENCIA (valor)
	PotenciaUnidad  string  // Solo para MANUAL_POTENCIA: "W", "KW", "KVA", "KVAR" (default: "KW")
	Tension         float64 `json:"tension"`        // Voltaje (ej: 220, 480, o 0.48 para kV)
	TensionUnidad   string  `json:"tension_unidad"` // Unidad de tensión: "V" o "kV" (default: "V")
	FactorPotencia  float64 // Solo para CARGA en MANUAL_POTENCIA
	ITM             int

	// Parámetros de instalación
	TipoCanalizacion      string   // "TUBERIA_PVC", "CHAROLA_CABLE_ESPACIADO", etc.
	TemperaturaOverride   *int     // nil = usar lógica por defecto, o valor override (60, 75, 90)
	HilosPorFase          int      // default: 1
	NumTuberias           int      // default: 1, para distribución de conductores
	Material              string   // "Cu" o "Al"; si vacío, default Cu
	LongitudCircuito      float64  // metros, para caída de tensión
	PorcentajeCaidaMaximo float64  // default: 3.0%
	DiametroControlMM     *float64 // opcional, para cables de control en charola

	// Sistema eléctrico (DTO con tipos primitivos, no entity)
	SistemaElectrico SistemaElectrico `json:"sistema_electrico" binding:"required"`
	Estado           string           `json:"estado" binding:"required"`

	// Tipo de voltaje (FASE_NEUTRO o FASE_FASE) — requerido para caída de tensión
	TipoVoltaje string `json:"tipo_voltaje" binding:"required"`
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

// ValidateForMemoria validates ALL fields required for the full memory calculation pipeline.
// This includes validation for: mode fields + tipo canalización + longitud circuito + tipo voltaje.
func (e EquipoInput) ValidateForMemoria() error {
	// First, validate basic fields
	if err := e.Validate(); err != nil {
		return err
	}

	// Validate TipoCanalizacion
	if e.TipoCanalizacion == "" {
		return ErrEquipoInputInvalido
	}
	tipoCanalizacion := e.ToEntityTipoCanalizacion()
	if err := entity.ValidarTipoCanalizacion(tipoCanalizacion); err != nil {
		return err
	}

	// Validate LongitudCircuito (required for voltage drop calculation)
	if e.LongitudCircuito <= 0 {
		return ErrEquipoInputInvalido
	}

	// Validate TipoVoltaje (required for voltage drop calculation)
	if e.TipoVoltaje == "" {
		return ErrEquipoInputInvalido
	}
	_, err := e.ToDomainTipoVoltaje()
	if err != nil {
		return err
	}

	// Validate PorcentajeCaidaMaximo (optional, but if provided must be > 0)
	if e.PorcentajeCaidaMaximo < 0 {
		return ErrEquipoInputInvalido
	}

	// Validate HilosPorFase (optional, but if provided must be >= 1)
	if e.HilosPorFase < 0 {
		return ErrEquipoInputInvalido
	}

	// Validate NumTuberias (optional, but if provided must be >= 1)
	if e.NumTuberias < 0 {
		return ErrEquipoInputInvalido
	}

	// Validate FactorPotencia for MANUAL_POTENCIA mode
	if e.Modo == ModoManualPotencia && (e.FactorPotencia <= 0 || e.FactorPotencia > 1) {
		return ErrEquipoInputInvalido
	}

	return nil
}

// ApplyDefaults sets default values for optional fields when they are zero/empty.
// Call this before validation to ensure defaults are applied.
func (e *EquipoInput) ApplyDefaults() {
	// Default: 1 hilo por fase
	if e.HilosPorFase <= 0 {
		e.HilosPorFase = 1
	}

	// Default: 1 tubería
	if e.NumTuberias <= 0 {
		e.NumTuberias = 1
	}

	// Default: 3% caída máxima
	if e.PorcentajeCaidaMaximo <= 0 {
		e.PorcentajeCaidaMaximo = 3.0
	}

	// Default: Cobre
	if e.Material == "" {
		e.Material = "Cu"
	}

	// Default: KW para potencia
	if e.Modo == ModoManualPotencia && e.PotenciaUnidad == "" {
		e.PotenciaUnidad = "KW"
	}

	// Default: V para tensión
	if e.TensionUnidad == "" {
		e.TensionUnidad = "V"
	}
}

// ToEntityTipoEquipo convierte el DTO string a entity.TipoEquipo.
func (e EquipoInput) ToEntityTipoEquipo() entity.TipoEquipo {
	return entity.TipoEquipo(e.TipoEquipo)
}

// ToEntityTipoCanalizacion convierte el DTO string a entity.TipoCanalizacion.
func (e EquipoInput) ToEntityTipoCanalizacion() entity.TipoCanalizacion {
	return entity.TipoCanalizacion(e.TipoCanalizacion)
}

// ToDomainTension converts the primitive float64 to valueobject.Tension.
// Requires TensionUnidad to be a valid unit: V, kV (default: V).
// The value must be one of the valid NOM values (127, 220, 240, 277, 440, 480, 600).
func (e EquipoInput) ToDomainTension() (valueobject.Tension, error) {
	// Apply default if not set (should be done by ApplyDefaults, but be safe)
	unidad := e.TensionUnidad
	if unidad == "" {
		unidad = "V"
	}
	return valueobject.NewTension(e.Tension, unidad)
}

// ToDomainPotencia convierte el primitivo float64 a valueobject.Potencia.
// Requiere que PotenciaUnidad sea una unidad válida: W, KW, KVA, KVAR.
func (e EquipoInput) ToDomainPotencia() (valueobject.Potencia, error) {
	// Apply default if not set
	if e.PotenciaUnidad == "" {
		e.PotenciaUnidad = "KW"
	}
	return valueobject.NewPotencia(e.PotenciaNominal, e.PotenciaUnidad)
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

// ToDomainTipoVoltaje convierte el string a entity.TipoVoltaje.
func (e EquipoInput) ToDomainTipoVoltaje() (entity.TipoVoltaje, error) {
	return entity.ParseTipoVoltaje(e.TipoVoltaje)
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
