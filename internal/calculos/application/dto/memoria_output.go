// internal/calculos/application/dto/memoria_output.go
package dto

// PasoMemoria representa un paso individual del cálculo.
type PasoMemoria struct {
	Numero      int
	Nombre      string
	Descripcion string
	Resultado   interface{}
}

// ═══════════════════════════════════════════════════════════════════════════
// ESTRUCTURAS EXISTENTES (sin cambios)
// ═══════════════════════════════════════════════════════════════════════════

// PasoDesarrollo representa un paso individual en el desarrollo del cálculo de corriente.
type PasoDesarrollo struct {
	// Numero es el número secuencial del paso (1, 2, 3, ...).
	Numero int `json:"numero"`

	// Descripcion describe qué se está calculando en este paso.
	// Ejemplo: "P = 45.08 kW = 45080 W"
	Descripcion string `json:"descripcion"`

	// Resultado es el resultado parcial o final de este paso.
	// Ejemplo: "I = 45.08 A"
	Resultado string `json:"resultado"`
}

// DatosDesarrolloCorriente contiene el desarrollo paso a paso del cálculo de corriente nominal.
type DatosDesarrolloCorriente struct {
	// TipoCalculo describe el tipo de cálculo realizado.
	// Valores posibles: "Amperaje directo", "Desde KVA (Transformador)",
	// "Desde KVAR (Filtro de Rechazo)", "Desde Potencia (Sistema Trifásico)",
	// "Desde Potencia (Sistema Monofásico)"
	TipoCalculo string `json:"tipo_calculo"`

	// FormulaUsada contiene la fórmula en formato legible.
	// Ejemplo: "I = KVA / (kV × √3)"
	FormulaUsada string `json:"formula_usada"`

	// PasosDesarrollo contiene cada paso del cálculo en orden secuencial.
	// Cada paso incluye la descripción y el resultado parcial.
	PasosDesarrollo []PasoDesarrollo `json:"pasos_desarrollo"`

	// ValoresReferencia contiene valores clave usados en el cálculo.
	// Mapa de clave -> valor formateado.
	// Incluye: KVA, KVAR, Potencia, Voltaje, Factor de Potencia, Sistema.
	ValoresReferencia map[string]string `json:"valores_referencia"`
}

// ResultadoConductor contiene la información del conductor seleccionado.
type ResultadoConductor struct {
	Calibre         string  `json:"calibre"`
	Material        string  `json:"material"`
	SeccionMM2      float64 `json:"seccion_mm2"`
	TipoAislamiento string  `json:"tipo_aislamiento"`
	Capacidad       float64 `json:"capacidad"`
	NumHilos        int     `json:"num_hilos"` // Número de hilos de tierra (1 para charola/tubería≤2, 2 para tubería>2)

	// Selección por caída de tensión (NOM-001-SEDE)
	// SeleccionPorCaidaTension indica si el calibre fue aumentado por caída de tensión
	SeleccionPorCaidaTension bool `json:"seleccion_por_caida_tension"`
	// CalibreOriginalAmpacidad es el calibre que hubiera correspondido solo por ampacidad
	CalibreOriginalAmpacidad string `json:"calibre_original_ampacidad,omitempty"`
	// NotaSeleccion explica el motivo del aumento de calibre
	NotaSeleccion string `json:"nota_seleccion,omitempty"`
}

// ResultadoConductores contiene los conductores seleccionados.
type ResultadoConductores struct {
	Alimentacion ResultadoConductor
	Tierra       ResultadoConductor
	TablaUsada   string
}

