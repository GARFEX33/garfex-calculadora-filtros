# Design: Agregar unidades de voltaje (V/kV)

## Technical Approach

Seguir el patrón existente del value object `Potencia`:
1. Agregar tipo `UnidadTension` con constantes `V` y `kV`
2. Modificar `NewTension` para aceptar `(valor float64, unidad string)`
3. Normalizar a volts internamente: `kV = V * 1000`
4. Agregar métodos conversores (similar a `KW()`, `KVA()`)
5. Actualizar DTO `EquipoInput` para incluir campo `TensionUnidad`
6. Mantener compatibilidad hacia atrás con default "V"

## Architecture Decisions

### Decision: Usar float64 para valor de tensión

**Choice**: `NewTension(valor float64, unidad string)`
**Alternatives considered**: `NewTension(valor int, unidad string)` - solo acepta enteros
**Rationale**: Necesitamos aceptar valores decimales cuando la unidad es kV (ej: 0.48 kV). Usar float64 permite esto mientras validamos que el valor normalizado a volts sea uno de los valores NOM válidos.

### Decision: Validación NOM después de normalización

**Choice**: Normalizar primero a volts, luego validar contra lista NOM
**Alternatives considered**: Validar antes de normalizar (0.48 kV → 480 V → validar si 480 está en lista)
**Rationale**: Simplifica la lógica - siempre trabajamos con volts para validación. El usuario envía kV, normalizamos a V, validamos contra NOM.

### Decision: Default "V" para compatibilidad hacia atrás

**Choice**: Si `TensionUnidad` está vacío, asumir "V"
**Alternatives considered**: Retornar error si no se especifica unidad
**Rationale**: Mantiene compatibilidad con clientes existentes que no envían `tension_unidad`.

## Data Flow

```
HTTP Request (JSON)
    │
    ▼
EquipoInput DTO ──► ApplyDefaults() ──► ToDomainTension()
    │                              │
    │                     Agregar default "V"
    │                              │
    ▼                              ▼
valueobject.NewTension(valor, unidad)
    │
    ├──► ParseUnidadTension(unidad) ──► UnidadTension{V, kV}
    │
    ├──► normalizarAVolts(valor, unidad) ──► valor en V
    │
    └──► Validar voltaje en map[NOM] ──► Tension o error
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/shared/kernel/valueobject/tension.go` | Modify | Agregar tipo UnidadTension, ParseUnidadTension, modificar NewTension |
| `internal/calculos/application/dto/equipo_input.go` | Modify | Agregar campo TensionUnidad, ApplyDefaults, ToDomainTension |
| `internal/shared/kernel/valueobject/tension_test.go` | Modify | Agregar tests para nuevas funcionalidades |

## Interfaces / Contracts

### Value Object Tension (modificado)

```go
// UnidadTension representa la unidad de voltaje.
type UnidadTension string

const (
    UnidadTensionV  UnidadTension = "V"
    UnidadTensionkV UnidadTension = "kV"
)

// ParseUnidadTension converts a string to UnidadTension.
func ParseUnidadTension(s string) (UnidadTension, error)

// Tension represents an electrical voltage value in Volts. Immutable.
type Tension struct {
    valor  int        // siempre en volts
    unidad UnidadTension
}

// NewTension creates a Tension value object from a value and unit.
// The value is normalized to volts internally.
func NewTension(valor float64, unidad string) (Tension, error)

// NewTensionFromJSON compatible constructor for JSON unmarshaling
// Returns error if the normalized voltage is not a valid NOM value
```

### DTO EquipoInput (modificado)

```go
type EquipoInput struct {
    // ... existing fields ...
    Tension         float64 // Voltaje (ej: 220, 480, o 0.48)
    TensionUnidad  string   // "V" o "kV" (default: "V")
    // ... existing fields ...
}

// ApplyDefaults modificado
func (e *EquipoInput) ApplyDefaults() {
    // ... existing defaults ...
    
    // Default: V para tensión
    if e.TensionUnidad == "" {
        e.TensionUnidad = "V"
    }
}

// ToDomainTension modificado
func (e EquipoInput) ToDomainTension() (valueobject.Tension, error) {
    return valueobject.NewTension(e.Tension, e.TensionUnidad)
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | NewTension con V, NewTension con kV, errores, Edge cases | Tests unitarios en tension_test.go |
| Unit | ApplyDefaults con TensionUnidad vacío | Test en equipo_input_test.go |
| Integration | API endpoint acepta tension_unidad: "kV" | Test de integración con handler |

## Migration / Rollout

No se requiere migración de datos. Es un cambio additive:
- Clientes existentes que no envían `tension_unidad` seguirán funcionando (default: "V")
- Nuevos clientes pueden enviar `tension_unidad: "kV"`

## Open Questions

- [ ] Ninguno - el diseño está completo
