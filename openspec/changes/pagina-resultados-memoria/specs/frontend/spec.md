# Delta for Frontend - Página de Resultados de Memoria

## ADDED Requirements

### Requirement: Redirección a Página de Resultados

Al presionar el botón "Calcular Memoria", el sistema DEBE redirigir al usuario a la página de resultados dedicada en lugar de mostrar los resultados en la misma página.

#### Scenario: Redirección exitosa tras cálculo

- GIVEN el usuario ha completado el formulario con datos válidos
- WHEN el usuario hace clic en "Calcular Memoria" y el servidor retorna éxito
- THEN el sistema DEBE redirigir a `/calculos/resultado`
- AND los datos del cálculo DEBEN estar disponibles en la página de resultados

#### Scenario: Error en cálculo

- GIVEN el usuario envía el formulario y el servidor retorna un error
- THEN el sistema DEBE mostrar el mensaje de error en la página del formulario
- AND NO DEBE redirigir a la página de resultados

### Requirement: Página de Resultados - Encabezado

La página de resultados DEBE mostrar un encabezado profesional con los datos generales del cálculo.

#### Scenario: Encabezado con todos los datos

- GIVEN el usuario accede a la página de resultados con datos válidos
- THEN el sistema DEBE mostrar:
  - Nombre de empresa o "Memoria de Cálculo"
  - Título del documento identificando el tipo (ej: "Memoria de Cálculo de Alimentador")
  - Proyecto (si está disponible)
  - Equipo o carga calculada
  - Capacidad (kVA o amperaje)
  - Voltaje de operación (V)
  - Longitud del alimentador (m)

### Requirement: Página de Resultados - Cálculo de Corriente Nominal

La página DEBE mostrar el cálculo de corriente nominal con fórmula y desarrollo.

#### Scenario: Cálculo de corriente para sistema trifásico

- GIVEN existen datos de potencia o amperaje nominal
- WHEN se muestra la sección de corriente nominal
- THEN el sistema DEBE mostrar:
  - La fórmula general (I = S / (√3 × V) o I = P / (V × √3 × cosθ))
  - La sustitución con valores numéricos
  - El resultado en amperes con unidad

#### Scenario: Cálculo para sistema monofásico

- GIVEN el sistema eléctrico es MONOFASICO
- WHEN se muestra la sección de corriente nominal
- THEN el sistema DEBE usar la fórmula monofásica (I = S / V)

### Requirement: Página de Resultados - Dimensionamiento del Alimentador

La página DEBE mostrar el dimensionamiento del conductor de alimentación citando la norma NOM.

#### Scenario: Dimensionamiento con factor de diseño

- GIVEN se ha calculado la corriente nominal
- WHEN se muestra la sección de dimensionamiento
- THEN el sistema DEBE mostrar:
  - Referencia al Artículo 310-15(b)(17) de NOM-001-SEDE-2012
  - Factor de diseño aplicado (125% o 135% según criterio)
  - Fórmula: I_diseño = Factor × I_nominal
  - Sustitución y resultado
  - Conductor seleccionado con:
    - Calibre (AWG/kcmil)
    - Material (Cu/Al)
    - Tipo de aislamiento
    - Temperatura de referencia
    - Ampacidad (A)
  - Justificación técnica

### Requirement: Página de Resultados - Conductor de Puesta a Tierra

La página DEBE mostrar el dimensionamiento del conductor de tierra conforme a NOM.

#### Scenario: Selección de conductor de tierra

- GIVEN se conoce la corriente del dispositivo de protección (ITM)
- WHEN se muestra la sección de puesta a tierra
- THEN el sistema DEBE mostrar:
  - Referencia a Tabla 250-122 de NOM-001-SEDE-2012
  - ITM utilizado
  - Calibre seleccionado
  - Justificación

### Requirement: Página de Resultados - Cálculo de Caída de Tensión

La página DEBE mostrar el cálculo de caída de tensión con todos los parámetros.

#### Scenario: Cálculo de caída de tensión trifásico

- GIVEN existen datos del conductor, longitud, corriente y tensión
- WHEN se muestra la sección de caída de tensión
- THEN el sistema DEBE mostrar:
  - Fórmula trifásica: %V = (√3 × I × Z × L / V) × 100
  - Desarrollo de la fórmula
  - Sustitución con valores numéricos
  - Resultado en porcentaje
  - Verificación contra límite (típicamente 3%)
  - Si aplica: cálculo de DMG, RMG, inductancia e impedancia

#### Scenario: Cálculo de caída de tensión monofásico

- GIVEN el sistema es monofásico
- WHEN se muestra la sección de caída de tensión
- THEN el sistema DEBE usar la fórmula monofásica: %V = (2 × I × Z × L / V) × 100

### Requirement: Página de Resultados - Cálculo de Canalización

La página DEBE mostrar el dimensionamiento de la canalización.

#### Scenario: Dimensionamiento de tubería

- GIVEN se conoce el tamaño y cantidad de conductores
- WHEN se muestra la sección de canalización
- THEN el sistema DEBE mostrar:
  - Tipo de arreglo
  - Fórmula empleada
  - Sustitución de valores
  - Tamaño comercial seleccionado

#### Scenario: Dimensionamiento de charola

- GIVEN el tipo de canalización es CHAROLA
- WHEN se muestra la sección de canalización
- THEN el sistema DEBE mostrar cálculo de charola si aplica

### Requirement: Página de Resultados - Conclusión Técnica

La página DEBE incluir una conclusión técnica que resuma el cumplimiento normativo.

#### Scenario: Sistema cumple todos los criterios

- GIVEN todos los cálculos cumplen los criterios de la norma
- WHEN se muestra la conclusión
- THEN el sistema DEBE indicar:
  - Que se aplicaron las fórmulas correctas
  - Que el sistema cumple criterios de ampacidad
  - Que la caída de tensión está dentro de límites
  - Que el diseño cumple con NOM-001-SEDE-2012

#### Scenario: Sistema no cumple algún criterio

- GIVEN algún cálculo no cumple los criterios
- WHEN se muestra la conclusión
- THEN el sistema DEBE indicar qué criterios no se cumplieron

### Requirement: Navegación de Regreso

La página de resultados DEBE permitir al usuario regresar al formulario.

#### Scenario: Regresar al formulario

- GIVEN el usuario está en la página de resultados
- WHEN el usuario hace clic en "Nuevo Cálculo" o botón de regresar
- THEN el sistema DEBE redirigir a la página del formulario

## MODIFIED Requirements

### Requirement: Flujo de Cálculo de Memoria (Modificado)

El flujo de cálculo de memoria DEBE ahora mostrar resultados en página dedicada.

(Anteriormente: Los resultados se mostraban inline en la misma página del formulario)

#### Scenario: Flujo modificado

- GIVEN el usuario completa el formulario
- WHEN hace clic en "Calcular Memoria"
- THEN los resultados se muestran en `/calculos/resultado` (no inline)

## REMOVED Requirements

### Requirement: Visualización Inline de Resultados (Eliminado)

La visualización inline de resultados en la página principal ya no aplica.

(Razón: Se迁移 a página dedicada para mejor presentación técnica)

## Technical Notes

- Los datos del cálculo se pasan via SvelteKit navigation state o query params
- Usar los tipos existentes de `MemoriaOutput` del frontend
- Mantener diseño responsive con tokens de Tailwind existentes
- El backend NO se modifica (el endpoint ya existe y funciona)
