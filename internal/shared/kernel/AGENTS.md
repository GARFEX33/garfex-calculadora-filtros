# Shared Kernel

Value objects y utilidades compartidas entre todas las features.

## Propósito

El kernel contiene conceptos que son **transversales** al dominio:

- Corriente, Potencia, Tensión, Temperatura (conceptos eléctricos básicos)
- MaterialConductor, Conductor
- Charola, ResistenciaReactancia

## Regla de Oro

**El kernel NO importa ninguna feature.**

```go
// ESTO ESTÁ PROHIBIDO:
import "github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"

// ESTO ESTÁ PERMITIDO:
import "github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
```

## Cómo agregar al kernel

### Criterio

- Si el VO es usado por múltiples features → `shared/kernel/`
- Si es específico de una feature → `internal/{feature}/domain/`

## Estructura

```
internal/shared/kernel/
└── valueobject/
    ├── corriente.go
    ├── potencia.go
    ├── tension.go
    ├── temperatura.go
    ├── material_conductor.go
    ├── conductor.go
    ├── charola.go
    ├── resistencia_reactancia.go
    └── tabla_entrada.go
```

## Características de un buen VO

```go
// Inmutable
type Corriente float64

// Validación en constructor
func NewCorriente(valor float64) (Corriente, error) {
    if valor <= 0 {
        return 0, ErrCorrienteInvalida
    }
    return Corriente(valor), nil
}

// Métodos de comportamiento
func (c Corriente) Amperes() float64 { return float64(c) }
```

## Value Objects con Unidades — Patrón

Cuando un VO acepta múltiples unidades, seguir el patrón de `Potencia` y `Tension`:

1. Definir tipo `Unidad*` como `string` con constantes
2. Implementar `ParseUnidad*(s string) (Unidad*, error)` con variantes case-insensitive
3. El constructor acepta `(valor float64, unidad string)` y normaliza internamente
4. Guardar la unidad original para poder reportarla con `Unidad()`
5. Exponer métodos de conversión (`EnKilovoltios()`, `KW()`, etc.)

**Value objects con unidades actualmente:**

| VO        | Unidades           | Normalización interna | Default |
| --------- | ------------------ | --------------------- | ------- |
| `Potencia` | W, KW, KVA, KVAR  | Watts                 | —       |
| `Tension`  | V, kV              | Volts                 | V       |

**Constructor de `Tension` (firma actual):**
```go
// NewTension acepta V o kV y normaliza a volts internamente.
// Valida que el resultado sea un voltaje NOM válido (127,220,240,277,440,480,600).
func NewTension(valor float64, unidad string) (Tension, error)
```

## Reglas Críticas — Shared Kernel

*Estas reglas son específicas para el Shared Kernel (value objects compartidos). Ver [docs/reference/structure.md](../../../docs/reference/structure.md) para reglas globales.*

- **NEVER**: Importar nada de `internal/calculos/`, `internal/equipos/` u otras features
- **NEVER**: Dependencias externas (sin Gin, sin pgx, sin CSV)
- **ALWAYS**: Value objects inmutables con constructores que validan
- **ALWAYS**: Constructores retornan `(T, error)`

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
- Ubicación: `internal/shared/kernel/valueobject/`
