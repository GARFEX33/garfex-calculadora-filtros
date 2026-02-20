# Sistema de Agentes Especializados

Cada capa tiene su propio agente especializado con su ciclo completo de trabajo. El coordinador delega a los agentes en este orden:

```
domain-agent → application-agent → infrastructure-agent
```

### Cuándo invocar cada agente

| Accion | Agente | Skills del agente |
| ------ | ------ | ----------------- |
| Crear/modificar entidades, value objects, servicios de dominio | `domain-agent` | `brainstorming-dominio` → `writing-plans-dominio` → `executing-plans-dominio` |
| Crear/modificar ports, use cases, DTOs | `application-agent` | `brainstorming-application` → `writing-plans-application` → `executing-plans-application` |
| Crear/modificar adapters, repositorios, handlers HTTP | `infrastructure-agent` | `brainstorming-infrastructure` → `writing-plans-infrastructure` → `executing-plans-infrastructure` |
| Actualizar/auditar archivos AGENTS.md y README.md | `agents-md-manager` | `agents-md-manager` |

### Reglas de delegación entre agentes

- **domain-agent** trabaja primero — no sabe que existen los otros agentes
- **application-agent** toma el output del domain-agent — no toca infraestructura
- **infrastructure-agent** toma el output del application-agent — implementa ports, no define reglas
- **Coordinador** es el único que conoce el orden y hace el wiring en `cmd/api/main.go`
- Cada agente crea sus propias tareas con TodoWrite antes de ejecutar
- Cada agente verifica con `go test` antes de entregar

### Flujo del Coordinador (este chat)

El coordinador orquesta TODO el trabajo. Los agentes especializados solo ejecutan su parte.

```
Usuario pide feature/cambio
         │
         ▼
┌─────────────────────────────────────────────┐
│         ORQUESTADOR (Coordinador)           │
│  1. Invocar skill `brainstorming`           │
│  2. Crear diseño + plan                     │
│  3. Crear rama de trabajo                   │
│  4. Despachar agentes en orden              │
│  5. Hacer wiring en main.go                 │
│  6. Auditar AGENTS.md con agents-md-manager │
│  7. Commit final                            │
└─────────────────────────────────────────────┘
         │
    ┌────┴────┬────────────┐
    ▼         ▼            ▼
domain-   application-  infrastructure-
agent     agent         agent
    │         │            │
    ▼         ▼            ▼
 dominio   aplicación   infraestructura
 completo  completa     completa
```

**Qué hace el orquestador (coordinador):**
- Brainstorming inicial con el usuario
- Crear documentos de diseño y plan
- Crear rama git para el trabajo
- Despachar cada agente con contexto completo
- Esperar que cada agente termine antes de despachar el siguiente
- Hacer el wiring final en `cmd/api/main.go`
- **Auditar AGENTS.md con agents-md-manager PRE-merge**
- Aplicar correcciones de documentación antes de mergear
- Commit y preparar para merge

> **Nota:** La documentación es parte de la "definition of done". Los cambios a AGENTS.md van en el mismo PR/feature, no después.

**Qué hace cada agente especializado:**
- Leer el plan que le corresponde
- Crear sus propias tareas con TodoWrite
- Ejecutar SOLO en su capa (domain, application, o infrastructure)
- Verificar con `go test` antes de terminar
- Reportar archivos creados y resultado de tests

---

## Template para despachar agente

```
Sos el {agente} de este proyecto. Tu trabajo es ejecutar {pasos} del plan.

## Proyecto
Repositorio: {ruta}
Rama: {rama}
Módulo Go: {modulo}

## Contexto — qué hicieron los agentes anteriores
{resumen de lo que ya existe}

## Tu scope
{carpetas que puede tocar}

**NO toques** {carpetas prohibidas}

## Plan a ejecutar
{ruta al plan}

## Instrucciones
1. Leé el plan y creá tus propias tareas con TodoWrite
2. Ejecutá cada tarea
3. Verificá con go test antes de terminar

## Al terminar
Reportá: archivos creados, output de tests, issues encontrados
```

> **Skill de referencia:** Ver `.agents/skills/orchestrating-agents/SKILL.md` para el proceso completo.
