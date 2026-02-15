# Ejemplo: Orquestar Agentes para Nueva Feature

Este ejemplo muestra el flujo completo de cómo usar el comando `orchestrate-agents` para implementar una nueva feature.

## Escenario

Queremos agregar una nueva feature llamada `proyectos` que permita guardar memorias de cálculo en proyectos con nombre, cliente, y fecha.

## Paso 1: Brainstorming (Coordinador)

```
Usuario: Quiero agregar soporte para proyectos que agrupen memorias de cálculo

Coordinador: Invoca skill brainstorming
→ Diseño aprobado: docs/plans/2026-02-20-proyectos-design.md
```

## Paso 2: Writing Plans (Coordinador)

```
Coordinador: Invoca skill writing-plans
→ Plan creado: docs/plans/2026-02-20-proyectos-plan.md
```

Contenido típico del plan:

```markdown
# Plan: Proyectos Feature

## Paso 1 — Crear shared/kernel/ (si aplica)
Ya existe, no se necesita.

## Paso 2 — Crear proyectos/domain/
- entity/proyecto.go (Proyecto, IDProyecto)
- entity/errors.go
- repository/proyecto_repository.go (interface)

## Paso 3 — Crear proyectos/application/
- port/proyecto_repository.go
- usecase/crear_proyecto.go
- usecase/agregar_memoria.go
- dto/proyecto_input.go, proyecto_output.go

## Paso 4 — Crear proyectos/infrastructure/
- adapter/driven/postgres/postgres_proyecto_repository.go
- adapter/driver/http/proyecto_handler.go
- router.go

## Paso 5 — Wiring en cmd/api/main.go
- Instanciar repositorio
- Crear use cases
- Configurar router
```

## Paso 3: Crear Rama (Coordinador)

```bash
git checkout -b feature/proyectos
```

## Paso 4: Despachar Domain-Agent

**Comando:**
```bash
orchestrate-agents --agent domain --steps "2" --plan docs/plans/2026-02-20-proyectos-plan.md --feature proyectos
```

**Prompt que se envía:**

---

Sos el **domain-agent** de este proyecto. Tu trabajo es ejecutar **Paso 2** del plan.

## Proyecto

- Repositorio: `C:\PROGRAMACION\garfex-calculadora-filtros`
- Rama activa: `feature/proyectos`
- Módulo Go: `github.com/garfex/calculadora-filtros`

## Contexto — qué hicieron los agentes anteriores

Empezando desde cero. No hay agentes previos.

El shared/kernel ya existe con:
- valueobject.Corriente, Tension, Temperatura, etc.

## Tu scope

- `internal/proyectos/domain/entity/`
- `internal/proyectos/domain/repository/`

**NO toques:**
- `internal/proyectos/application/`
- `internal/proyectos/infrastructure/`
- `cmd/api/main.go`

## Plan a ejecutar

`docs/plans/2026-02-20-proyectos-plan.md`

## Instrucciones

1. Leé el plan completo
2. Creá tus propias tareas con TodoWrite antes de empezar
3. Ejecutá cada tarea marcando `in_progress` → `completed`
4. Verificá con `go test ./internal/proyectos/domain/...` antes de terminar
5. Si algo falla, arreglalo antes de seguir

## Al terminar

Reportá:
- Lista exacta de archivos creados/modificados
- Output de `go test ./internal/proyectos/domain/...`
- Issues encontrados (si hay)

---

**Respuesta esperada del agente:**

```
✅ Paso 2 completado

Archivos creados:
- internal/proyectos/domain/entity/proyecto.go
- internal/proyectos/domain/entity/errors.go
- internal/proyectos/domain/repository/proyecto_repository.go

Tests:
ok  github.com/garfex/calculadora-filtros/internal/proyectos/domain/entity
ok  github.com/garfex/calculadora-filtros/internal/proyectos/domain/repository

Sin issues.
```

## Paso 5: Despachar Application-Agent

**Comando:**
```bash
orchestrate-agents --agent application --steps "3" --plan docs/plans/2026-02-20-proyectos-plan.md --feature proyectos
```

**Prompt que se envía:**

---

Sos el **application-agent** de este proyecto. Tu trabajo es ejecutar **Paso 3** del plan.

## Proyecto

- Repositorio: `C:\PROGRAMACION\garfex-calculadora-filtros`
- Rama activa: `feature/proyectos`
- Módulo Go: `github.com/garfex/calculadora-filtros`

## Contexto — qué hizo domain-agent

