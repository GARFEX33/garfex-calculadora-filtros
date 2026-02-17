# Plan de Implementación: Refactorización a Vertical Slices

**Fecha:** 2026-02-15  
**Diseño base:** `docs/plans/2026-02-15-vertical-slices-refactor-design.md`  
**Estado:** Listo para ejecutar  

---

## Resumen

Migrar la estructura por capas horizontales (`domain/`, `application/`, `infrastructure/`, `presentation/`) a vertical slices por feature (`calculos/`, `equipos/`) con un `shared/kernel/` para value objects compartidos.

**Principio de migración:** De adentro hacia afuera. Cada paso termina con `go test ./...` verde antes de avanzar.

**Agentes responsables por paso:**

| Pasos | Agente |
|-------|--------|
| 1–2 | `domain-agent` |
| 3 | `application-agent` |
| 4 | `infrastructure-agent` |
| 5–8 | Coordinador |

---

## Paso 1 — Crear `shared/kernel/` y mover value objects

**Agente:** `domain-agent`  
**Scope:** `internal/shared/kernel/valueobject/`

### Tareas

**1.1** Crear carpeta `internal/shared/kernel/valueobject/`

**1.2** Mover los siguientes archivos a `internal/shared/kernel/valueobject/`:
- `internal/domain/valueobject/corriente.go` → `internal/shared/kernel/valueobject/corriente.go`
- `internal/domain/valueobject/corriente_test.go` → `internal/shared/kernel/valueobject/corriente_test.go`
- `internal/domain/valueobject/tension.go` → `internal/shared/kernel/valueobject/tension.go`
- `internal/domain/valueobject/tension_test.go` → `internal/shared/kernel/valueobject/tension_test.go`
- `internal/domain/valueobject/temperatura.go` → `internal/shared/kernel/valueobject/temperatura.go`
- `internal/domain/valueobject/temperatura_test.go` → `internal/shared/kernel/valueobject/temperatura_test.go`
- `internal/domain/valueobject/material_conductor.go` → `internal/shared/kernel/valueobject/material_conductor.go`
- `internal/domain/valueobject/material_conductor_test.go` → `internal/shared/kernel/valueobject/material_conductor_test.go`
- `internal/domain/valueobject/conductor.go` → `internal/shared/kernel/valueobject/conductor.go`
- `internal/domain/valueobject/conductor_test.go` → `internal/shared/kernel/valueobject/conductor_test.go`
- `internal/domain/valueobject/resistencia_reactancia.go` → `internal/shared/kernel/valueobject/resistencia_reactancia.go`
- `internal/domain/valueobject/charola.go` → `internal/shared/kernel/valueobject/charola.go`
- `internal/domain/valueobject/charola_test.go` → `internal/shared/kernel/valueobject/charola_test.go`
- `internal/domain/valueobject/tabla_entrada.go` → `internal/shared/kernel/valueobject/tabla_entrada.go`

**1.3** Actualizar `package` declaration en cada archivo movido:
- De: `package valueobject`
- A: `package valueobject` ← mismo nombre, solo cambia el import path

**1.4** Crear `internal/shared/kernel/AGENTS.md` con reglas del kernel

**1.5** Verificar:
```bash
go test ./internal/shared/...
```

---

## Paso 2 — Crear `calculos/domain/` y mover entities + services

**Agente:** `domain-agent`  
**Scope:** `internal/calculos/domain/`

### Tareas

**2.1** Crear carpetas:
- `internal/calculos/domain/entity/`
- `internal/calculos/domain/service/`

**2.2** Mover entities a `internal/calculos/domain/entity/`:
- Todos los archivos de `internal/domain/entity/` → `internal/calculos/domain/entity/`
- Actualizar imports de `valueobject` al nuevo path: `github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject`

**2.3** Mover services a `internal/calculos/domain/service/`:
- Todos los archivos de `internal/domain/service/` → `internal/calculos/domain/service/`
- Actualizar imports de `valueobject` y `entity` a los nuevos paths

**2.4** Crear `internal/calculos/domain/AGENTS.md`

