# Delta Spec: desarrollo-formula-corriente-pdf

## Resumen Ejecutivo

Este cambio añade la visualización del desarrollo paso a paso de la fórmula de corriente nominal en el PDF de la memoria de cálculo. Actualmente, el PDF solo muestra el resultado final (Iₙ = X A) sin el desglose del cálculo. El objetivo es que el PDF muestre la fórmula utilizada, los valores de entrada, y el desarrollo completo del cálculo — igual que como se muestra en el frontend.

## 1. Análisis del Contexto Actual

### Estado Actual

- **Frontend (SeccionCorriente.svelte)**: Ya posee la lógica completa en la función `getInfoCalculo()` que retorna:
  - `tipo`: Descripción del tipo de cálculo
  - `formula`: Fórmula usada en formato legible
  - `desarrollo`: Array de pasos del cálculo
  - `valores`: Mapa de valores de referencia

- **Template PDF (seccion_corriente.html)**: Muestra resultado final y parámetros básicos, pero NO muestra el desarrollo paso a paso

- **DTO (memoria_output.go)**: Solo contiene `ResultadoCorriente.CorrienteNominal` — carece de campos para el desarrollo

### Gap Identificado

El frontend tiene toda la lógica de cálculo por tipo de equipo, pero el backend no provee los datos necesarios para el PDF. El template PDF no puede mostrar el desarrollo porque el DTO no se lo proporciona.

## 2. Nuevo Struct Go para el Desarrollo

### ADDED Requirements

### Requirement: DatosDesarrolloCorriente Struct

El DTO `MemoriaOutput` DEBE incluir un nuevo campo `DesarrolloCorriente` que contenga toda la información necesaria para mostrar el desarrollo de la fórmula en el PDF.

```go
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

// PasoDesarrollo representa un paso individual en el desarrollo del cálculo.
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
```

#### Scenario: Estructura vacía para filtro activo

- GIVEN MemoriaOutput con TipoEquipo = "FILTRO_ACTIVO" y CorrienteNominal = 60.5
- WHEN se serializa el campo DesarrolloCorriente
- THEN DesarrolloCorriente.TipoCalculo = "Amperaje directo"
- AND DesarrolloCorriente.FormulaUsada = "I = Iₙominal"
- AND DesarrolloCorriente.PasosDesarrollo tiene 1 elemento
- AND DesarrolloCorriente.PasosDesarrollo[0].Descripcion = "I = 60.50 A (dato del equipo)"
- AND DesarrolloCorriente.ValoresReferencia["Amperaje"] = "60.50 A"

#### Scenario: Estructura completa para transformador trifásico

- GIVEN MemoriaOutput con TipoEquipo = "TRANSFORMADOR", Tension = 480, CorrienteNominal = 120.5
- WHEN se calcula el campo DesarrolloCorriente
- THEN DesarrolloCorriente.TipoCalculo = "Desde KVA (Transformador)"
- AND DesarrolloCorriente.FormulaUsada = "I = KVA / (kV × √3)"
- AND DesarrolloCorriente.PasosDesarrollo tiene 3 elementos
- AND DesarrolloCorriente.PasosDesarrollo[0].Descripcion contiene "kVA"
- AND DesarrolloCorriente.PasosDesarrollo[1].Descripcion contiene divisor
- AND DesarrolloCorriente.PasosDesarrollo[2].Resultado = "I = 120.50 A"

## 3. Lógica de Cálculo por Tipo de Equipo

### ADDED Requirements

### Requirement: DesarrolloCorriente use case

El sistema DEBE calcular los campos de `DatosDesarrolloCorriente` basándose en el tipo de equipo y el sistema eléctrico. La lógica DEBE ser idéntica a la implementada en el frontend `SeccionCorriente.getInfoCalculo()`.

#### Scenario: FILTRO_ACTIVO - amperaje directo

- GIVEN TipoEquipo = "FILTRO_ACTIVO" y CorrienteNominal = 45.0 A
- WHEN se genera el desarrollo
- THEN tipo_calculo = "Amperaje directo"
- AND formula = "I = Iₙominal"
- AND pasos = ["I = 45.00 A (dato del equipo)"]
- AND valores["Tipo"] = "Filtro Activo (FP = 1.0)"

