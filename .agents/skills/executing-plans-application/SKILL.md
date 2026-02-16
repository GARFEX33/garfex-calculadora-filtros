---
name: executing-plans-application
description: "Usar cuando tengas un plan de implementación de aplicación escrito para ejecutarlo en una sesión separada con puntos de revisión"
---

# Ejecución de Planes de Aplicación

## Descripción General

Cargar un plan de **aplicación** aprobado, revisarlo críticamente, ejecutar tareas de aplicación en lotes, y reportar para revisión entre cada lote. **Solo se implementa lógica de aplicación**, nada de infraestructura, UI o base de datos.

**Principio central:** Ejecución por lotes con puntos de control para revisión de lógica de aplicación, siguiendo la filosofía **fail-fast, batch incremental y verificación de invariantes**.

**Anunciar al inicio:** "Estoy usando la skill executing-plans-application para implementar este plan de aplicación siguiendo la filosofía **fail-fast, batch incremental y verificación de invariantes**."

## El Proceso

### Paso 1: Cargar y Revisar Plan de Aplicación

1. Leer el archivo del plan de aplicación (`docs/plans/YYYY-MM-DD-<tema>-application-implementation-plan.md`)
2. Revisarlo críticamente — identificar preguntas o preocupaciones sobre casos de uso, servicios, reglas de negocio o integraciones de dominio
3. Si hay preocupaciones: plantearlas al usuario antes de comenzar
4. Si no hay preocupaciones: crear lista de tareas `TodoWrite` y continuar

### Paso 2: Ejecutar Lote de Tareas de Aplicación

**Por defecto: primeras 3 tareas**

Para cada tarea de aplicación:

1. Marcar como `in_progress`
2. Seguir cada paso exactamente (el plan tiene pasos detallados de aplicación)
3. **Verificar criterios de aplicación antes de continuar (fail-fast, Design by Contract, testeo incremental). Si la verificación falla, detenerse y pedir aclaración.**
4. Marcar como `completed`
5. **Después de CADA tarea completada:**
   - Anunciar: "Estoy usando la skill finishing-a-development-branch para completar esta tarea de aplicación."
   - **SUB-SKILL REQUERIDA:** superpowers:finishing-a-development-branch
   - Seguir esa skill para: verificar tests → auditar AGENTS.md → hacer commit de la tarea

### Paso 3: Reportar

Al terminar el lote:

- Mostrar qué tareas de aplicación se implementaron
- Mostrar resultados de las verificaciones de aplicación
- Decir: "Listo para feedback"

### Paso 4: Continuar

Basado en feedback:

- Aplicar cambios si es necesario
- Ejecutar siguiente lote
- Repetir hasta completar todas las tareas de aplicación

### Paso 5: Finalizar

Después de completar todas las tareas de aplicación:

- Verificar que todas las tareas estén commiteadas
- Anunciar: "Trabajo de aplicación completado. Todas las tareas fueron verificadas y commiteadas individualmente."

## Cuándo Detenerse y Pedir Ayuda

- Bloqueo a mitad de lote (dependencia de aplicación no clara, verificación falla, instrucción poco clara)
- Plan de aplicación incompleto
- No se entiende una instrucción de aplicación
- Verificación de criterios falla repetidamente

**Siempre pedir aclaración en lugar de asumir.**

## Cuándo Revisar Pasos Anteriores

- Usuario actualiza plan de aplicación basado en feedback
- Enfoque fundamental de aplicación necesita replantearse

## Principios Clave

- Revisar plan críticamente antes de ejecutar
- Seguir pasos de aplicación exactamente
- **Verificar cada paso antes de continuar (fail-fast + batch incremental + invariantes)**
- No omitir verificaciones de reglas de negocio o casos de uso
- Entre lotes: solo reportar y esperar
- Detenerse ante bloqueos, no asumir
- Nunca implementar en main/master sin consentimiento explícito
- Prohibido cualquier mecanismo de worktree o espacios de trabajo aislados

## Integración

**Skills requeridas:**

- **writing-plans-application** — crea el plan que esta skill ejecuta
