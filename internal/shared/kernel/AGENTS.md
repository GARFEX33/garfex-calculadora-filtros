# Shared Kernel

Value objects compartidos entre múltiples features del sistema.

## Reglas Críticas

- NEVER: Importar nada de `internal/calculos/`, `internal/equipos/` u otras features
- NEVER: Dependencias externas (sin Gin, sin pgx, sin encoding/csv)
- ALWAYS: Value objects inmutables con constructores que validan
- ALWAYS: Constructores retornan `(T, error)`

## Contenido

### `valueobject/`

Value objects puros del dominio eléctrico NOM:

| Tipo | Descripción |
|------|-------------|
| `Corriente` | Corriente eléctrica en Amperes. Valida > 0. |
| `Tension` | Tensión eléctrica en Volts. Solo valores NOM permitidos. |
| `Temperatura` | Rating de temperatura (60, 75, 90°C). |
| `MaterialConductor` | Material del conductor (Cu=0, Al=1). |
| `Conductor` | Conductor con propiedades físicas/eléctricas completas. |
| `ConductorParams` | Params para construir un Conductor. |
| `ResistenciaReactancia` | Valores R y X para cálculo de caída de tensión. |
| `CableControl` | Cable de control para charola. |
| `ConductorCharola` | Conductor con diámetro para cálculo de charola. |
| `EntradaTablaConductor` | Fila de tabla NOM 310-15(b)(16). |
| `EntradaTablaTierra` | Fila de tabla NOM 250-122. |
| `EntradaTablaCanalizacion` | Fila de tabla de canalizaciones. |

## Auto-invocación

| Acción | Referencia |
|--------|-----------|
| Crear/modificar value object | `golang-patterns` skill |
| Agregar nuevo tipo compartido | Verificar que NO depende de ninguna feature |