#### Scenario: TRANSFORMADOR - cálculo desde KVA

- GIVEN TipoEquipo = "TRANSFORMADOR", Tension = 480 V, CorrienteNominal = 100.0 A
- WHEN se genera el desarrollo
- THEN tipo_calculo = "Desde KVA (Transformador)"
- AND formula = "I = KVA / (kV × √3)"
- AND KVA calculado = (100 × 480 × 1.732) / 1000 = 83.14 kVA
- AND pasos[0] = "I = 83.14 kVA / (0.480 kV × 1.732)"
- AND pasos[1] = "I = 83.14 / 0.831"
- AND pasos[2] = "I = 100.00 A"
- AND valores["KVA"] = "83.14 kVA"
- AND valores["Voltaje"] = "480 V (0.480 kV)"

#### Scenario: FILTRO_RECHAZO - cálculo desde KVAR

- GIVEN TipoEquipo = "FILTRO_RECHAZO", Tension = 480 V, CorrienteNominal = 75.0 A
- WHEN se genera el desarrollo
- THEN tipo_calculo = "Desde KVAR (Filtro de Rechazo)"
- AND formula = "I = KVAR / (kV × √3)"
- AND KVAR calculado = (75 × 480 × 1.732) / 1000 = 62.35 kVAR
- AND pasos contienen el desarrollo completo
- AND valores["KVAR"] = "62.35 kVAR"

#### Scenario: CARGA trifásico (ESTRELLA)

- GIVEN TipoEquipo = "CARGA", SistemaElectrico = "ESTRELLA", Tension = 220 V, FactorPotencia = 0.85, CorrienteNominal = 45.08 A
- WHEN se genera el desarrollo
- THEN tipo_calculo = "Desde Potencia (Sistema Trifásico)"
- AND formula = "I = P / (V × √3 × cosθ)"
- AND Potencia calculada = (45.08 × 220 × 1.732 × 0.85) / 1000 = 14.50 kW
- AND pasos[0] = "P = 14.50 kW = 14500 W"
- AND pasos[1] = "I = 14500 / (220 × 1.732 × 0.85)"
- AND pasos[2] = "I = 14500 / 323.69"
- AND pasos[3] = "I = 45.08 A"
- AND valores["Potencia"] = "14.50 kW"
- AND valores["Factor de Potencia"] = "0.85"
- AND valores["Sistema"] = "ESTRELLA"

#### Scenario: CARGA trifásico (DELTA)

- GIVEN TipoEquipo = "CARGA", SistemaElectrico = "DELTA", Tension = 480 V, FactorPotencia = 0.90, CorrienteNominal = 30.0 A
- WHEN se genera el desarrollo
- THEN usa la fórmula trifásica "I = P / (V × √3 × cosθ)"
- AND valores["Sistema"] = "DELTA"

#### Scenario: CARGA monofásico/bifásico

- GIVEN TipoEquipo = "CARGA", SistemaElectrico = "MONOFASICO", Tension = 127 V, FactorPotencia = 0.80, CorrienteNominal = 20.0 A
- WHEN se genera el desarrollo
- THEN tipo_calculo = "Desde Potencia (Sistema Monofásico)"
- AND formula = "I = P / (V × cosθ)"
- AND Potencia calculada = (20 × 127 × 0.80) / 1000 = 2.03 kW
- AND pasos usan divisor simple (V × FP), no √3
- AND valores["Sistema"] = "MONOFASICO"

### Requirement: Integración con memoria calculation use case

El `OrquestadorMemoriaCalculoUseCase` DEBE invocar la lógica de desarrollo de corriente después de calcular la corriente nominal y almacenar el resultado en `MemoriaOutput.DesarrolloCorriente`.

#### Scenario: Memoria completa incluye desarrollo de corriente

- GIVEN todos los datos de entrada válidos para un cálculo completo
- WHEN se ejecuta OrquestadorMemoriaCalculoUseCase
- THEN el resultado incluye el campo DesarrolloCorriente populated
- AND todos los campos de DesarrolloCorriente tienen valores válidos (no vacíos)

## 4. Cambios en el Template HTML

### ADDED Requirements

### Requirement: Sección de desarrollo en template PDF

