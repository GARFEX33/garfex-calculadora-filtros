# Delta for Tension (Value Object)

## Purpose

Agregar soporte para unidades de voltaje (V y kV) al value object Tension, permitiendo que la API acepte voltajes en kilovoltios además de volts.

## ADDED Requirements

### Requirement: El value object Tension debe soportar unidades V y kV

El sistema DEBE aceptar un valor numérico junto con una unidad (V o kV) y normalizar internamente a volts.

#### Scenario: Crear tensión en volts

- GIVEN un valor de voltaje válido según NOM (127, 220, 240, 277, 440, 480, 600)
- WHEN se llama `NewTension(480, "V")`
- THEN retorna un Tension con valor 480 y unidad "V"
- AND `Valor()` retorna 480
- AND `Unidad()` retorna "V"
- AND `EnKilovoltios()` retorna 0.48

#### Scenario: Crear tensión en kilovoltios

- GIVEN un valor en kilovoltios que equivalga a un voltaje válido según NOM (0.48 kV = 480 V)
- WHEN se llama `NewTension(0.48, "kV")`
- THEN retorna un Tension con valor 480 y unidad "kV"
- AND `Valor()` retorna 480 (valor normalizado en volts)
- AND `Unidad()` retorna "kV"
- AND `EnKilovoltios()` retorna 0.48

#### Scenario: Crear tensión sin especificar unidad (default V)

- GIVEN un valor de voltaje válido según NOM
- WHEN se llama `NewTension(220, "")` o `NewTension(220)`
- THEN debe comportarse como si fuera "V" (compatibilidad hacia atrás)

### Requirement: Validación de valores NOM debe continuar funcionando

El sistema DEBE validar que el voltaje normalizado sea uno de los valores permitidos por NOM.

#### Scenario: Voltaje en kV que no corresponde a valor NOM válido

- GIVEN un valor de 0.5 kV (500 V, que NO está en lista NOM)
- WHEN se llama `NewTension(0.5, "kV")`
- THEN retorna error ErrVoltajeInvalido

#### Scenario: Voltaje en V que no corresponde a valor NOM válido

- GIVEN un valor de 230 V (NO está en lista NOM)
- WHEN se llama `NewTension(230, "V")`
- THEN retorna error ErrVoltajeInvalido

### Requirement: La API debe aceptar el campo tension_unidad

El sistema DEBE aceptar el campo opcional `tension_unidad` en el input JSON.

#### Scenario: Input con tensión en volts

- GIVEN un JSON con `tension: 480` y `tension_unidad: "V"`
- WHEN se procesa el input
- THEN se crea correctamente el value object Tension

#### Scenario: Input con tensión en kilovoltios

- GIVEN un JSON con `tension: 0.48` y `tension_unidad: "kV"`
- WHEN se procesa el input
- THEN se crea correctamente el value object Tension con valor 480

#### Scenario: Input sin tensión (compatibilidad hacia atrás)

- GIVEN un JSON con solo `tension: 480` sin `tension_unidad`
- WHEN se procesa el input
- THEN se usa "V" como unidad por defecto

## MODIFIED Requirements

### Requirement: Constructor NewTension

(New description — reemplaza el constructor actual)

El constructor DEBE aceptar dos parámetros: valor numérico y unidad como string. Si la unidad es vacía, debe asumir "V".

(Previously: `NewTension(valor int) (Tension, error)` - solo aceptaba valor entero)

#### Scenario: Constructor recibe valores válidos

- GIVEN un valor de 480 y unidad "V"
- WHEN se llama `NewTension(480, "V")`
- THEN retorna Tension{valor: 480, unidad: "V"}

## REMOVED Requirements

(Ninguno - es un cambio aditivo)
