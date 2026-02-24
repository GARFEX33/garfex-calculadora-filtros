# Equipos — Infrastructure Layer

Implementa los ports definidos en `application/port/`. Tecnologías: PostgreSQL via pgx v5, HTTP via Gin.

## Estructura

```
internal/equipos/infrastructure/
├── router.go                             ← RegisterEquiposRoutes() monta rutas bajo /api/v1
└── adapter/
    ├── driven/postgres/
    │   ├── pool.go                       ← NewPool() pgxpool, max 10 conns, 5s timeout
    │   └── equipo_filtro_repository.go   ← PostgresEquipoFiltroRepository
    └── driver/http/
        └── equipo_handler.go             ← EquipoHandler, 5 métodos HTTP
```

> Las subcarpetas heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Adapters

### Driven — `PostgresEquipoFiltroRepository`

Implementa `port.EquipoFiltroRepository`. Detalles importantes:

- **Pool**: `pgxpool.Pool` inyectado por constructor. Máximo 10 conexiones, timeout 5s.
- **Columna especial**: `"qn/In"` tiene barra diagonal — siempre escapar con comillas dobles en SQL.
- **Unique violation**: detección via `pgconn.PgError` con SQLSTATE `23505` (NO string matching).
- **`buildWhereClause()`**: función interna que evita duplicación entre `Listar` y `Contar`.
- **`Eliminar` es idempotente**: no retorna error si el ID no existe.
- **`conexion` nullable**: mapeada como `*string` en SQL y `*entity.Conexion` en dominio via `mapConexionToDB/mapConexionFromDB`.

```go
// Construcción del repositorio
repo := postgres.NewPostgresEquipoFiltroRepository(pool)
```

### Driver — `EquipoHandler`

Handler HTTP con 5 métodos:

| Método handler | HTTP | Path |
|----------------|------|------|
| `Crear` | `POST` | `/api/v1/equipos` |
| `Listar` | `GET` | `/api/v1/equipos` |
| `ObtenerPorID` | `GET` | `/api/v1/equipos/:id` |
| `Actualizar` | `PUT` | `/api/v1/equipos/:id` |
| `Eliminar` | `DELETE` | `/api/v1/equipos/:id` |

Todas las respuestas de error incluyen campo `timestamp` en ISO 8601:

```json
{
  "error": "equipo no encontrado",
  "timestamp": "2026-02-23T01:00:00Z"
}
```

## Conexión a PostgreSQL

Variables de entorno requeridas:

| Variable | Valor típico | Descripción |
|----------|-------------|-------------|
| `DB_HOST` | `192.168.1.60` | Host Supabase self-hosted |
| `DB_PORT` | `5434` | Puerto mapeado en docker-compose del servidor |
| `DB_USER` | `postgres` | Usuario PostgreSQL |
| `DB_PASSWORD` | — | Contraseña |
| `DB_NAME` | `postgres` | Base de datos |

> ⚠️ El puerto es **5434**, no 5432. El 5432 está ocupado por otro PostgreSQL en el servidor.

## Mapeo de Errores HTTP

| Error application | HTTP | Código |
|-------------------|------|--------|
| `dto.ErrInputInvalido` | 400 | Bad Request |
| `dto.ErrIDInvalido` | 400 | Bad Request |
| `dto.ErrEquipoNoEncontrado` | 404 | Not Found |
| `dto.ErrClaveYaExiste` | 409 | Conflict |
| Error interno | 500 | Internal Server Error |

## Reglas de Oro — Capa Infrastructure

1. **Implementar exactamente el port** — no agregar métodos no definidos en la interfaz
2. **Sin lógica de negocio** — solo traducción entre DB rows y entities
3. **Handlers solo coordinan**: bind → use case → JSON response
4. **Inyección de dependencias** — constructor, no variables globales
5. **Context.Context** — primer parámetro en toda operación I/O

## QA Checklist

- [ ] `go test ./internal/equipos/...` pasa
- [ ] Repositorio implementa el port exactamente (sin métodos extra)
- [ ] Sin estado global
- [ ] Sin lógica de negocio en handler ni repositorio

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
