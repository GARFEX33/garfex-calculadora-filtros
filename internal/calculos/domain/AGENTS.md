# Calculos — Domain Layer

Capa de negocio pura para la feature de cálculos eléctricos. Sin dependencias externas (sin Gin, pgx, CSV).

## Estructura

| Subdirectorio | Contenido |
|---------------|-----------|
| `entity/` | Entidades, tipos, interfaces del dominio de cálculos |
| `service/` | Servicios de cálculo puros (sin I/O) |

## Dependencias permitidas

- `internal/shared/kernel/valueobject` — value objects compartidos
- stdlib de Go

## Dependencias prohibidas

- `internal/application/` o `internal/calculos/application/`
- `internal/infrastructure/` o `internal/calculos/infrastructure/`
- Gin, pgx, encoding/csv, cualquier framework externo

## Guías por Subdirectorio

| Subdirectorio | Ver |
|---------------|-----|
| `entity/` | Entidades: TipoEquipo, TipoCanalizacion, SistemaElectrico, MemoriaCalculo |
| `service/` | 8 servicios de cálculo NOM + IEEE-141 |

## Auto-invocación

| Acción | Referencia |
|--------|-----------|
| Crear/modificar entidad o tipo | `golang-patterns` skill |
| Crear/modificar servicio de cálculo | `golang-patterns` skill |
| Agregar nueva fórmula NOM | Verificar que NO depende de infrastructure |
