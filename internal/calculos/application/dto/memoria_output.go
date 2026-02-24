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
	AreaTotalMM2     float64
	AreaRequeridaMM2 float64
	NumeroDeTubos    int
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
	Canalizacion ResultadoCanalizacion `json:"canalizacion"`
	FillFactor   float64               `json:"fill_factor"`

	// Paso 7: Caída de Tensión
	LongitudCircuito float64               `json:"longitud_circuito"`
	CaidaTension     ResultadoCaidaTension `json:"caida_tension"`

	// Resumen de cumplimiento
	CumpleNormativa bool     `json:"cumple_normativa"`
	Observaciones   []string `json:"observaciones"`

	// Todos los pasos para el reporte
	Pasos []PasoMemoria `json:"pasos"`
}
