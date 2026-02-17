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

## Flujo Obligatorio

invoke-agent: orchestrating-agents

context:
mandatory_skills:

- brainstorming
  architecture: "clean-ddd-hexagonal-vertical"
  require_full_flow: true
  requirement: "{{lo_que_el_usuario_solicita}}"

## Secuencia Forzada del Orquestador

1. **Revisar Planes Pendientes** - Ver si hay planes en `docs/plans/` que ya están completados
2. Invocar skill `brainstorming`
3. Generar diseño técnico
4. Generar plan de implementación
5. Crear rama de trabajo
6. Despachar agentes en orden:
   - domain-agent
   - application-agent
   - infrastructure-agent
7. Realizar wiring en main.go
8. **Verificación post-wiring**: `go build ./...` + `go test ./...`
9. **Pruebas manuales del endpoint** (si es API)
10. **Auditoría de código** (OBLIGATORIO) - Invocar auditores por capa
11. **Mover planes completados** a `docs/plans/completed/`
12. Invocar skill `finishing-a-development-branch`

## Reglas

- Nunca saltar brainstorming
- Nunca ejecutar agentes sin diseño previo
- Nunca hacer wiring antes de terminar infraestructura
- Siempre ejecutar verificación post-wiring
- Siempre hacer pruebas manuales para APIs
- **Siempre hacer auditoría de código antes de finalizar**
- Siempre mover planes completados a `completed/` al final
- **finishing-a-development-branch аудит AGENTS.md** (no duplicar)
- Mantener separación estricta de capas
