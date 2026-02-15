---
name: writing-plans-application
description: "Usar cuando tengas un diseño de aplicación aprobado para crear un plan detallado de implementación de la capa de aplicación. Solo escribe planes, no codifiques ni implementes, y nunca planifiques dominio."
---

# Creación de Planes de Implementación de Aplicación

## Descripción General

Esta skill toma un **diseño de aplicación aprobado** y genera un **plan de implementación detallado**, dividido en tareas concretas y secuenciales, listo para ser ejecutado por `executing-plans-application`. **No escribir código de dominio ni estructuras internas de entidades**, solo planificar la lógica, orquestación y flujos de la capa de aplicación.

**Principio central:** Generar un plan claro, secuencial y verificable, enfocado estrictamente en la capa de aplicación.

**Anunciar al inicio:** "Estoy usando la skill writing-plans-application para crear un plan de implementación de aplicación."

## Flujo del Proceso

### Paso 1: Analizar Diseño de Aplicación

1. Leer el diseño aprobado (`docs/plans/YYYY-MM-DD-<tema>-application-design.md`).
2. Identificar casos de uso, servicios de aplicación, DTOs y flujos de orquestación.
3. Detectar dependencias con interfaces de dominio ya definidas y otras capas de infraestructura.

> ⚠️ Si falta información de dominio para implementar un caso de uso, **detenerse y solicitarla al equipo de dominio**.

### Paso 2: Descomponer en Tareas de Aplicación

1. Crear tareas separadas para cada caso de uso, servicio de aplicación o flujo de orquestación.

2. Cada tarea debe incluir:
   - Nombre descriptivo
   - Objetivo de la tarea
   - Pasos concretos y secuenciales para su implementación en la capa de aplicación
   - Criterios de verificación (tests de integración de aplicación o mocks de dominio si aplica)

3. Priorizar tareas según dependencias de flujos y servicios de dominio.

### Paso 3: Revisar Secuencia y Complejidad

1. Asegurarse de que el plan respeta la integridad de los casos de uso y flujos de aplicación.
2. Ajustar secuencia para que cada tarea tenga sentido y pueda ejecutarse en lotes por `executing-plans-application`.
3. Identificar puntos de control para revisión del usuario y coordinación con el equipo de dominio.

### Paso 4: Presentar Plan de Aplicación

1. Mostrar la lista de tareas con pasos detallados.
2. Solicitar aprobación incremental del usuario antes de finalizar el plan.
3. Permitir ajustes según feedback, **sin tocar nada de dominio**.

### Paso 5: Guardar Plan

1. Guardar plan en `docs/plans/YYYY-MM-DD-<tema>-application-implementation-plan.md`.
2. Hacer commit del archivo para referencia futura.

### Principios Clave

- **Una tarea a la vez** — No abrumar con muchas acciones simultáneas.
- **Secuencia lógica y verificable** — Cada tarea debe ser clara y ejecutable sin suposiciones.
- **Solo aplicación** — No incluir código de dominio, entidades o reglas de negocio.
- **Dependencia de dominio explícita** — Solicitar información de dominio cuando sea necesario, no asumir nada.
- **Validación incremental** — Presentar el plan para aprobación antes de finalizar.
- **Preparación para ejecución** — El plan debe ser directamente ejecutable por `executing-plans-application`.
