# Equipos — Application Layer

Orquesta el dominio de equipos filtros. Define contratos (ports), no implementaciones.

## Estructura

```
internal/equipos/application/
├── port/
│   └── equipo_filtro_repository.go  ← Interface EquipoFiltroRepository
├── dto/
│   ├── equipo_filtro_input.go   ← CreateEquipoInput, UpdateEquipoInput, ListEquiposQuery
│   ├── equipo_filtro_output.go  ← EquipoOutput, ListEquiposOutput, PaginationMeta
│   └── errors.go                ← Errores sentinel de application
└── usecase/
    ├── crear_equipo.go
    ├── obtener_equipo.go
    ├── listar_equipos.go
    ├── actualizar_equipo.go
    └── eliminar_equipo.go
```

> Las subcarpetas `port/`, `usecase/` y `dto/` heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Port: `EquipoFiltroRepository`

Interface que infrastructure implementa:

```go
type EquipoFiltroRepository interface {
    Crear(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error)
    ObtenerPorID(ctx context.Context, id uuid.UUID) (*entity.EquipoFiltro, error)
    Listar(ctx context.Context, filtros FiltrosListado) ([]*entity.EquipoFiltro, error)
    Contar(ctx context.Context, filtros FiltrosListado) (int, error)
    Actualizar(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error)
    Eliminar(ctx context.Context, id uuid.UUID) error
}
```

`FiltrosListado` contiene `Tipo *entity.TipoFiltro`, `Voltaje *int`, `Limit int`, `Offset int`.

## DTOs

### Input

| DTO | Uso |
|-----|-----|
| `CreateEquipoInput` | Body de POST — tiene `Validate()` y `ToDomain()` |
| `UpdateEquipoInput` | Body de PUT — misma validación, ID viene del path |
| `ListEquiposQuery` | Query params de GET — tiene `ApplyDefaults()` y `Offset()` |

Campos `Clave`, `Bornes` y `Conexion` son `*string` / `*int` (punteros para nullability). `Conexion` acepta `"MONOFASICA"` o `"TRIFASICA"` cuando se provee.

### Output

| DTO | Uso |
|-----|-----|
| `EquipoOutput` | Respuesta de un solo equipo — todos primitivos, `created_at` en ISO 8601 |
| `ListEquiposOutput` | Respuesta paginada con `equipos []EquipoOutput` y `pagination PaginationMeta` |
| `PaginationMeta` | `page`, `page_size`, `total`, `total_pages`, `has_next`, `has_prev` |

Constructores de mapeo: `FromDomain(e)` y `FromDomainList(entities, page, pageSize, total)`.

### Errores sentinel

| Error | Descripción |
|-------|-------------|
| `ErrEquipoNoEncontrado` | UUID no existe en DB |
| `ErrClaveYaExiste` | Violación UNIQUE en `clave` (SQLSTATE 23505) |
| `ErrInputInvalido` | Campo requerido faltante o valor inválido |
| `ErrIDInvalido` | UUID malformado en path param |

## Use Cases

| Use Case | Responsabilidad única |
|----------|-----------------------|
| `CrearEquipoUseCase` | Validar input → ToDomain → repo.Crear → FromDomain |
| `ObtenerEquipoUseCase` | Parse UUID → repo.ObtenerPorID → FromDomain |
| `ListarEquiposUseCase` | ApplyDefaults → repo.Contar + repo.Listar → FromDomainList |
| `ActualizarEquipoUseCase` | Parse UUID → Validate → repo.Actualizar → FromDomain |
| `EliminarEquipoUseCase` | Parse UUID → repo.Eliminar |

`ListarEquiposUseCase` llama tanto a `Contar()` como a `Listar()` para construir la `PaginationMeta` completa.

## Dependencias permitidas

- `internal/equipos/domain/entity`
- `github.com/google/uuid`
- stdlib de Go

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../../docs/reference/structure.md)

- Sin `internal/shared/kernel` (equipos no usan value objects eléctricos)
- Sin Gin, pgx
- Sin `internal/calculos/`

## Reglas de Oro — Capa Application

1. **Un use case = una responsabilidad** — sin lógica de negocio, solo orquestación
2. **DTOs con primitivos** — nunca exponer `entity.TipoFiltro` ni `uuid.UUID` fuera de domain
3. **Context como primer parámetro** en toda operación con I/O
4. **Ports driven** — application define la interfaz, infrastructure la implementa

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../../docs/reference/structure.md)
