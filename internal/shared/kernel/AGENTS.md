# Shared Kernel

Value objects y utilidades compartidas entre todas las features.

## Propósito

El kernel contiene conceptos que son **transversales** al dominio:

- Corriente, Tensión, Temperatura (conceptos eléctricos básicos)
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

### Opción A: Durante nueva feature

Cuando el `domain-agent` detecta que un VO es cross-feature, lo crea aquí:

```
# domain-agent automáticamente:
# Si el VO es usado por múltiples features → shared/kernel/
# Si es específico de una feature → internal/{feature}/domain/
```

### Opción B: Agregar VO nuevo

```bash
# Orquestador:
# "domain-agent: agregar VO PotenciaActiva al kernel"
# "Se usará en calculos y en equipos"
```

## Estructura

```
internal/shared/kernel/
└── valueobject/
    ├── corriente.go
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

## Reglas Críticas

- **NEVER**: Importar nada de `internal/calculos/`, `internal/equipos/` u otras features
- **NEVER**: Dependencias externas (sin Gin, sin pgx, sin CSV)
- **ALWAYS**: Value objects inmutables con constructores que validan
- **ALWAYS**: Constructores retornan `(T, error)`

## Referencias

- Agente responsable: `domain-agent`
- Comando: `orchestrate-agents`
- Ubicación: `internal/shared/kernel/valueobject/`