// ResultadoConductorCaidaTension contiene el resultado de la selección de conductor
// por criterio de caída de tensión (NOM-001-SEDE).
// Se utiliza cuando el conductor seleccionado por ampacidad no cumple la caída de tensión.
type ResultadoConductorCaidaTension struct {
	// CalibreOriginal es el calibre seleccionado originalmente por ampacidad
	CalibreOriginal string `json:"calibre_original"`
	// CalibreSeleccionado es el calibre superior que cumple la caída de tensión
	CalibreSeleccionado string  `json:"calibre_seleccionado"`
	SeccionMM2          float64 `json:"seccion_mm2"`
	TipoAislamiento     string  `json:"tipo_aislamiento"`
	Capacidad           float64 `json:"capacidad"`
	// CaidaTension es el resultado de caída de tensión ya verificado con el nuevo calibre
	CaidaTension ResultadoCaidaTension `json:"caida_tension"`
	// Nota describe el motivo del aumento: "Calibre aumentado de X a Y por caída de tensión (NOM-001-SEDE)"
	Nota string `json:"nota"`
	// Cumple indica si se encontró un calibre que cumple. False si se agotaron los intentos.
	Cumple bool `json:"cumple"`
	// IntentosRealizados es el número de calibres probados antes de encontrar uno válido
	IntentosRealizados int `json:"intentos_realizados"`
}

// ResultadoCanalizacion contiene la información de la canalización.
type ResultadoCanalizacion struct {
	Tamano           string  `json:"tamano"`
	AnchoComercialMM float64 `json:"ancho_comercial_mm,omitempty"`
	AreaTotalMM2     float64 `json:"area_total_mm2"`
	AreaRequeridaMM2 float64 `json:"area_requerida_mm2"`
	NumeroDeTubos    int     `json:"numero_de_tubos"`
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
	FactorControl    float64 `json:"factor_control,omitempty"`

	// Diagrama contiene el SVG generado para la charola
	// Solo se填充 cuando se solicita explícitamente
	Diagrama *DiagramaCharola `json:"diagrama,omitempty"`
}

// DiagramaCharola contiene los datos del diagrama SVG de una charola.
type DiagramaCharola struct {
	// Posiciones de los conductores en el diagrama (coordenadas mm)
	Posiciones []ConductorPosicionDTO `json:"posiciones"`
	// AnchoOcupado es el ancho total ocupado por los conductores en mm
	AnchoOcupado float64 `json:"ancho_ocupado_mm"`
	// ViewBox es el string del viewBox SVG
	ViewBox string `json:"viewBox"`
	// Cotas contiene las líneas de dimensión del diagrama
	Cotas []LineaCotaDTO `json:"cotas"`
	// SVG es el string completo del SVG generado
	SVG string `json:"svg"`
}

// ConductorPosicionDTO es la versión DTO de la posición de conductor para JSON.
type ConductorPosicionDTO struct {
	CX       float64 `json:"cx"`
	CY       float64 `json:"cy"`
	Radio    float64 `json:"radio"`
	Color    string  `json:"color"`
	Etiqueta string  `json:"etiqueta"`
	Tipo     string  `json:"tipo"`
}

// LineaCotaDTO es la versión DTO de una línea de cota para JSON.
type LineaCotaDTO struct {
	X1            float64 `json:"x1"`
	Y1            float64 `json:"y1"`
	X2            float64 `json:"x2"`
	Y2            float64 `json:"y2"`
	Valor         float64 `json:"valor"`
	Texto         string  `json:"texto"`
	PosicionTexto string  `json:"posicionTexto"`
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
	NumTuberias       int `json:"num_tuberias"`         // Número de tubos en paralelo

	// Tubo seleccionado de la tabla NOM (Cap. 9)
	AreaOcupacionTuboMM2 float64 `json:"area_ocupacion_tubo_mm2"` // área_ocupacion_mm2 del CSV (40% interior ya aplicado)
	DesignacionMetrica   string  `json:"designacion_metrica"`     // ej: "63" → mostrar como "63 mm"
	FillFactor           float64 `json:"fill_factor"`             // 0.40 para >2 conductores

	// Dimensiones físicas del tubo (para visualización SVG del diagrama de arreglo)
	// Leídos de tuberia-pvc-dimensiones-fisicas.csv — referencia visual, no para cálculo NOM.
	DiametroInteriorMM float64 `json:"diametro_interior_mm"`
	DiametroExteriorMM float64 `json:"diametro_exterior_mm"`

	// Diagrama contiene el SVG generado para la tubería
	// Solo se llena cuando se solicita explícitamente
	Diagrama *DiagramaTuberia `json:"diagrama,omitempty"`
}