Ya están creados y testeados:
- `internal/proyectos/domain/entity/proyecto.go` (Proyecto, IDProyecto)
- `internal/proyectos/domain/entity/errors.go`
- `internal/proyectos/domain/repository/proyecto_repository.go`

Los imports correctos que debés usar:
- Entities: `github.com/garfex/calculadora-filtros/internal/proyectos/domain/entity`

## Tu scope

- `internal/proyectos/application/port/`
- `internal/proyectos/application/usecase/`
- `internal/proyectos/application/dto/`

**NO toches:**
- `internal/proyectos/domain/`
- `internal/proyectos/infrastructure/`
- `cmd/api/main.go`

## Plan a ejecutar

`docs/plans/2026-02-20-proyectos-plan.md`

## Instrucciones

1. Leé el plan completo
2. Creá tus propias tareas con TodoWrite antes de empezar
3. Ejecutá cada tarea marcando `in_progress` → `completed`
4. Verificá con `go test ./internal/proyectos/application/...` antes de terminar
5. Si algo falla, arreglalo antes de seguir

## Al terminar

Reportá:
- Lista exacta de archivos creados/modificados
- Output de `go test ./internal/proyectos/application/...`
- Issues encontrados (si hay)

---

## Paso 6: Despachar Infrastructure-Agent

**Comando:**
```bash
orchestrate-agents --agent infrastructure --steps "4" --plan docs/plans/2026-02-20-proyectos-plan.md --feature proyectos
```

**Prompt que se envía:**

---

Sos el **infrastructure-agent** de este proyecto. Tu trabajo es ejecutar **Paso 4** del plan.

## Proyecto

- Repositorio: `C:\PROGRAMACION\garfex-calculadora-filtros`
- Rama activa: `feature/proyectos`
- Módulo Go: `github.com/garfex/calculadora-filtros`

## Contexto — qué hicieron los agentes anteriores

Ya están creados y testeados:
- Domain completo (entity, repository interface)
- Application completo (ports, use cases: CrearProyecto, AgregarMemoria, DTOs)

Los ports que debés implementar están en:
- `internal/proyectos/application/port/proyecto_repository.go`

## Tu scope

- `internal/proyectos/infrastructure/adapter/driver/http/`
- `internal/proyectos/infrastructure/adapter/driven/postgres/`

**NO toches:**
- `internal/proyectos/domain/`
- `internal/proyectos/application/`
- `cmd/api/main.go`

## Plan a ejecutar

`docs/plans/2026-02-20-proyectos-plan.md`

## Instrucciones

1. Leé el plan completo
2. Creá tus propias tareas con TodoWrite antes de empezar
3. Ejecutá cada tarea marcando `in_progress` → `completed`
4. Verificá con `go test ./internal/proyectos/infrastructure/...` antes de terminar
5. Si algo falla, arreglalo antes de seguir

## Al terminar

Reportá:
- Lista exacta de archivos creados/modificados
- Output de `go test ./internal/proyectos/infrastructure/...`
- Issues encontrados (si hay)

---

## Paso 7: Coordinador Finaliza

```bash
# Actualizar cmd/api/main.go (coordinador)
# - Importar nuevos paquetes
# - Crear instancias de repositorio
# - Crear use cases
# - Configurar router

# Verificar todo
go test ./...
go build ./...
go vet ./...

# Commit
git add -A
git commit -m "feat: add proyectos feature with CRUD operations

- Create domain/ with Proyecto entity and repository interface
- Create application/ with use cases and DTOs
- Create infrastructure/ with Postgres repo and HTTP handlers
- Update main.go with dependency wiring
- All tests passing"
```

## Resultado Final

```
internal/
  proyectos/
    domain/
      entity/
        proyecto.go
        errors.go
      repository/
        proyecto_repository.go
    application/
      port/
        proyecto_repository.go
      usecase/
        crear_proyecto.go
        agregar_memoria.go
      dto/
        proyecto_input.go
        proyecto_output.go
    infrastructure/
      adapter/
        driver/http/
          proyecto_handler.go
        driven/postgres/
          postgres_proyecto_repository.go
```

## Checklist para Coordinador

- [ ] Brainstorming completado y diseño aprobado
- [ ] Plan de implementación escrito
- [ ] Rama creada (`feature/nombre`)
- [ ] Domain-agent despachado y reportó éxito
- [ ] Application-agent despachado y reportó éxito
- [ ] Infrastructure-agent despachado y reportó éxito
- [ ] `cmd/api/main.go` actualizado con wiring
- [ ] `go test ./...` pasa completamente
- [ ] `go build ./...` sin errores
- [ ] Commit realizado
