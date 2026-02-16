# Application Ports

Interfaces que definen contratos con la infraestructura.

## Puertos

| Puerto | Tipo | Descripción |
|--------|------|-------------|
| `TablaNOMRepository` | Driven | Acceso a tablas NOM (CSV/DB) |
| `EquipoRepository` | Driven | Acceso a catálogo de equipos |
| `SeleccionarTemperatura` | Driven | Selección de temperatura por estado |

## Reglas

- Solo interfaces (sin implementación)
- Definidas en `application/port/`
- Implementadas en `infrastructure/`

## Ejemplo

```go
type TablaNOMRepository interface {
    ObtenerTablaAmpacidad(ctx context.Context, ...) ([]valueobject.EntradaTablaConductor, error)
}
```