El template `seccion_corriente.html` DEBE incluir una sección que muestre el desarrollo paso a paso del cálculo, utilizando los datos del nuevo campo `DesarrolloCorriente`.

#### Scenario: Template renderiza desarrollo completo

- GIVEN MemoriaOutput con DesarrolloCorriente populated
- WHEN se renderiza el template seccion_corriente.html
- THEN muestra el bloque "Tipo de Cálculo" con el valor de TipoCalculo
- AND muestra la fórmula en formato destacado (fórmula usada)
- AND muestra cada paso del desarrollo en orden secuencial
- AND muestra los valores de referencia en una tabla o lista

#### Scenario: Template renderiza con estilos apropiados

- GIVEN DesarrolloCorriente con datos válidos
- WHEN se renderiza el template
- THEN la fórmula usa estilo monoespaciado
- AND cada paso del desarrollo se muestra con indentación o numeración
- AND los valores de referencia se organizan en formato de tabla o grid

#### Scenario: Template maneja valores missing gracefully

- GIVEN DesarrolloCorriente con algunos ValoresReferencia vacíos
- WHEN se renderiza el template
- THEN solo muestra los valores que tienen contenido
- AND no muestra campos vacíos o con valor ""

### MODIFIED Requirements

### Requirement: Tipo de cálculo del equipo

El bloque existente de "Tipo de Cálculo" DEBE cambiar de una serie de `{{if}}` anidados a usar los datos estructurados de `DesarrolloCorriente`.

(Anteriormente: lógica condicional manual para cada tipo de equipo)

#### Scenario: Template usa datos estructurados

- GIVEN DesarrolloCorriente.TipoCalculo = "Desde KVA (Transformador)"
- WHEN se renderiza el template
- THEN muestra "Tipo de Cálculo: Desde KVA (Transformador)"
- AND muestra "Fórmula: I = KVA / (kV × √3)"
- AND NO usa {{if eq .Memoria.TipoEquipo "TRANSFORMADOR"}}

## 5. Casos Edge

### Requirement: Manejo de valores nil o cero

El sistema DEBE manejar gracefully los casos donde algunos valores no aplican.

#### Scenario: FILTRO_ACTIVO sin amperaje en equipo

- GIVEN TipoEquipo = "FILTRO_ACTIVO" pero Equipo.Amperaje es 0 o nil
- WHEN se genera el desarrollo
- THEN usa CorrienteNominal como valor de amperaje
- AND muestra el valor calculado, no cero

#### Scenario: Factor de potencia no aplica

- GIVEN TipoEquipo = "FILTRO_ACTIVO"
- WHEN se generan los valores de referencia
- THEN NO incluye "Factor de Potencia" en ValoresReferencia
- AND el template no muestra campo de FP

#### Scenario: Sistema eléctrico no determinado

- GIVEN SistemaElectrico es vacío o inválido
- WHEN se genera el desarrollo
- THEN usa comportamiento por defecto (monofásico)
- AND el cálculo funciona sin crash

## 6. Especificación de Valores de Referencia por Tipo

| Tipo Equipo | Valores Incluidos |
|-------------|-------------------|
| FILTRO_ACTIVO | Amperaje, Tipo |
| TRANSFORMADOR | KVA, Voltaje, Fórmula |
| FILTRO_RECHAZO | KVAR, Voltaje, Fórmula |
| CARGA (trifásico) | Potencia, Voltaje, Factor de Potencia, Sistema, Fórmula |
| CARGA (monofásico) | Potencia, Voltaje, Factor de Potencia, Sistema, Fórmula |

## 7. Validación de Consistencia

### Requirement: Consistencia entre frontend y backend

Los valores calculados en Go DEBEN coincidir exactamente con los calculados en el frontend (redondeo a 2 decimales).

#### Scenario: Validación de precisión

- GIVEN los mismos valores de entrada: Tension=480, CorrienteNominal=120.5, FactorPotencia=0.85
- WHEN se calculan los valores en Go y en TypeScript
- THEN los valores de KVA/Potencia difieren en máximo 0.01
- AND ambos usan Math.sqrt(3) = 1.7320508075688772
- AND ambos usan toFixed(2) para display
