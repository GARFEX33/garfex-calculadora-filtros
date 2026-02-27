// internal/calculos/application/dto/equipo_input.go
package dto

import (
	"fmt"

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

// TipoFiltro representa el tipo de filtro según viene de la BD/equipos.
// Valores: "A", "KVA", "KVAR"
type TipoFiltro string

const (
	TipoFiltroA    TipoFiltro = "A"    // Filtro activo en Amperes
	TipoFiltroKVA  TipoFiltro = "KVA"  // Filtro calificado en KVA
	TipoFiltroKVAR TipoFiltro = "KVAR" // Filtro de rechazo en KVAR
)

// ToTipoEquipo mapea TipoFiltro (de BD/equipos) a entity.TipoEquipo.
func (t TipoFiltro) ToTipoEquipo() (entity.TipoEquipo, error) {
	switch t {
	case TipoFiltroA:
		return entity.TipoEquipoFiltroActivo, nil
	case TipoFiltroKVA:
		return entity.TipoEquipoTransformador, nil
	case TipoFiltroKVAR:
		return entity.TipoEquipoFiltroRechazo, nil
	default:
		return "", fmt.Errorf("tipo de filtro inválido: '%s' (válidos: A, KVA, KVAR)", t)
	}
}

// DatosEquipo contiene los datos del equipo obtenidos del endpoint de equipos.
// En modo LISTADO, el frontend envía estos datos tal cual los recibió de GET /equipos.
type DatosEquipo struct {
	Clave    string     `json:"clave"`    // Clave comercial (ej: ACTISINE48D400MMB)
	Tipo     TipoFiltro `json:"tipo"`     // "A", "KVA", "KVAR"
	Voltaje  int        `json:"voltaje"`  // Voltaje nominal (ej: 480)
	Amperaje int        `json:"amperaje"` // Qn/In - amperaje nominal o KVA/KVAR según tipo
	ITM      int        `json:"itm"`      // Interruptor termomagnético
	Bornes   *int       `json:"bornes"`   // Número de bornes (opcional)
}

// ToTipoEquipo mapea el tipo de filtro a TipoEquipo.
func (d DatosEquipo) ToTipoEquipo() (entity.TipoEquipo, error) {
	return d.Tipo.ToTipoEquipo()
}

// EquipoInput contiene todos los datos necesarios para calcular una memoria.
type EquipoInput struct {
	// Modo indica cómo se proporcionan los datos del equipo
	Modo ModoCalculo

	// ═══════════════════════════════════════════════════════════════════════
	// DATOS DEL EQUIPO
	// ═══════════════════════════════════════════════════════════════════════
	// LISTADO: El frontend envía DatosEquipo tal cual de GET /equipos
	// MANUAL_AMPERAJE: Solo se usa TipoEquipo y AmperajeNominal
	// MANUAL_POTENCIA: Se usa TipoEquipo, PotenciaNominal, PotenciaUnidad, FactorPotencia
	Equipo          DatosEquipo // Datos del equipo (LISTADO)
	TipoEquipo      string      // MANUAL_*: FILTRO_ACTIVO, TRANSFORMADOR, FILTRO_RECHAZO, CARGA
	AmperajeNominal float64     // MANUAL_AMPERAJE: amperaje directo
	PotenciaNominal float64     // MANUAL_POTENCIA: valor de potencia
	PotenciaUnidad  string      // MANUAL_POTENCIA: W, KW, KVA, KVAR
	FactorPotencia  float64     // MANUAL_POTENCIA: solo para CARGA

	// ═══════════════════════════════════════════════════════════════════════
	// DATOS DE INSTALACIÓN (comunes a todos los modos)
	// ═══════════════════════════════════════════════════════════════════════
	Tension               float64  // Voltaje de referencia para cálculos
	TensionUnidad         string   `json:"tension_unidad"` // "V" o "kV" (default: "V")
	TipoCanalizacion      string   // "TUBERIA_PVC", "CHAROLA_CABLE_ESPACIADO", etc.
	TemperaturaOverride   *int     // nil = usar lógica por defecto
	HilosPorFase          int      // default: 1
	NumTuberias           int      // default: 1
	Material              string   // "Cu" o "Al"; default: Cu
	LongitudCircuito      float64  // metros
	PorcentajeCaidaMaximo float64  // default: 3.0%
	DiametroControlMM     *float64 // opcional, para cables de control en charola

	// Sistema eléctrico
	SistemaElectrico SistemaElectrico
	Estado           string
	TipoVoltaje      string // "FASE_NEUTRO" o "FASE_FASE"
}

// Validate verifica que el input tenga los campos requeridos según el modo.
func (e EquipoInput) Validate() error {
	switch e.Modo {
	case ModoListado:
		// Validar datos del equipo
		if e.Equipo.Tipo == "" {
			return fmt.Errorf("%w: tipo de equipo requerido", ErrEquipoInputInvalido)
		}
		if e.Equipo.Voltaje <= 0 {
			return fmt.Errorf("%w: voltaje debe ser mayor que cero", ErrEquipoInputInvalido)
		}
		if e.Equipo.Amperaje <= 0 {
			return fmt.Errorf("%w: amperaje debe ser mayor que cero", ErrEquipoInputInvalido)
		}
		if e.Equipo.ITM <= 0 {
			return fmt.Errorf("%w: ITM debe ser mayor que cero", ErrEquipoInputInvalido)
		}
		// Validar que el tipo sea mapeable
		if _, err := e.Equipo.ToTipoEquipo(); err != nil {
			return err
		}

	case ModoManualAmperaje:
		if e.AmperajeNominal <= 0 {
			return fmt.Errorf("%w: amperaje_nominal requerido en modo MANUAL_AMPERAJE", ErrEquipoInputInvalido)
		}
		if e.TipoEquipo == "" {
			return fmt.Errorf("%w: tipo_equipo requerido en modo MANUAL_AMPERAJE", ErrEquipoInputInvalido)
		}

	case ModoManualPotencia:
		if e.PotenciaNominal <= 0 {
			return fmt.Errorf("%w: potencia_nominal requerido en modo MANUAL_POTENCIA", ErrEquipoInputInvalido)
		}
		if e.TipoEquipo == "" {
			return fmt.Errorf("%w: tipo_equipo requerido en modo MANUAL_POTENCIA", ErrEquipoInputInvalido)
		}

	default:
		return ErrModoInvalido
	}

	// Validar tensión (común a todos los modos)
	if e.Tension <= 0 {
		return fmt.Errorf("%w: tension requerida", ErrEquipoInputInvalido)
	}

	// Validar estado
	if e.Estado == "" {
		return fmt.Errorf("%w: estado requerido", ErrEquipoInputInvalido)
	}

	// Validar sistema eléctrico
	entitySistema := e.SistemaElectrico.ToEntity()
	if err := entity.ValidarSistemaElectrico(entitySistema); err != nil {
		return err
	}

	return nil
}

// ValidateForMemoria validates ALL fields required for the full memory calculation pipeline.
func (e EquipoInput) ValidateForMemoria() error {
	if err := e.Validate(); err != nil {
		return err
	}

	// Validate TipoCanalizacion
	if e.TipoCanalizacion == "" {
		return fmt.Errorf("%w: tipo_canalizacion requerido", ErrEquipoInputInvalido)
	}
	tipoCanalizacion := e.ToEntityTipoCanalizacion()
	if err := entity.ValidarTipoCanalizacion(tipoCanalizacion); err != nil {
		return err
	}

	// Validate LongitudCircuito
	if e.LongitudCircuito <= 0 {
		return fmt.Errorf("%w: longitud_circuito requerida", ErrEquipoInputInvalido)
	}

	// Validate TipoVoltaje
	if e.TipoVoltaje == "" {
		return fmt.Errorf("%w: tipo_voltaje requerido", ErrEquipoInputInvalido)
	}
	if _, err := e.ToDomainTipoVoltaje(); err != nil {
		return err
	}

	// Validate PorcentajeCaidaMaximo
	if e.PorcentajeCaidaMaximo < 0 {
		return fmt.Errorf("%w: porcentaje_caida_maximo no puede ser negativo", ErrEquipoInputInvalido)
	}

	// Validate HilosPorFase
	if e.HilosPorFase < 0 {
		return fmt.Errorf("%w: hilos_por_fase no puede ser negativo", ErrEquipoInputInvalido)
	}

	// Validate NumTuberias
	if e.NumTuberias < 0 {
		return fmt.Errorf("%w: num_tuberias no puede ser negativo", ErrEquipoInputInvalido)
	}

	// Validate FactorPotencia for MANUAL_POTENCIA mode
	if e.Modo == ModoManualPotencia && (e.FactorPotencia <= 0 || e.FactorPotencia > 1) {
		return fmt.Errorf("%w: factor_potencia debe estar entre 0 y 1", ErrEquipoInputInvalido)
	}

	return nil
}

// ApplyDefaults sets default values for optional fields.
func (e *EquipoInput) ApplyDefaults() {
	if e.HilosPorFase <= 0 {
		e.HilosPorFase = 1
	}
	if e.NumTuberias <= 0 {
		e.NumTuberias = 1
	}
	if e.PorcentajeCaidaMaximo <= 0 {
		e.PorcentajeCaidaMaximo = 3.0
	}
	if e.Material == "" {
		e.Material = "Cu"
	}
	if e.Modo == ModoManualPotencia && e.PotenciaUnidad == "" {
		e.PotenciaUnidad = "KW"
	}
	if e.TensionUnidad == "" {
		e.TensionUnidad = "V"
	}
	if e.FactorPotencia <= 0 {
		e.FactorPotencia = 1.0
	}
}

// GetTipoEquipo retorna el TipoEquipo según el modo.
func (e EquipoInput) GetTipoEquipo() (entity.TipoEquipo, error) {
	switch e.Modo {
	case ModoListado:
		return e.Equipo.ToTipoEquipo()
	case ModoManualAmperaje, ModoManualPotencia:
		return entity.ParseTipoEquipo(e.TipoEquipo)
	default:
		return "", ErrModoInvalido
	}
}

// GetAmperajeNominal retorna el amperaje nominal según el modo.
func (e EquipoInput) GetAmperajeNominal() float64 {
	switch e.Modo {
	case ModoListado:
		return float64(e.Equipo.Amperaje)
	case ModoManualAmperaje:
		return e.AmperajeNominal
	default:
		return 0
	}
}

// GetITM retorna el ITM según el modo.
func (e EquipoInput) GetITM() int {
	if e.Modo == ModoListado {
		return e.Equipo.ITM
	}
	// Para MANUAL_*, el ITM debe venir en el input de instalación
	// (se maneja en el handler)
	return 0
}

// ToEntityTipoCanalizacion convierte el DTO string a entity.TipoCanalizacion.
func (e EquipoInput) ToEntityTipoCanalizacion() entity.TipoCanalizacion {
	return entity.TipoCanalizacion(e.TipoCanalizacion)
}

// ToDomainTension converts the primitive float64 to valueobject.Tension.
func (e EquipoInput) ToDomainTension() (valueobject.Tension, error) {
	unidad := e.TensionUnidad
	if unidad == "" {
		unidad = "V"
	}
	return valueobject.NewTension(e.Tension, unidad)
}

// ToDomainPotencia convierte el primitivo float64 a valueobject.Potencia.
func (e EquipoInput) ToDomainPotencia() (valueobject.Potencia, error) {
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
		return valueobject.ParseMaterialConductor("Cu")
	}
	return valueobject.ParseMaterialConductor(e.Material)
}

// ToDomainTipoVoltaje convierte el string a entity.TipoVoltaje.
func (e EquipoInput) ToDomainTipoVoltaje() (entity.TipoVoltaje, error) {
	return entity.ParseTipoVoltaje(e.TipoVoltaje)
}
