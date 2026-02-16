---
description: "Inicializa y orquesta cualquier nueva feature o cambio siguiendo el flujo obligatorio de coordinación."
disable-model-invocation: false
---

# Orquestación Obligatoria de Cambios

## Regla Crítica

SIEMPRE seguir este flujo sin excepción cuando el usuario solicite:

- Nueva feature
- Nueva función
- Cambio estructural
- Ajuste en arquitectura
- Modificación relevante

## Flujo Obligatorio

invoke-agent: orchestrating-agents

context:
mandatory_skills:

- brainstorming
  architecture: "clean-ddd-hexagonal-vertical"
  require_full_flow: true
  requirement: "{{lo_que_el_usuario_solicita}}"

## Secuencia Forzada del Orquestador

1. Invocar skill `brainstorming`
2. Generar diseño técnico
3. Generar plan de implementación
4. Crear rama de trabajo
5. Despachar agentes en orden:
   - domain-agent
   - application-agent
   - infrastructure-agent
6. Realizar wiring en main.go
7. Auditar AGENTS.md con agents-md-curator
8. Commit final

## Reglas

- Nunca saltar brainstorming
- Nunca ejecutar agentes sin diseño previo
- Nunca hacer wiring antes de terminar infraestructura
- Siempre auditar AGENTS.md antes del commit
- Mantener separación estricta de capas
