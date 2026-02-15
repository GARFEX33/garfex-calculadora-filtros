---
name: executing-plans
description: Usar cuando tengas un plan de implementación escrito para ejecutarlo en una sesión separada con puntos de revisión
---

# Ejecución de Planes

## Descripción General

Cargar el plan, revisarlo críticamente, ejecutar tareas en lotes y reportar para revisión entre cada lote.

**Principio central:** Ejecución por lotes con puntos de control para revisión del arquitecto.

**Anunciar al inicio:** "Estoy usando la skill executing-plans para implementar este plan."

## El Proceso

### Paso 1: Cargar y Revisar el Plan

1. Leer el archivo del plan
2. Revisarlo críticamente — identificar preguntas o preocupaciones sobre el plan
3. Si hay preocupaciones: plantearlas al usuario antes de comenzar
4. Si no hay preocupaciones: crear TodoWrite y continuar

### Paso 2: Ejecutar Lote

**Por defecto: Las primeras 3 tareas**

Para cada tarea:

1. Marcar como in_progress
2. Seguir cada paso exactamente (el plan tiene pasos pequeños y concretos)
3. Ejecutar las verificaciones especificadas
4. Marcar como completed

### Paso 3: Reportar

Cuando el lote esté completo:

- Mostrar qué se implementó
- Mostrar la salida de las verificaciones
- Decir: "Listo para feedback."

### Paso 4: Continuar

Basado en el feedback:

- Aplicar cambios si es necesario
- Ejecutar el siguiente lote
- Repetir hasta completar

### Paso 5: Finalizar el Desarrollo

Después de que todas las tareas estén completas y verificadas:

- Anunciar: "Estoy usando la skill finishing-a-development-branch para completar este trabajo."
- **SUB-SKILL REQUERIDA:** Usar superpowers:finishing-a-development-branch
- Seguir esa skill para verificar tests, presentar opciones y ejecutar la elección

## Cuándo Detenerse y Pedir Ayuda

**DETENER la ejecución inmediatamente cuando:**

- Se encuentra un bloqueo a mitad de lote (dependencia faltante, test falla, instrucción poco clara)
- El plan tiene vacíos críticos que impiden comenzar
- No se entiende una instrucción
- Una verificación falla repetidamente

**Pedir aclaración en lugar de asumir.**

## Cuándo Revisar Pasos Anteriores

**Volver a la Revisión (Paso 1) cuando:**

- El usuario actualiza el plan basado en tu feedback
- El enfoque fundamental necesita replantearse

**No forzar bloqueos** — detenerse y preguntar.

## Recordar

- Revisar el plan críticamente primero
- Seguir los pasos exactamente
- No omitir verificaciones
- Referenciar skills cuando el plan lo indique
- Entre lotes: solo reportar y esperar
- Detenerse ante bloqueos, no asumir
- Nunca iniciar implementación en main/master sin consentimiento explícito del usuario
- **Prohibido usar cualquier mecanismo de worktree, espacios de trabajo aislados basados en worktree o configuraciones derivadas de worktree en cualquier circunstancia**

## Integración

**Skills de flujo requeridas:**

- **superpowers:writing-plans** — Crea el plan que esta skill ejecuta
- **superpowers:finishing-a-development-branch** — Completa el desarrollo después de todas las tareas
