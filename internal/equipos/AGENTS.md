# Feature: equipos

Catálogo de equipos de filtros eléctricos Garfex (filtros activos y de rechazo) con CRUD completo persistido en PostgreSQL (Supabase via pgx v5).

## Propósito

Gestionar el catálogo de filtros eléctricos que se instalan en instalaciones industriales según normativa NOM. Cada equipo tiene: clave comercial, tipo de filtro, voltaje nominal, amperaje, capacidad ITM, bornes opcionales y tipo de conexión eléctrica opcional.

## Endpoints

| Método | Path | Descripción |
|--------|------|-------------|
| `POST` | `/api/v1/equipos` | Crear nuevo equipo filtro |
| `GET` | `/api/v1/equipos` | Listar con paginación y filtros opcionales |
| `GET` | `/api/v1/equipos/:id` | Obtener equipo por UUID |
| `PUT` | `/api/v1/equipos/:id` | Actualizar equipo completo |
| `DELETE` | `/api/v1/equipos/:id` | Eliminar equipo (idempotente) |

### Query params de listado

| Param | Tipo | Descripción |
|-------|------|-------------|
| `page` | int | Página (default: 1) |
| `page_size` | int | Registros por página (default: 20, max: 100) |
| `tipo` | string | Filtrar por tipo: `A`, `KVA`, `KVAR` |
| `voltaje` | int | Filtrar por voltaje exacto (V) |

### Tipos de filtro (`TipoFiltro`)

| Valor | Descripción |
|-------|-------------|
| `A` | Filtro activo — calificado en Amperes |
| `KVA` | Filtro calificado en KVA |
| `KVAR` | Filtro de rechazo — calificado en KVAR reactivos |

> Los valores del enum coinciden exactamente con el enum PostgreSQL `public.tipo_filtro`.

### Tipos de conexión (`Conexion`) — nullable

| Valor | Descripción |
|-------|-------------|
| `DELTA` | Conexión trifásica en triángulo (∆) |
| `ESTRELLA` | Conexión trifásica en estrella (Y) |
| `MONOFASICO` | Conexión monofásica |
| `BIFASICO` | Conexión bifásica |

> Los valores del enum coinciden exactamente con el enum PostgreSQL `public.conexion`. El campo es opcional (nullable).

## Estructura

```
internal/equipos/
├── domain/
│   └── entity/
│       ├── tipo_filtro.go       ← TipoFiltro enum (A, KVA, KVAR)
│       ├── conexion.go          ← Conexion enum (MONOFASICA, TRIFASICA)
│       ├── equipo_filtro.go     ← EquipoFiltro entity, NewEquipoFiltro()
│       └── errors.go            ← ErrTipoFiltroInvalido, ErrConexionInvalida, ErrVoltajeInvalido, etc.
├── application/
│   ├── port/
│   │   └── equipo_filtro_repository.go  ← Interface: Crear, ObtenerPorID, Listar, Contar, Actualizar, Eliminar
│   ├── dto/
│   │   ├── equipo_filtro_input.go   ← CreateEquipoInput, UpdateEquipoInput, ListEquiposQuery
│   │   ├── equipo_filtro_output.go  ← EquipoOutput, ListEquiposOutput, PaginationMeta
│   │   └── errors.go                ← ErrEquipoNoEncontrado, ErrClaveYaExiste, ErrInputInvalido, ErrIDInvalido
│   └── usecase/
│       ├── crear_equipo.go
│       ├── obtener_equipo.go
│       ├── listar_equipos.go      ← usa Contar() + Listar() para pagination metadata
│       ├── actualizar_equipo.go
│       └── eliminar_equipo.go
└── infrastructure/
    ├── router.go                  ← RegisterEquiposRoutes() monta 5 endpoints bajo /api/v1
    └── adapter/
        ├── driven/postgres/
        │   ├── pool.go                       ← NewPool() pgxpool, max 10 conns, 5s timeout
        │   └── equipo_filtro_repository.go   ← PostgresEquipoFiltroRepository
        └── driver/http/
            └── equipo_handler.go             ← EquipoHandler, 5 métodos, errorResponse con timestamp ISO 8601
```

## Tabla PostgreSQL

```
public.equipos_filtros
├── id           uuid          PK, DEFAULT gen_random_uuid()
├── created_at   timestamptz   DEFAULT now()
├── clave        text          UNIQUE, NULLABLE
├── tipo         tipo_filtro   NOT NULL  ← enum: 'A', 'KVA', 'KVAR'
├── voltaje      integer       NOT NULL
├── "qn/In"      integer       NOT NULL  ← columna con barra diagonal — escapar con comillas en SQL
├── itm          integer       NOT NULL
├── bornes       integer       NULLABLE
├── conexion     conexion      NULLABLE  ← enum: 'DELTA', 'ESTRELLA', 'MONOFASICO', 'BIFASICO'
└── tipo_voltaje tipo_voltaje  NULLABLE  ← enum: 'FF' (fase-fase), 'FN' (fase-neutro). DB default: 'FF'
```

> ⚠️ El campo `qn/In` tiene barra diagonal. En SQL siempre usar comillas dobles: `"qn/In"`.

## Conexión a PostgreSQL

- Supabase self-hosted en servidor Ubuntu, puerto **5434** (mapeado en docker-compose.yml del servidor)
- Variables de entorno: `DB_HOST`, `DB_PORT=5434`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Pool: `internal/equipos/infrastructure/adapter/driven/postgres/pool.go`

## Integración con cálculos (memoria de cálculo)

El catálogo de equipos se consume desde la feature `calculos` mediante el adapter `CalcEquipoFiltroRepository` (`internal/calculos/infrastructure/adapter/driven/postgres/equipo_filtro_repository.go`).

El pool de PostgreSQL se comparte entre ambos repositorios sin duplicar conexiones.

## Mapeo de Errores HTTP

| Error application | HTTP | Causa |
|-------------------|------|-------|
| `ErrInputInvalido` | 400 | Campo requerido faltante o inválido |
| `ErrIDInvalido` | 400 | UUID malformado en path param |
| `ErrEquipoNoEncontrado` | 404 | ID no existe |
| `ErrClaveYaExiste` | 409 | Violación UNIQUE en columna `clave` (SQLSTATE 23505) |
| Error interno | 500 | Error no manejado |

Todas las respuestas de error incluyen campo `timestamp` en ISO 8601.

## Deuda Técnica Conocida

- `PUT` semánticamente debería ser `PATCH` (actualmente acepta campos parciales)
- Sin rate limiting
- Sin autenticación
- Sin OpenAPI/Swagger spec

## Cómo modificar esta feature

Ver guías por capa:
- [domain/AGENTS.md](domain/AGENTS.md)
- [application/AGENTS.md](application/AGENTS.md)
- [infrastructure/AGENTS.md](infrastructure/AGENTS.md)

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
