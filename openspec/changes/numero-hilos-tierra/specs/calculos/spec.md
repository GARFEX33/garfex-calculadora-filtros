# Delta for Cálculos — Número de Hilos de Tierra

## Overview

Este cambio elimina el valor hardcodeado de `1` hilo de tierra en el cálculo de dimensionamiento de canalización, implementando cálculo automático según las reglas de la normativa NOM.

**Tipo de cambio**: MODIFIED (reemplaza comportamiento hardcodeado existente)

## Contexto

Actualmente, en `calcular_tamanio_tuberia.go` (línea 76) y en el orquestador, el número de hilos de tierra está hardcodeado a `1`. Este valor no es correcto cuando hay más de 2 tubos en paralelo.

## ADDED Requirements

### Requirement: Cálculo Automático de Hilos de Tierra

El sistema DEBE calcular automáticamente el número de hilos de tierra según el tipo de canalización y la cantidad de tubos, sin requerir que el usuario envíe este valor.

#### Scenario: Charola — Un hilo de tierra

- GIVEN el tipo de canalización es "charola" (cable espaciado o triangular)
- WHEN se ejecuta el dimensionamiento de canalización
- THEN el sistema DEBE usar exactamente 1 hilo de tierra para el cálculo

#### Scenario: Tubería con 1-2 tubos — Un hilo de tierra

- GIVEN el tipo de canalización es "tubería"
- AND el número de tubos en paralelo es 1 o 2
- WHEN se ejecuta el dimensionamiento de canalización
- THEN el sistema DEBE usar exactamente 1 hilo de tierra para el cálculo

#### Scenario: Tubería con más de 2 tubos — Dos hilos de tierra

- GIVEN el tipo de canalización es "tubería"
- AND el número de tubos en paralelo es mayor a 2 (3, 4, 5, ...)
- WHEN se ejecuta el dimensionamiento de canalización
- THEN el sistema DEBE usar exactamente 2 hilos de tierra para el cálculo

### Requirement: Integración con Memoria de Cálculo

El sistema DEBE calcular el número de hilos de tierra durante la ejecución de la memoria de cálculo antes de llamar al servicio de dominio.

#### Scenario: Flujo de memoria de cálculo con tubería

- GIVEN el usuario ejecuta la memoria de cálculo con tipo canalización "tubería"
- AND especifica 3 tubos en paralelo
- WHEN se llega al paso de dimensionamiento de canalización
- THEN el sistema DEBE calcular 2 hilos de tierra
- AND pasar ese valor al servicio de dominio `CalcularTamanioTuberiaWithMultiplePipes`

#### Scenario: Flujo de memoria de cálculo con charola

- GIVEN el usuario ejecuta la memoria de cálculo con tipo canalización "charola"
- WHEN se llega al paso de dimensionamiento de canalización
- THEN el sistema DEBE calcular 1 hilo de tierra
- AND usar ese valor para el cálculo de charola

## MODIFIED Requirements

### Requirement: Parámetro tierass en CalcularTamanioTuberiaUseCase

El use case DEBE aceptar y utilizar el valor calculado de hilos de tierra en lugar del valor hardcodeado.

(Anteriormente: el número de tierras estaba hardcodeado a `1` en la línea 76)

#### Scenario: Use case recibe valor calculado

- GIVEN el use case `CalcularTamanioTuberiaUseCase` recibe un `TuberiaInput`
- WHEN ejecuta el cálculo de tamaño de tubería
- THEN el sistema DEBE pasar el número de tierras calculado al servicio de dominio `CalcularTamanioTuberiaWithMultiplePipes`
- AND NO debe usar un valor hardcodeado de 1

## Edge Cases

### Scenario: Número de tubos inválido (cero o negativo)

- GIVEN el número de tubos especificado es 0 o negativo
- WHEN se intenta ejecutar el dimensionamiento
- THEN el sistema DEBE retornar un error de validación

### Scenario: Número de tubos no especificado (default)

- GIVEN el número de tubos no está especificado en la entrada
- WHEN se ejecuta el dimensionamiento
- THEN el sistema DEBE usar 1 tubo como valor por defecto
- AND calcular 1 hilo de tierra (porque numTuberias ≤ 2)

## Error Handling

### Scenario: Error en cálculo de número de tierras

- GIVEN ocurre un error al calcular el número de tierras
- WHEN se ejecuta el dimensionamiento
- THEN el sistema DEBE retornar el error correspondiente
- AND no繼續 con el cálculo de canalización

## Implementation Notes (No incluir en specs)

> **Nota**: Las siguientes notas son solo para referencia del implementador y NO forman parte de las especificaciones.

- El dominio ya soporta el parámetro `tierras` en `CalcularTamanioTuberiaWithMultiplePipes`
- El DTO `TuberiaInput` no necesita cambios (el campo no se expone al usuario)
- La lógica de cálculo debe implementarse en `OrquestadorMemoriaCalculoUseCase` antes de construir el `TuberiaInput`
