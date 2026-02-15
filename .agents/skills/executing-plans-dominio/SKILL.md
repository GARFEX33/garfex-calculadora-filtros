---
name: executing-plans-dominio
description: "Usar cuando tengas un plan de implementación de dominio escrito para ejecutarlo en una sesión separada con puntos de revisión"
---

# Ejecución de Planes de Dominio

## Descripción General

Cargar un plan de **dominio** aprobado, revisarlo críticamente, ejecutar tareas de dominio en lotes, y reportar para revisión entre cada lote. **Solo se implementa dominio**, nada de infraestructura, UI o base de datos.

**Principio central:** Ejecución por lotes con puntos de control para revisión de modelo de dominio, siguiendo la filosofía **fail-fast, batch incremental y verificación de invariantes**.

**Anunciar al inicio:** "Estoy usando la skill executing-plans-dominio para implementar este plan de dominio siguiendo la filosofía **fail-fast, batch incremental y verificación de invariantes**."

## El Proceso

### Paso 1: Cargar y Revisar Plan de Dominio

1. Leer el archivo del plan de dominio (`docs/plans/YYYY-MM-DD-<tema>-domain-implementation-plan.md`)
2. Revisarlo críticamente — identificar preguntas o preocupaciones sobre reglas de negocio, entidades, agregados o value objects
3. Si hay preocupaciones: plantearlas al usuario antes de comenzar
4. Si no hay preocupaciones: crear lista de tareas `TodoWrite` y continuar

### Paso 2: Ejecutar Lote de Tareas de Dominio

**Por defecto: primeras 3 tareas**

Para cada tarea de dominio:

1. Marcar como `in_progress`
2. Seguir cada paso exactamente (el plan tiene pasos detallados de dominio)
3. **Verificar criterios de dominio antes de continuar (fail-fast, Design by Contract, testeo incremental). Si la verificación falla, detenerse y pedir aclaración.**
4. Marcar como `completed`

### Paso 3: Reportar

Al terminar el lote:

- Mostrar qué tareas de dominio se implementaron
- Mostrar resultados de las verificaciones de dominio
- Decir: "Listo para feedback"

### Paso 4: Continuar

Basado en feedback:

- Aplicar cambios si es necesario
- Ejecutar siguiente lote
- Repetir hasta completar todas las tareas de dominio

### Paso 5: Finalizar

Después de completar todas las tareas de dominio y verificarlas:

- Anunciar: "Estoy usando la skill finishing-a-development-branch para completar este trabajo de dominio."
- **SUB-SKILL REQUERIDA:** superpowers:finishing-a-development-branch
- Seguir esa skill para verificar tests, presentar opciones y ejecutar elección

## Cuándo Detenerse y Pedir Ayuda

- Bloqueo a mitad de lote (dependencia de dominio no clara, verificación falla, instrucción poco clara)
- Plan de dominio incompleto
- No se entiende una instrucción de dominio
- Verificación de invariantes falla repetidamente

**Siempre pedir aclaración en lugar de asumir.**

## Cuándo Revisar Pasos Anteriores

- Usuario actualiza plan de dominio basado en feedback
- Enfoque fundamental de dominio necesita replantearse

## Principios Clave

- Revisar plan críticamente antes de ejecutar
- Seguir pasos de dominio exactamente
- **Verificar cada paso antes de continuar (fail-fast + batch incremental + invariantes)**
- No omitir verificaciones de reglas de negocio
- Entre lotes: solo reportar y esperar
- Detenerse ante bloqueos, no asumir
- Nunca implementar en main/master sin consentimiento explícito
- Prohibido cualquier mecanismo de worktree o espacios de trabajo aislados

## Integración

**Skills requeridas:**

- **writing-plans-dominio** — crea el plan que esta skill ejecuta
