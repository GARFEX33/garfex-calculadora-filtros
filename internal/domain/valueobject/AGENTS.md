# Domain — Value Objects

Inmutables. Encapsulan validacion y semantica. Sin dependencias externas.

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — constructores con validación, inmutabilidad, error handling

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Crear nuevo value object | `golang-patterns` |
| Definir constructor con validación | `golang-patterns` |

## Value Objects

### Corriente

- Valida > 0, unidad "A"
- Metodos: `Multiplicar`, `Dividir`
- Constructor: `NewCorriente(valor float64)`

### Tension

- Solo valores NOM permitidos: `127, 220, 240, 277, 440, 480, 600`
- Constructor: `NewTension(valor float64)` — retorna error si valor no esta en lista

### Conductor

- Constructor: `NewConductor(ConductorParams{})`
- Campos requeridos: `Calibre`, `Material`, `SeccionMM2`
- `TipoAislamiento` vacio para conductores desnudos (tierra)

### MaterialConductor

- Tipo: `int` con constantes `MaterialCobre` (0) y `MaterialAluminio` (1)
- Serialización JSON: string "CU" / "AL" (via `MarshalJSON()`)
- Deserialización JSON: case-insensitive, acepta "Cu", "cu", "cobre", "Al", "al", "aluminio"
- Método `String()` retorna: "CU" o "AL"

## Reglas de Value Objects

- Son inmutables — no exponer setters
- Igualdad por valor, no por referencia
- Constructores siempre validan y retornan `(T, error)`

---

## CRITICAL RULES

### Inmutabilidad
- ALWAYS: Campos privados — solo accesibles via métodos getter
- ALWAYS: Constructor `NewXxx(...)` valida y retorna `(T, error)`
- NEVER: Setters o mutacion post-construccion
- NEVER: Zero value valido sin pasar por constructor

### Validacion
- ALWAYS: Validar en el constructor — si llega mal, falla rapido con error descriptivo
- ALWAYS: Valores NOM permitidos hardcodeados en el VO (ej: tensiones validas en `Tension`)
- NEVER: Validacion duplicada fuera del VO — quien lo construye confía en que es valido

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Struct VO | `PascalCase` | `Corriente`, `Tension`, `Conductor` |
| Constructor | `NewPascalCase` | `NewCorriente`, `NewTension` |
| Params struct | `PascalCaseParams` | `ConductorParams` |
| Getter | `PascalCase` (sin Get) | `Valor()`, `Unidad()` |
| Error sentinel | `ErrPascalCase` | `ErrTensionInvalida` |

---

## QA CHECKLIST

- [ ] `go test ./internal/domain/valueobject/...` pasa
- [ ] Constructor retorna error para valores invalidos
- [ ] Test cubre valores en limite (ej: 0, negativos, valores no NOM)
- [ ] Sin setters expuestos
- [ ] Sin imports externos al stdlib de Go
