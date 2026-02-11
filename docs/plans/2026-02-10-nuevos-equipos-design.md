# Diseño: Extensión del modelo de equipos (Transformador + Carga)

**Fecha:** 2026-02-10
**Contexto:** Brainstorming sobre extensibilidad del modelo de equipos. Se decidió ampliar de 2 a 4 tipos de equipo en la misma fase.

---

## Decisiones de diseño

1. **`Equipo` base se mantiene** con campos universales: `Clave`, `Tipo`, `Voltaje`, `ITM`
2. **Cada tipo de equipo define su dato de entrada** específico para calcular corriente nominal
3. **Las interfaces `CalculadorCorriente` y `CalculadorPotencia` son el contrato polimórfico** — los servicios (Tasks 9-14) trabajan con estas interfaces, agnósticos al tipo
4. **Solo `Carga` soporta multifase** (1, 2, 3) en esta iteración. Los demás asumen trifásico (√3). Si se necesita multifase en otros equipos, se mueve `Fases` al `Equipo` base (refactor menor)
5. **Renombrar `FiltroActivo.Amperaje` → `AmperajeNominal`** para mayor claridad semántica

---

## Tabla resumen

| Equipo | Dato de entrada | Fórmula In | KVA | KW | KVAR |
|--------|----------------|------------|-----|----|----|
| FiltroActivo | AmperajeNominal | directo | In×V×√3/1000 | =KVA (PF=1) | 0 |
| FiltroRechazo | KVAR | KVAR/(KV×√3) | =KVAR | 0 | dado |
| Transformador | KVA | KVA/(KV×√3) | dado | 0 | 0 |
| Carga | KW + FP + Fases | KW/(KV×factor×FP) | KW/FP | dado | √(KVA²-KW²) |

Factor de fases para Carga:
- 3 fases → √3
- 2 fases → 2
- 1 fase → 1

---

## Cambios al código existente

### 1. `FiltroActivo.Amperaje` → `AmperajeNominal`
- Campo, constructor param, validación, tests

### 2. `tipo_equipo.go` — nuevas constantes
```go
TipoEquipoTransformador TipoEquipo = "TRANSFORMADOR"
TipoEquipoCarga         TipoEquipo = "CARGA"
```
- Actualizar `ParseTipoEquipo` con los nuevos cases

---

## Nuevas entidades

### Transformador
```go
type Transformador struct {
    Equipo
    KVA int
}

func NewTransformador(clave string, voltaje, kva int, itm ITM) (*Transformador, error)

// In = KVA / (KV × √3)
func (t *Transformador) CalcularCorrienteNominal() (valueobject.Corriente, error)

// Potencias: KVA=dado, KW=0, KVAR=0
func (t *Transformador) PotenciaKVA() float64   // float64(KVA)
func (t *Transformador) PotenciaKW() float64    // 0
func (t *Transformador) PotenciaKVAR() float64  // 0
```

### Carga
```go
type Carga struct {
    Equipo
    KW             int
    FactorPotencia float64
    Fases          int // 1, 2, o 3
}

func NewCarga(clave string, voltaje, kw int, fp float64, fases int, itm ITM) (*Carga, error)

// In = KW / (KV × factorFases × FP)
// factorFases: 1→1, 2→2, 3→√3
func (c *Carga) CalcularCorrienteNominal() (valueobject.Corriente, error)

// Potencias: KW=dado, KVA=KW/FP, KVAR=√(KVA²-KW²)
func (c *Carga) PotenciaKVA() float64   // float64(KW) / FP
func (c *Carga) PotenciaKW() float64    // float64(KW)
func (c *Carga) PotenciaKVAR() float64  // √(KVA²-KW²)
```

**Validaciones de `NewCarga`:**
- `kw > 0`
- `0 < fp <= 1`
- `fases ∈ {1, 2, 3}`
- `voltaje > 0` (por división)

---

## Impacto en servicios (Tasks 9-14)

**Ninguno.** Los servicios trabajan con las interfaces `CalculadorCorriente` y `CalculadorPotencia`. Los 4 tipos de equipo las implementan. No hay cambios en el diseño de los servicios.

---

## Impacto en BD (futuro)

- Enum `tipo_equipo` necesitará valores: `FILTRO_ACTIVO`, `FILTRO_RECHAZO`, `TRANSFORMADOR`, `CARGA`
- La tabla `equipos_filtros` podría necesitar renombrarse a `equipos` y agregar columnas para los nuevos tipos
- Esto es infrastructure (Fase 2+), no impacta domain layer
