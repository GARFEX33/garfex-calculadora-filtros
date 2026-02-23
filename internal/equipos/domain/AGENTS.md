# Equipos — Domain Layer

Capa de negocio pura para el catálogo de filtros eléctricos. Sin dependencias externas (sin Gin, pgx, CSV).

## Estructura

| Subdirectorio | Contenido |
|---------------|-----------|
| `entity/` | Entidades y tipos del dominio de equipos filtros |

> No existe `service/` — la lógica de negocio de equipos es de validación en constructores, no de cálculo.

## Entidades

### `EquipoFiltro`

Entidad principal del catálogo. Campos:

| Campo | Tipo Go | Descripción |
|-------|---------|-------------|
| `ID` | `uuid.UUID` | Generado por PostgreSQL en INSERT |
| `CreatedAt` | `time.Time` | Generado por PostgreSQL en INSERT |
| `Clave` | `*string` | Clave comercial (nullable, UNIQUE en DB) |
| `Tipo` | `TipoFiltro` | Enum: A, KVA, KVAR |
| `Voltaje` | `int` | Voltaje nominal en Volts |
| `Amperaje` | `int` | Corriente nominal Qn/In en Amperes |
| `ITM` | `int` | Capacidad del interruptor termomagnético en Amperes |
| `Bornes` | `*int` | Número de bornes (nullable) |
| `Conexion` | `*Conexion` | Tipo de conexión eléctrica (nullable): MONOFASICA, TRIFASICA |

Constructor: `NewEquipoFiltro(clave, tipo, voltaje, amperaje, itm, bornes, conexion)` — valida que voltaje, amperaje e ITM sean > 0. `Bornes` y `Conexion` son nullable.

### `TipoFiltro`

Enum que mapea exactamente al enum PostgreSQL `public.tipo_filtro`:

| Constante | Valor DB | Descripción |
|-----------|---------|-------------|
| `TipoFiltroA` | `"A"` | Filtro activo en Amperes |
| `TipoFiltroKVA` | `"KVA"` | Filtro calificado en KVA |
| `TipoFiltroKVAR` | `"KVAR"` | Filtro de rechazo reactivo |

`ParseTipoFiltro(s string)` convierte string → TipoFiltro con error claro si el valor no es válido.

### `Conexion`

Enum que mapea exactamente al enum PostgreSQL `public.conexion`:

| Constante | Valor DB | Descripción |
|-----------|---------|-------------|
| `ConexionDelta` | `"DELTA"` | Conexión trifásica en triángulo (∆) |
| `ConexionEstrella` | `"ESTRELLA"` | Conexión trifásica en estrella (Y) |
| `ConexionMonofasico` | `"MONOFASICO"` | Conexión monofásica |
| `ConexionBifasico` | `"BIFASICO"` | Conexión bifásica |

`ParseConexion(s string)` convierte string → Conexion con error claro si el valor no es válido.

## Dependencias permitidas

- `github.com/google/uuid`
- stdlib de Go

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../../docs/reference/structure.md)

- Sin `internal/shared/kernel` — los equipos no usan value objects eléctricos del kernel
- Sin `internal/calculos/` ni ninguna otra feature
- Sin Gin, pgx, CSV

## Reglas de Oro — Capa Domain

1. Domain nunca depende de Application ni Infrastructure
2. Sin I/O (no leer archivos, no HTTP, no DB)
3. Puro Go + lógica de validación

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../../docs/reference/structure.md)
