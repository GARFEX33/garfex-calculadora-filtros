---
name: writing-plans-infrastructure
description: "Usar cuando tengas un diseño de infraestructura aprobado para crear un plan detallado de implementación de la capa de infraestructura. Solo escribe planes, no codifiques ni implementes, y nunca planifiques dominio o aplicación."
---

# Creación de Planes de Implementación de Infraestructura

## Descripción General

Esta skill toma un **diseño de infraestructura aprobado** y genera un **plan de implementación detallado**, dividido en tareas concretas y secuenciales, listo para ser ejecutado por `executing-plans-infrastructure`. **No escribir código de dominio ni de aplicación**, solo planificar la provisión, configuración y orquestación de recursos de infraestructura.

**Principio central:** Generar un plan claro, secuencial y verificable, enfocado estrictamente en la capa de infraestructura.

**Anunciar al inicio:** "Estoy usando la skill writing-plans-infrastructure para crear un plan de implementación de infraestructura."

## Flujo del Proceso

### Paso 1: Analizar Diseño de Infraestructura

1. Leer el diseño aprobado (`docs/plans/YYYY-MM-DD-<tema>-infrastructure-design.md`).
2. Identificar recursos, servicios, redes, bases de datos, colas, almacenamiento y demás componentes de infraestructura.
3. Detectar dependencias entre recursos y posibles interacciones con la capa de aplicación (interfaces, endpoints, configuraciones).

> ⚠️ Si falta información de aplicación o dominio para decidir configuraciones críticas, **detenerse y solicitarla al equipo correspondiente**.

### Paso 2: Descomponer en Tareas de Infraestructura

1. Crear tareas separadas para cada recurso, servicio o configuración de infraestructura.

2. Cada tarea debe incluir:
   - Nombre descriptivo
   - Objetivo de la tarea
   - Pasos concretos y secuenciales para su implementación en la capa de infraestructura
   - Criterios de verificación (tests de despliegue, conectividad, permisos, backups)

3. Priorizar tareas según dependencias de recursos y servicios.

### Paso 3: Revisar Secuencia y Complejidad

1. Asegurarse de que el plan respeta la integridad de la infraestructura y dependencias críticas.
2. Ajustar secuencia para que cada tarea tenga sentido y pueda ejecutarse en lotes por `executing-plans-infrastructure`.
3. Identificar puntos de control para revisión del usuario y coordinación con equipos de dominio o aplicación si aplica.

### Paso 4: Presentar Plan de Infraestructura

1. Mostrar la lista de tareas con pasos detallados.
2. Solicitar aprobación incremental del usuario antes de finalizar el plan.
3. Permitir ajustes según feedback, **sin tocar nada de dominio o aplicación**.

### Paso 5: Guardar Plan

1. Guardar plan en `docs/plans/YYYY-MM-DD-<tema>-infrastructure-implementation-plan.md`.
2. Hacer commit del archivo para referencia futura.

### Principios Clave

- **Una tarea a la vez** — No abrumar con muchas acciones simultáneas.
- **Secuencia lógica y verificable** — Cada tarea debe ser clara y ejecutable sin suposiciones.
- **Solo infraestructura** — No incluir código de dominio o aplicación.
- **Dependencias explícitas** — Solicitar información de dominio o aplicación cuando sea necesario, no asumir nada.
- **Validación incremental** — Presentar el plan para aprobación antes de finalizar.
- **Preparación para ejecución** — El plan debe ser directamente ejecutable por `executing-plans-infrastructure`.
