// internal/calculos/application/dto/memoria_output.go
package dto

// PasoMemoria representa un paso individual del cálculo.
type PasoMemoria struct {
	Numero      int
	Nombre      string
	Descripcion string
	Resultado   interface{}
}

// ResultadoConductor contiene la información del conductor seleccionado.
type ResultadoConductor struct {
	Calibre         string
	Material        string
	SeccionMM2      float64
	TipoAislamiento string
	Capacidad       float64
	NumHilos        int // Número de hilos de tierra (1 para charola/tubería≤2, 2 para tubería>2)
}

// ResultadoConductores contiene los conductores seleccionados.
type ResultadoConductores struct {
	Alimentacion ResultadoConductor
	Tierra       ResultadoConductor
	TablaUsada   string
}

// ResultadoCanalizacion contiene la información de la canalización.
type ResultadoCanalizacion struct {
	Tamano           string
	AnchoComercialMM float64 `json:"ancho_comercial_mm,omitempty"`
	AreaTotalMM2     float64
	AreaRequeridaMM2 float64
	NumeroDeTubos    int
}

// DetalleCharola contiene los valores intermedios del cálculo de charola
// para mostrar el desarrollo completo en la memoria de cálculo.
type DetalleCharola struct {
	// Diámetros de los conductores (mm)
	DiametroFaseMM    float64  `json:"diametro_fase_mm"`
	DiametroTierraMM  float64  `json:"diametro_tierra_mm"`
	DiametroControlMM *float64 `json:"diametro_control_mm,omitempty"`

	// Charola espaciado
	NumHilosTotal    int     `json:"num_hilos_total,omitempty"`
	EspacioFuerzaMM  float64 `json:"espacio_fuerza_mm"`
	AnchoFuerzaMM    float64 `json:"ancho_fuerza_mm,omitempty"`
	EspacioControlMM float64 `json:"espacio_control_mm,omitempty"`
	AnchoControlMM   float64 `json:"ancho_control_mm,omitempty"`
	AnchoTierraMM    float64 `json:"ancho_tierra_mm"`

	// Charola triangular (adicional)
	AnchoPotenciaMM  float64 `json:"ancho_potencia_mm,omitempty"`
	FactorTriangular float64 `json:"factor_triangular,omitempty"`
}

// DetalleTuberia contiene los valores intermedios del cálculo de tubería
// para mostrar el desarrollo completo en la memoria de cálculo.
type DetalleTuberia struct {
	// Áreas físicas de conductores (mm²)
	// Fase: Tabla 5 NOM (área con aislamiento THW)
	// Tierra: Tabla 8 NOM (conductor desnudo)
	AreaFaseMM2   float64  `json:"area_fase_mm2"`
	AreaNeutroMM2 *float64 `json:"area_neutro_mm2,omitempty"` // nil si DELTA (sin neutro)
	AreaTierraMM2 float64  `json:"area_tierra_mm2"`

	// Distribución de conductores por tubo
	NumFasesPorTubo   int `json:"num_fases_por_tubo"`
	NumNeutrosPorTubo int `json:"num_neutros_por_tubo"` // 0 si DELTA
	NumTierras        int `json:"num_tierras"`          // 1 o 2 según NOM

	// Tubo seleccionado de la tabla NOM (Cap. 9)
	AreaOcupacionTuboMM2 float64 `json:"area_ocupacion_tubo_mm2"` // área_ocupacion_mm2 del CSV (40% interior ya aplicado)
	DesignacionMetrica   string  `json:"designacion_metrica"`     // ej: "63" → mostrar como "63 mm"
	FillFactor           float64 `json:"fill_factor"`             // 0.40 para >2 conductores
}

// ResultadoCaidaTension contiene el resultado del cálculo de caída.
type ResultadoCaidaTension struct {
	Porcentaje       float64
	CaidaVolts       float64
	Cumple           bool
	LimitePorcentaje float64
	Impedancia       float64
}

// EntradaDimensionarCanalizacion es el DTO de entrada para DimensionarCanalizacionUseCase.
type EntradaDimensionarCanalizacion struct {
	ConductorAlimentacionSeccionMM2 float64
	ConductorTierraSeccionMM2       float64
	HilosPorFase                    int
	TipoCanalizacion                string
}