**2.5** Verificar:
```bash
go test ./internal/calculos/domain/...
go test ./internal/shared/...
```

---

## Paso 3 — Crear `calculos/application/` y mover ports + use cases + DTOs

**Agente:** `application-agent`  
**Scope:** `internal/calculos/application/`

### Tareas

**3.1** Crear carpetas:
- `internal/calculos/application/port/`
- `internal/calculos/application/usecase/`
- `internal/calculos/application/usecase/helpers/`
- `internal/calculos/application/dto/`

**3.2** Mover ports a `internal/calculos/application/port/`:
- `internal/application/port/tabla_nom_repository.go` → `internal/calculos/application/port/`
- `internal/application/port/equipo_repository.go` → `internal/calculos/application/port/`
- `internal/application/port/seleccionar_temperatura.go` → `internal/calculos/application/port/`
- Actualizar imports a nuevos paths de `domain/` y `shared/kernel/`

**3.3** Mover DTOs a `internal/calculos/application/dto/`:
- `internal/application/dto/equipo_input.go` → `internal/calculos/application/dto/`
- `internal/application/dto/equipo_input_test.go` → `internal/calculos/application/dto/`
- `internal/application/dto/memoria_output.go` → `internal/calculos/application/dto/`
- `internal/application/dto/errors.go` → `internal/calculos/application/dto/`
- Actualizar imports a nuevos paths

**3.4** Mover use cases a `internal/calculos/application/usecase/`:
- Todos los archivos de `internal/application/usecase/` → `internal/calculos/application/usecase/`
- `internal/application/usecase/helpers/nombre_tabla.go` → `internal/calculos/application/usecase/helpers/`
- Actualizar imports a nuevos paths de `domain/`, `port/`, `dto/` y `shared/kernel/`

**3.5** Crear `internal/calculos/application/AGENTS.md`

**3.6** Verificar:
```bash
go test ./internal/calculos/application/...
go test ./internal/calculos/domain/...
go test ./internal/shared/...
```

---

## Paso 4 — Crear `calculos/infrastructure/` y mover adapters

**Agente:** `infrastructure-agent`  
**Scope:** `internal/calculos/infrastructure/`

### Tareas

**4.1** Crear carpetas:
- `internal/calculos/infrastructure/adapter/driver/http/`
- `internal/calculos/infrastructure/adapter/driver/http/formatters/`
- `internal/calculos/infrastructure/adapter/driver/http/middleware/`
- `internal/calculos/infrastructure/adapter/driven/csv/`
- `internal/calculos/infrastructure/adapter/driven/csv/testdata/`
- `internal/calculos/infrastructure/adapter/driven/postgres/`

**4.2** Mover CSV repository a `internal/calculos/infrastructure/adapter/driven/csv/`:
- `internal/infrastructure/repository/csv_tabla_nom_repository.go` → `internal/calculos/infrastructure/adapter/driven/csv/`
- `internal/infrastructure/repository/csv_tabla_nom_repository_test.go` → `internal/calculos/infrastructure/adapter/driven/csv/`
- `internal/infrastructure/repository/seleccionar_temperatura.go` → `internal/calculos/infrastructure/adapter/driven/csv/`
- `internal/infrastructure/repository/testdata/` → `internal/calculos/infrastructure/adapter/driven/csv/testdata/`
- Actualizar imports a nuevos paths de `port/` y `shared/kernel/`

**4.3** Mover HTTP handler a `internal/calculos/infrastructure/adapter/driver/http/`:
- `internal/presentation/handler/calculo_handler.go` → `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`
- `internal/presentation/handler/calculo_handler_test.go` → `internal/calculos/infrastructure/adapter/driver/http/calculo_handler_test.go`
- Actualizar imports a nuevos paths

**4.4** Mover formatters a `internal/calculos/infrastructure/adapter/driver/http/formatters/`:
- `internal/presentation/formatters/nombre_tabla.go` → `internal/calculos/infrastructure/adapter/driver/http/formatters/`
- `internal/presentation/formatters/nombre_tabla_test.go` → `internal/calculos/infrastructure/adapter/driver/http/formatters/`
- `internal/presentation/formatters/observaciones.go` → `internal/calculos/infrastructure/adapter/driver/http/formatters/`
- `internal/presentation/formatters/observaciones_test.go` → `internal/calculos/infrastructure/adapter/driver/http/formatters/`

