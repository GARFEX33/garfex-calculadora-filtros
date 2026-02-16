---
name: executing-plans-infrastructure
description: "Usar cuando tengas un plan de implementación de infraestructura escrito para ejecutarlo en una sesión separada con puntos de revisión"
---

# Ejecución de Planes de Infraestructura

## Descripción General

Cargar un plan de **infraestructura** aprobado, revisarlo críticamente, ejecutar tareas de infraestructura en lotes, y reportar para revisión entre cada lote. **Solo se implementa infraestructura**, nada de dominio, UI o lógica de aplicación.

**Principio central:** Ejecución por lotes con puntos de control para revisión de infraestructura, siguiendo la filosofía **fail-fast, batch incremental y verificación de invariantes**.

**Anunciar al inicio:** "Estoy usando la skill executing-plans-infrastructure para implementar este plan de infraestructura siguiendo la filosofía **fail-fast, batch incremental y verificación de invariantes**."

## El Proceso

### Paso 1: Cargar y Revisar Plan de Infraestructura

1. Leer el archivo del plan de infraestructura (`docs/plans/YYYY-MM-DD-<tema>-infrastructure-implementation-plan.md`)
2. Revisarlo críticamente — identificar preguntas o preocupaciones sobre despliegue, entornos, configuraciones, dependencias o seguridad
3. Si hay preocupaciones: plantearlas al usuario antes de comenzar
4. Si no hay preocupaciones: crear lista de tareas `TodoWrite` y continuar

### Paso 2: Ejecutar Lote de Tareas de Infraestructura

**Por defecto: primeras 3 tareas**

Para cada tarea de infraestructura:

1. Marcar como `in_progress`
2. Seguir cada paso exactamente (el plan tiene pasos detallados de infraestructura)
3. **Verificar criterios de infraestructura antes de continuar (fail-fast, Design by Contract, testeo incremental). Si la verificación falla, detenerse y pedir aclaración.**
4. Marcar como `completed`
5. **Después de CADA tarea completada:**
   - Anunciar: "Estoy usando la skill finishing-a-development-branch para completar esta tarea de infraestructura."
   - **SUB-SKILL REQUERIDA:** superpowers:finishing-a-development-branch
   - Seguir esa skill para: verificar tests → auditar AGENTS.md → hacer commit de la tarea

### Paso 3: Reportar

Al terminar el lote:

- Mostrar qué tareas de infraestructura se implementaron
- Mostrar resultados de las verificaciones de infraestructura
- Decir: "Listo para feedback"

### Paso 4: Continuar

Basado en feedback:

- Aplicar cambios si es necesario
- Ejecutar siguiente lote
- Repetir hasta completar todas las tareas de infraestructura

### Paso 5: Finalizar

Después de completar todas las tareas de infraestructura:

- Verificar que todas las tareas estén commiteadas
- Anunciar: "Trabajo de infraestructura completado. Todas las tareas fueron verificadas y commiteadas individualmente."

## Cuándo Detenerse y Pedir Ayuda

- Bloqueo a mitad de lote (dependencia de infraestructura no clara, verificación falla, instrucción poco clara)
- Plan de infraestructura incompleto
- No se entiende una instrucción de infraestructura
- Verificación de criterios falla repetidamente

**Siempre pedir aclaración en lugar de asumir.**

## Cuándo Revisar Pasos Anteriores

- Usuario actualiza plan de infraestructura basado en feedback
- Enfoque fundamental de infraestructura necesita replantearse

## Principios Clave

- Revisar plan críticamente antes de ejecutar
- Seguir pasos de infraestructura exactamente
- **Verificar cada paso antes de continuar (fail-fast + batch incremental + invariantes)**
- No omitir verificaciones de despliegue, dependencias o seguridad
- Entre lotes: solo reportar y esperar
- Detenerse ante bloqueos, no asumir
- Nunca implementar en main/master sin consentimiento explícito
- Prohibido cualquier mecanismo de worktree o espacios de trabajo aislados

## Integración

**Skills requeridas:**

- **writing-plans-infrastructure** — crea el plan que esta skill ejecuta