// ResultadoAjusteCorriente contiene el resultado del ajuste de corriente.
type ResultadoAjusteCorriente struct {
	CorrienteAjustada        float64 `json:"corriente_ajustada"`
	FactorTemperatura        float64 `json:"factor_temperatura"`
	FactorAgrupamiento       float64 `json:"factor_agrupamiento"`
	FactorUso                float64 `json:"factor_uso"`
	FactorTotal              float64 `json:"factor_total"`
	Temperatura              int     `json:"temperatura"`
	ConductoresPorTubo       int     `json:"conductores_por_tubo"`
	CantidadConductoresTotal int     `json:"cantidad_conductores_total"`
	TemperaturaAmbiente      int     `json:"temperatura_ambiente"`
}

// ResultadoCorriente contains the result of the current calculation.
type ResultadoCorriente struct {
	CorrienteNominal float64 `json:"corriente_nominal"`
}

// MemoriaOutput contiene el resultado completo de la memoria de cálculo.
type MemoriaOutput struct {
	// ═══════════════════════════════════════════════════════════════════════
	// DATOS DEL EQUIPO (reflejan el input)
	// ═══════════════════════════════════════════════════════════════════════
	Equipo DatosEquipo `json:"equipo"`

	// Información del cálculo
	TipoEquipo          string           `json:"tipo_equipo"`
	Tension             int              `json:"tension"`
	FactorPotencia      float64          `json:"factor_potencia"`
	Estado              string           `json:"estado"`
	TemperaturaAmbiente int              `json:"temperatura_ambiente"`
	SistemaElectrico    SistemaElectrico `json:"sistema_electrico"`
	CantidadConductores int              `json:"cantidad_conductores"`

	// Factores calculados
	FactorTemperaturaCalculado  float64 `json:"factor_temperatura_calculato"`
	FactorAgrupamientoCalculado float64 `json:"factor_agrupamiento_calculado"`

	// Paso 1: Corriente Nominal
	CorrienteNominal float64 `json:"corriente_nominal"`

	// Paso 2: Ajuste de Corriente
	FactorAgrupamiento float64 `json:"factor_agrupamiento"`
	FactorTemperatura  float64 `json:"factor_temperatura"`
	FactorTotalAjuste  float64 `json:"factor_total_ajuste"`
	CorrienteAjustada  float64 `json:"corriente_ajustada"`
	HilosPorFase       int     `json:"hilos_por_fase"`
	CorrientePorHilo   float64 `json:"corriente_por_hilo"`
	ConductoresPorTubo int     `json:"conductores_por_tubo"`

	// Paso 3: Tipo de Canalización
	TipoCanalizacion string `json:"tipo_canalizacion"`

	// Material del conductor
	Material string `json:"material"`

	// Paso 4: Conductor de Alimentación
	TemperaturaUsada      int                `json:"temperatura_usada"`
	ConductorAlimentacion ResultadoConductor `json:"conductor_alimentacion"`
	TablaAmpacidadUsada   string             `json:"tabla_ampacidad_usada"`

	// Paso 5: Conductor de Tierra
	ConductorTierra ResultadoConductor `json:"conductor_tierra"`
	ITM             int                `json:"itm"`

	// Paso 6: Canalización
	Canalizacion   ResultadoCanalizacion `json:"canalizacion"`
	FillFactor     float64               `json:"fill_factor"`
	DetalleCharola *DetalleCharola       `json:"detalle_charola,omitempty"`
	DetalleTuberia *DetalleTuberia       `json:"detalle_tuberia,omitempty"`

	// Paso 7: Caída de Tensión
	LongitudCircuito float64               `json:"longitud_circuito"`
	CaidaTension     ResultadoCaidaTension `json:"caida_tension"`

	// Resumen de cumplimiento
	CumpleNormativa bool     `json:"cumple_normativa"`
	Observaciones   []string `json:"observaciones"`

	// Todos los pasos para el reporte
	Pasos []PasoMemoria `json:"pasos"`
}