**4.5** Mover middleware:
- `internal/presentation/middleware/.gitkeep` → `internal/calculos/infrastructure/adapter/driver/http/middleware/.gitkeep`

**4.6** Mover router:
- `internal/presentation/router.go` → `internal/calculos/infrastructure/router.go`
- Actualizar imports al nuevo path del handler

**4.7** Crear `internal/calculos/infrastructure/AGENTS.md`

**4.8** Verificar:
```bash
go test ./internal/calculos/infrastructure/...
go test ./internal/calculos/...
```

---

## Paso 5 — Actualizar `cmd/api/main.go`

**Agente:** Coordinador  
**Scope:** `cmd/api/main.go`

### Tareas

**5.1** Leer el `main.go` actual y mapear todos los imports

**5.2** Actualizar todos los imports a los nuevos paths:
- `internal/domain/...` → `internal/calculos/domain/...`
- `internal/application/...` → `internal/calculos/application/...`
- `internal/infrastructure/...` → `internal/calculos/infrastructure/...`
- `internal/presentation/...` → `internal/calculos/infrastructure/...`
- `internal/domain/valueobject/...` → `internal/shared/kernel/valueobject/...`

**5.3** Verificar compilación y servidor:
```bash
go build ./...
go test ./...
```

---

## Paso 6 — Crear `equipos/` placeholder vacío

**Agente:** Coordinador  
**Scope:** `internal/equipos/`

### Tareas

**6.1** Crear estructura de carpetas con `.gitkeep`:
```
internal/equipos/
  domain/.gitkeep
  application/port/.gitkeep
  application/dto/.gitkeep
  infrastructure/adapter/driven/postgres/.gitkeep
```

**6.2** Crear `internal/equipos/AGENTS.md` con scope futuro y reglas de la feature

**6.3** Verificar que no rompe nada:
```bash
go build ./...
```

---

## Paso 7 — Eliminar carpetas viejas

**Agente:** Coordinador

### Tareas

**7.1** Eliminar SOLO cuando `go test ./...` esté completamente verde:
- `internal/domain/`
- `internal/application/`
- `internal/infrastructure/`
- `internal/presentation/`

**7.2** Verificar final:
```bash
go test ./...
go build ./...
go vet ./...
```

---

## Paso 8 — Actualizar AGENTS.md

**Agente:** Coordinador

### Tareas

**8.1** Actualizar `AGENTS.md` raíz:
- Nueva tabla de estructura de carpetas
- Nueva sección "Sistema de Agentes": cuándo invocar cada uno, orden de delegación
- Actualizar tabla de auto-invocación con nuevos paths
- Actualizar tabla de "Guias por Capa" con nuevas rutas

**8.2** Crear `internal/shared/kernel/AGENTS.md`

**8.3** Crear `internal/calculos/AGENTS.md` (overview de la feature)

**8.4** Actualizar `internal/calculos/domain/AGENTS.md` (mover desde `internal/domain/AGENTS.md`, actualizar paths)

**8.5** Actualizar `internal/calculos/application/AGENTS.md` (mover desde `internal/application/AGENTS.md`, actualizar paths)

**8.6** Crear `internal/calculos/infrastructure/AGENTS.md`

**8.7** Crear `internal/equipos/AGENTS.md`

**8.8** Verificar que todos los links internos en AGENTS.md apuntan a paths correctos

---

## Criterios de Éxito Finales

```bash
go test ./...        # todos verdes
go build ./...       # sin errores
go vet ./...         # sin warnings
```

- [ ] Ningún archivo importa desde `internal/domain/`, `internal/application/`, `internal/infrastructure/`, `internal/presentation/`
- [ ] Ningún cross-import entre `calculos/` y `equipos/`
- [ ] `shared/kernel/` no importa ninguna feature
- [ ] `cmd/api/main.go` es el único que conoce múltiples features
- [ ] Todos los AGENTS.md actualizados
