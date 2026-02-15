---
name: writing-plans-dominio
description: "Usar cuando tengas un diseño de dominio aprobado para crear un plan detallado de implementación de dominio. Solo escribe planes, no codifiques ni implementes."
---

# Creación de Planes de Implementación de Dominio

## Descripción General

Esta skill toma un **diseño de dominio aprobado** y genera un **plan de implementación detallado**, dividido en tareas concretas y secuenciales, listo para ser ejecutado por `executing-plans-dominio`. **No escribir código ni estructuras de proyecto**, solo planificar los pasos de implementación de dominio.

**Principio central:** Generar un plan claro, secuencial y verificable, enfocado estrictamente en el dominio.

**Anunciar al inicio:** "Estoy usando la skill writing-plans-dominio para crear un plan de implementación de dominio."

## Flujo del Proceso

### Paso 1: Analizar Diseño de Dominio

1. Leer el diseño aprobado (`docs/plans/YYYY-MM-DD-<tema>-domain-design.md`).
2. Identificar agregados, entidades, value objects y reglas de negocio.
3. Detectar dependencias y relaciones entre componentes de dominio.

### Paso 2: Descomponer en Tareas de Dominio

1. Crear tareas separadas para cada entidad, agregado, value object o regla de negocio a implementar.
2. Cada tarea debe incluir:
   - Nombre descriptivo
   - Objetivo de la tarea
   - Pasos concretos y secuenciales para su implementación en dominio
   - Criterios de verificación o tests de dominio (sin implementarlos aún)

3. Priorizar tareas según dependencias de agregados o reglas.

### Paso 3: Revisar Secuencia y Complejidad

1. Asegurarse de que el plan respeta la integridad del modelo de dominio.
2. Ajustar secuencia para que cada tarea tenga sentido y pueda ejecutarse en lotes por `executing-plans-dominio`.
3. Identificar puntos de control para revisión del usuario.

### Paso 4: Presentar Plan de Dominio

1. Mostrar la lista de tareas con pasos detallados.
2. Solicitar aprobación incremental del usuario antes de finalizar el plan.
3. Permitir ajustes según feedback.

### Paso 5: Guardar Plan

1. Guardar plan en `docs/plans/YYYY-MM-DD-<tema>-domain-implementation-plan.md`.
2. Hacer commit del archivo para referencia futura.

### Principios Clave

- **Una tarea a la vez** — No abrumar con muchas acciones simultáneas.
- **Secuencia lógica y verificable** — Cada tarea debe ser clara y ejecutable sin suposiciones.
- **Solo dominio** — No incluir infraestructura, UI, bases de datos ni código.
- **Validación incremental** — Presentar el plan para aprobación antes de finalizar.
- **Preparación para ejecución** — El plan debe ser directamente ejecutable por `executing-plans-dominio`.
