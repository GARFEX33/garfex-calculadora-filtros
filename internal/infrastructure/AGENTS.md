# Infrastructure Layer

Implementa los ports definidos en `application/port/`.
Tecnologias: PostgreSQL (pgx/v5), CSV (encoding/csv).

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — error handling, interfaces, convenciones de repositorios
> - [`golang-pro`](.agents/skills/golang-pro/SKILL.md) — connection pooling, concurrencia en queries

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Crear o modificar repositorio | `golang-patterns` |
| Configurar pgx pool o conexion BD | `golang-pro` |
| Implementar nuevo CSV reader | `golang-patterns` |

## Estructura

- `repository/` — PostgresEquipoRepository, CSVTablaNOMRepository
- `client/` — PostgresClient (pgx pool)

## PostgresEquipoRepository

- Tabla `equipos_filtros`: clave, tipo, voltaje, "qn/In", itm, bornes
- Mapea "qn/In" segun tipo: ACTIVO → Amperaje, RECHAZO → KVAR
- `context.Context` en todas las queries

## CSVTablaNOMRepository

Lee tablas NOM desde `data/tablas_nom/`.

### Mapeo canalizacion → tabla ampacidad

| TipoCanalizacion | Archivo CSV |
|---|---|
| TUBERIA_PVC / ALUMINIO / ACERO_PG / ACERO_PD | 310-15-b-16.csv |
| CHAROLA_CABLE_ESPACIADO | 310-15-b-17.csv |
| CHAROLA_CABLE_TRIANGULAR | 310-15-b-20.csv |

### Mapeo canalizacion → columna R (Tabla 9)

| TipoCanalizacion | Columna resistencia |
|---|---|
| TUBERIA_PVC / CHAROLA_ESPACIADO / CHAROLA_TRIANGULAR | `res_{material}_pvc` |
| TUBERIA_ALUMINIO | `res_{material}_al` |
| TUBERIA_ACERO_PG / ACERO_PD | `res_{material}_acero` |

Charola no tiene conduit metalico → usa columna PVC (sin efecto de proximidad).

## Variables de Entorno

`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `ENVIRONMENT`
Ver topologia completa en `docs/plans/2026-02-09-arquitectura-inicial-design.md`.

---

## CRITICAL RULES

- ALWAYS: Implementar exactamente el port definido en `application/port/`
- ALWAYS: `context.Context` como primer parametro en todas las operaciones
- ALWAYS: Inyeccion de dependencias via constructor — sin globals ni singletons
- ALWAYS: Usar pgx pool; cerrar rows con `defer rows.Close()`
- ALWAYS: Validar columnas requeridas al cargar CSV — fallar rapido si falta columna
- NEVER: Importar `domain/service` — solo `entity` y `valueobject`
- NEVER: Logica de negocio en repositorios — solo traduccion datos <-> domain
- NEVER: Estado global mutable (`var db *pgxpool.Pool` a nivel de package)
- NEVER: Escribir datos NOM hardcodeados en Go — siempre leer del CSV

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Repositorio PostgreSQL | `PostgresPascalCaseRepository` | `PostgresEquipoRepository` |
| Repositorio CSV | `CSVPascalCaseRepository` | `CSVTablaNOMRepository` |
| Cliente BD | `PascalCaseClient` | `PostgresClient` |
| Archivo | `snake_case.go` | `postgres_equipo_repository.go` |

---

## QA CHECKLIST

- [ ] `go test ./internal/infrastructure/...` pasa
- [ ] Variables de entorno seteadas (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
- [ ] Nuevo repositorio implementa el port completo
- [ ] Sin estado global mutable
- [ ] Rows cerradas con defer en todas las queries