// DiagramaTuberia contiene los datos del diagrama SVG de una tubería.
type DiagramaTuberia struct {
	// Posiciones de los conductores en el diagrama (coordenadas mm)
	Posiciones []ConductorPosicionDTO `json:"posiciones"`
	// DiametroInterior es el diámetro interior del tubo en mm
	DiametroInterior float64 `json:"diametro_interior_mm"`
	// DiametroExterior es el diámetro exterior del tubo en mm
	DiametroExterior float64 `json:"diametro_exterior_mm"`
	// ViewBox es el string del viewBox SVG
	ViewBox string `json:"viewBox"`
	// SVG es el string completo del SVG generado
	SVG string `json:"svg"`
}

// ResultadoCaidaTension contiene el resultado del cálculo de caída.
type ResultadoCaidaTension struct {
	Porcentaje       float64 `json:"porcentaje"`
	CaidaVolts       float64 `json:"caida_volts"`
	Cumple           bool    `json:"cumple"`
	LimitePorcentaje float64 `json:"limite_porcentaje"`
	Impedancia       float64 `json:"impedancia"`
	Resistencia      float64 `json:"resistencia"`
	Reactancia       float64 `json:"reactancia"`
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

// ═══════════════════════════════════════════════════════════════════════════
// NUEVAS ESTRUCTURAS AGRUPADAS (Phase 1)
// ═══════════════════════════════════════════════════════════════════════════

// DatosInstalacion agrupa los parámetros de instalación ingresados por el usuario.
// Refleja los valores de entrada más algunos cálculos derivados del entorno.
type DatosInstalacion struct {
	// Tension es el voltaje de operación en volts.
	// Valor de entrada del usuario.
	Tension int `json:"tension"`

	// SistemaElectrico indica el tipo de sistema eléctrico (DELTA, ESTRELLA, BIFASICO, MONOFASICO).
	// Valor de entrada del usuario.
	SistemaElectrico SistemaElectrico `json:"sistema_electrico"`

	// TipoCanalizacion indica el tipo de canalización a utilizar (TUBERIA_PVC, CHAROLA_CABLE_ESPACIADO, etc.).
	// Valor de entrada del usuario.
	TipoCanalizacion string `json:"tipo_canalizacion"`

	// Material es el material del conductor (CU para cobre, AL para aluminio).
	// Valor de entrada del usuario.
	Material string `json:"material"`

	// LongitudCircuito es la longitud del circuito en metros.
	// Valor de entrada del usuario.
	LongitudCircuito float64 `json:"longitud_circuito"`

	// HilosPorFase es el número de conductores por fase en paralelo.
	// Valor de entrada del usuario (default: 1).
	HilosPorFase int `json:"hilos_por_fase"`

	// PorcentajeCaidaMaximo es el límite de caída de tensión permitido en porcentaje.
	// Valor de entrada del usuario (default: 3.0%).
	PorcentajeCaidaMaximo float64 `json:"porcentaje_caida_maximo"`
}

// DatosCorrientes agrupa todos los cálculos relacionados con corriente eléctrica,
// incluyendo la corriente nominal, ajustada, y los factores de corrección.
type DatosCorrientes struct {
	// CorrienteNominal es la corriente nominal calculada en el paso 1 (Step 1).
	// Calculada a partir de la potencia o amperaje del equipo.
	CorrienteNominal float64 `json:"corriente_nominal"`

	// CorrienteAjustada es la corriente ajustada por factores de corrección en el paso 2 (Step 2).
	// Calculada: CorrienteNominal / (FactorTemperatura × FactorAgrupamiento × FactorUso)
	CorrienteAjustada float64 `json:"corriente_ajustada"`

	// CorrientePorHilo es la corriente que circula por cada hilo cuando hay conductores en paralelo.
	// Calculada: CorrienteAjustada / HilosPorFase
	CorrientePorHilo float64 `json:"corriente_por_hilo"`

	// FactorTemperatura es el factor de corrección por temperatura ambiente.
	// Valor de las tablas NOM según temperatura ambiente y temperatura del cable.
	FactorTemperatura float64 `json:"factor_temperatura"`

	// FactorAgrupamiento es el factor de corrección por agrupamiento de circuitos.
	// Valor de las tablas NOM según cantidad de conductores.
	FactorAgrupamiento float64 `json:"factor_agrupamiento"`

	// FactorTotalAjuste es el factor total combinado de corrección.
	// Calculado: FactorTemperatura × FactorAgrupamiento × FactorUso
	FactorTotalAjuste float64 `json:"factor_total_ajuste"`

	// TemperaturaAmbiente es la temperatura ambiente en grados Celsius.
	// Valor de entrada del usuario (determinado por el estado de la República).
	TemperaturaAmbiente int `json:"temperatura_ambiente"`

	// TemperaturaReferencia es la temperatura de operación del cable seleccionada (60, 75 o 90°C).
	// Determinada por el tipo de aislamiento y las tablas NOM.
	TemperaturaReferencia int `json:"temperatura_referencia"`

	// ConductoresPorTubo es el número de conductores por tubo o canal.
	// Valor de entrada del usuario o calculado según el sistema eléctrico.
	ConductoresPorTubo int `json:"conductores_por_tubo"`

	// CantidadConductores es el total de conductores en el sistema.
	// Calculado según el sistema eléctrico (DELTA, ESTRELLA, etc.) y hilos por fase.
	CantidadConductores int `json:"cantidad_conductores"`

	// TablaAmpacidadUsada es la tabla NOM utilizada para la selección de ampacidad.
	// Ejemplo: "Tabla 310-16" o "Tabla 310-17".
	TablaAmpacidadUsada string `json:"tabla_ampacidad_usada"`
}

// DatosCanalizacionCompleta agrupa el resultado del dimensionamiento de canalización
// junto con los detalles de cálculo (charola o tubería).
type DatosCanalizacionCompleta struct {
	// Resultado contiene el resultado del dimensionamiento de canalización.
	Resultado ResultadoCanalizacion `json:"resultado"`

	// FillFactor es el factor de llenado calculado (área ocupante / área disponible).
	// Debe ser ≤ 0.40 para más de 2 conductores según NOM.
	FillFactor float64 `json:"fill_factor"`

	// DetalleCharola contiene los valores intermedios del cálculo de charola.
	// Es nil cuando la canalización es tubería.
	DetalleCharola *DetalleCharola `json:"detalle_charola,omitempty"`

	// DetalleTuberia contiene los valores intermedios del cálculo de tubería.
	// Es nil cuando la canalización es charola.
	DetalleTuberia *DetalleTuberia `json:"detalle_tuberia,omitempty"`
}

// DatosProteccion agrupa los datos de protección eléctrica del circuito.
type DatosProteccion struct {
	// ITM es el interruptor termomagnético en Amperes.
	// Valor de entrada del usuario o obtenido del catálogo de equipos.
	ITM int `json:"itm"`
}

// MemoriaOutput contiene el resultado completo de la memoria de cálculo.
// Estructura reorganizada por entidad de dominio en lugar de pasos secuenciales.
type MemoriaOutput struct {
	// ═══════════════════════════════════════════════════════════════════════
	// DATOS DEL EQUIPO
	// ═══════════════════════════════════════════════════════════════════════

	// Equipo contiene los datos del equipo/filtro del catálogo.
	// Valor de entrada del usuario (modo LISTADO).
	Equipo DatosEquipo `json:"equipo"`

	// TipoEquipo indica el tipo de equipo (FILTRO_ACTIVO, TRANSFORMADOR, FILTRO_RECHAZO, CARGA).
	// Valor de entrada del usuario.
	TipoEquipo string `json:"tipo_equipo"`

	// FactorPotencia es el factor de potencia del equipo (cosθ ∈ (0, 1]).
	// Valor de entrada del usuario (solo para modo MANUAL_POTENCIA).
	FactorPotencia float64 `json:"factor_potencia"`

	// Estado indica el estado de la República Mexicana para determinar la temperatura ambiente.
	// Valor de entrada del usuario.
	Estado string `json:"estado"`

	// ═══════════════════════════════════════════════════════════════════════
	// PARÁMETROS DE INSTALACIÓN
	// ═══════════════════════════════════════════════════════════════════════

	// Instalacion agrupa los parámetros de instalación ingresados por el usuario.
	Instalacion DatosInstalacion `json:"instalacion"`

	// ═══════════════════════════════════════════════════════════════════════
	// CÁLCULOS DE CORRIENTE
	// ═══════════════════════════════════════════════════════════════════════

	// Corrientes agrupa los datos de corrientes y factores de ajuste.
	Corrientes DatosCorrientes `json:"corrientes"`

	// DesarrolloCorriente contiene el desarrollo paso a paso del cálculo de corriente.
	// Calculado en el orquestador después de determinar la corriente nominal.
	DesarrolloCorriente *DatosDesarrolloCorriente `json:"desarrollo_corriente,omitempty"`

	// ═══════════════════════════════════════════════════════════════════════
	// CONDUCTORES
	// ═══════════════════════════════════════════════════════════════════════

	// CableFase es el conductor de alimentación (fase).
	// Resultado del paso 4 de selección de conductor.
	CableFase ResultadoConductor `json:"cable_fase"`

	// CableNeutro es el conductor neutro.
	// Es nil para sistemas DELTA (sin neutro).
	CableNeutro *ResultadoConductor `json:"cable_neutro,omitempty"`

	// CableTierra es el conductor de tierra.
	// Resultado del paso 5 de selección de conductor de tierra.
	CableTierra ResultadoConductor `json:"cable_tierra"`

	// ═══════════════════════════════════════════════════════════════════════
	// CANALIZACIÓN
	// ═══════════════════════════════════════════════════════════════════════

	// Canalizacion agrupa el resultado del dimensionamiento de canalización.
	Canalizacion DatosCanalizacionCompleta `json:"canalizacion"`

	// ═══════════════════════════════════════════════════════════════════════
	// PROTECCIÓN
	// ═══════════════════════════════════════════════════════════════════════

	// Proteccion agrupa los datos de protección eléctrica.
	Proteccion DatosProteccion `json:"proteccion"`

	// ═══════════════════════════════════════════════════════════════════════
	// CAÍDA DE TENSIÓN
	// ═══════════════════════════════════════════════════════════════════════

	// CaidaTension contiene el resultado del cálculo de caída de tensión.
	// Resultado del paso 7.
	CaidaTension ResultadoCaidaTension `json:"caida_tension"`

	// ═══════════════════════════════════════════════════════════════════════
	// RESUMEN Y METADATOS
	// ═══════════════════════════════════════════════════════════════════════

	// CumpleNormativa indica si la instalación cumple toda la normativa NOM.
	// Es true si todos los criterios de selección son válidos.
	CumpleNormativa bool `json:"cumple_normativa"`

	// Observaciones contiene notas y advertencias sobre la instalación.
	Observaciones []string `json:"observaciones"`

	// Pasos contiene el detalle de todos los pasos del cálculo para debugging.
	Pasos []PasoMemoria `json:"pasos"`
}
