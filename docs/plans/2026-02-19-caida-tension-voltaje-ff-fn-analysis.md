# Análisis: Problema Voltaje Fase-Fase (FF) vs Fase-Neutro (FN) en Caída de Tensión

**Fecha:** 2026-02-19  
**Estado:** En análisis  
**Prioridad:** CRÍTICA

---

## Problema Identificado

La implementación actual de `CalcularCaidaTension` **no distingue entre voltaje fase-fase (Vff) y fase-neutro (Vfn)**, lo que puede causar errores de hasta **73%** en el cálculo del porcentaje de caída de tensión para sistemas trifásicos.

---

## Fórmulas NOM Correctas

Según NOM-001-SEDE-2012, la caída de tensión se calcula como:

```
e = factor × I × Z × L    (caída en volts)
%e = (e / V_referencia) × 100
```

Donde `V_referencia` **depende del sistema eléctrico**:

| Sistema Eléctrico | Factor | V_referencia | Fórmula completa |
|-------------------|--------|--------------|------------------|
| **Monofásico** 1F-2H | 2 | **Vfn** (fase-neutro) | %e = (2·I·Z·L / Vfn) × 100 |
| **Bifásico** 2F-3H | 1 | **Vfn** (fase-neutro) | %e = (I·Z·L / Vfn) × 100 |
| **Trifásico DELTA** 3F-3H | √3 | **Vff** (fase-fase) | %e = (√3·I·Z·L / Vff) × 100 |
| **Trifásico ESTRELLA** 3F-4H | 1 | **Vfn** (fase-neutro) | %e = (I·Z·L / Vfn) × 100 |

### Relaciones de Voltaje

Para sistemas trifásicos balanceados:
```
Vff = √3 × Vfn
Vfn = Vff / √3
```

**Ejemplos comunes:**
- 220V trifásico → Vff=220V, Vfn=127V
- 440V trifásico → Vff=440V, Vfn=254V
- 480V trifásico → Vff=480V, Vfn=277V

---

## Problema en Código Actual

### Código actual (INCORRECTO):
```go
// Línea 101 en calculo_caida_tension.go
porcentaje := (caida / float64(tension.Valor())) * 100
```

**Asume que `tension` es siempre el voltaje correcto**, pero:
- ¿Qué ingresa el usuario? ¿Vff o Vfn?
- ¿Cómo sabe el sistema cuál usar?

---

## Casos Problemáticos

### Caso 1: Usuario ingresa **Vff** pensando en sistema trifásico

**Request:**
```json
{
  "tension": 220,        ← Usuario piensa "220V trifásico"
  "sistema_electrico": "ESTRELLA"
}
```

**¿Qué debería pasar?**
- Sistema ESTRELLA requiere Vfn
- Si 220V es Vff → Vfn = 220/√3 = 127V
- Fórmula: %e = (I·Z·L / 127) × 100

**¿Qué pasa actualmente?**
- El código usa: %e = (I·Z·L / 220) × 100 ❌
- **Error: resultado 73% menor** (127/220 = 0.577)

---

### Caso 2: Usuario ingresa **Vfn**

**Request:**
```json
{
  "tension": 127,        ← Usuario especifica Vfn
  "sistema_electrico": "DELTA"
}
```

**¿Qué debería pasar?**
- Sistema DELTA requiere Vff
- Si 127V es Vfn → Vff = 127×√3 = 220V
- Fórmula: %e = (√3·I·Z·L / 220) × 100

**¿Qué pasa actualmente?**
- El código usa: %e = (√3·I·Z·L / 127) × 100 ❌
- **Error: resultado 73% mayor** (220/127 = 1.732)

---

## Verificación de Tests Actuales

Los tests actuales **NO detectan este problema** porque usan el mismo voltaje arbitrario (220V) para todos los sistemas sin validar si es Vfn o Vff.

### Test actual (NO valida correctamente):
```go
tension, _ := valueobject.NewTension(220)  // ¿Es Vfn o Vff?

// Casos con diferentes sistemas
TestMonofasico:  // usa 220V como Vfn ✅ (puede ser correcto)
TestBifasico:    // usa 220V como Vfn ✅ (puede ser correcto)
TestDelta:       // usa 220V como ¿Vff o Vfn? ⚠️ (ambiguo)
TestEstrella:    // usa 220V como Vfn ✅ (puede ser correcto)
```

**Problema:** Los tests solo verifican **relaciones entre sistemas**, no valores absolutos según normativa NOM.

---

## Soluciones Posibles

### Opción 1: Campo Adicional `tipo_voltaje` (RECOMENDADA)

**API Request:**
```json
{
  "tension": 220,
  "tipo_voltaje": "FASE_FASE",  ← Nuevo campo obligatorio
  "sistema_electrico": "DELTA"
}
```

**Ventajas:**
- ✅ Explícito y sin ambigüedad
- ✅ Usuario controla qué voltaje ingresa
- ✅ Conversión automática interna

**Desventajas:**
- ❌ Breaking change (nuevo campo obligatorio)
- ❌ Mayor complejidad en request

**Lógica interna:**
```go
type TipoVoltaje string
const (
    TipoVoltajeFaseFase   TipoVoltaje = "FASE_FASE"   // Vff
    TipoVoltajeFaseNeutro TipoVoltaje = "FASE_NEUTRO" // Vfn
)

// Convertir a voltaje de referencia según sistema
func ObtenerVoltajeReferencia(
    voltajeIngresado int,
    tipoVoltaje TipoVoltaje,
    sistema SistemaElectrico,
) float64 {
    vIngresado := float64(voltajeIngresado)
    
    // Determinar qué necesita el sistema
    var necesitaVfn bool
    switch sistema {
    case SistemaElectricoMonofasico, SistemaElectricoBifasico, SistemaElectricoEstrella:
        necesitaVfn = true
    case SistemaElectricoDelta:
        necesitaVfn = false  // necesita Vff
    }
    
    // Convertir si es necesario
    if necesitaVfn && tipoVoltaje == TipoVoltajeFaseFase {
        return vIngresado / math.Sqrt(3)  // Vff → Vfn
    }
    if !necesitaVfn && tipoVoltaje == TipoVoltajeFaseNeutro {
        return vIngresado * math.Sqrt(3)  // Vfn → Vff
    }
    
    return vIngresado  // Ya es el tipo correcto
}
```

---

### Opción 2: Asumir siempre **Vfn** y convertir internamente

**API Request:**
```json
{
  "tension": 127,  ← Siempre Vfn (documentado)
  "sistema_electrico": "DELTA"
}
```

**Ventajas:**
- ✅ API simple
- ✅ No breaking change adicional

**Desventajas:**
- ❌ Confuso para usuarios acostumbrados a especificar Vff en trifásicos
- ❌ Requiere conversión mental del usuario (220V → 127V)

---

### Opción 3: Asumir siempre **Vff** y convertir internamente

**API Request:**
```json
{
  "tension": 220,  ← Siempre Vff (documentado)
  "sistema_electrico": "ESTRELLA"
}
```

**Ventajas:**
- ✅ API simple
- ✅ Intuitivo para sistemas trifásicos

**Desventajas:**
- ❌ Confuso para monofásico (¿220V monofásico es Vff=220 o Vfn=220?)
- ❌ Monofásico no tiene "fase-fase", solo fase-neutro

---

### Opción 4: Dos campos separados `vfn` y `vff` (EXCESIVO)

**API Request:**
```json
{
  "vfn": 127,
  "vff": 220,
  "sistema_electrico": "ESTRELLA"
}
```

**Ventajas:**
- ✅ Máxima claridad

**Desventajas:**
- ❌ Redundante (uno se calcula del otro)
- ❌ Riesgo de inconsistencia (usuario ingresa Vfn=127 pero Vff=240)

---

## Problema con Múltiples Hilos por Fase

### Test actual con 2 hilos por fase:

```go
// calculo_caida_tension_nom_test.go línea 108
entrada := EntradaCalculoCaidaTension{
    ResistenciaOhmPorKm: 0.62,
    ReactanciaOhmPorKm:  0.148,
    HilosPorFase:        2,  ← 2 conductores en paralelo
    SistemaElectrico:    SistemaElectricoMonofasico,
}

// ¿Qué pasa internamente?
// Línea 75-76 en calculo_caida_tension.go
rEf := entrada.ResistenciaOhmPorKm / n  // 0.62 / 2 = 0.31 ✅
xEf := entrada.ReactanciaOhmPorKm / n   // 0.148 / 2 = 0.074 ✅
```

**Verificación:**
- ✅ La impedancia se divide correctamente entre hilos en paralelo
- ✅ Z_efectiva = √(0.31² + 0.074²) = 0.319 Ω/km
- ✅ Caída = 2 × 70 × 0.319 × 0.030 = 1.338V
- ✅ Resultado: **mitad de la caída con 1 hilo** ✅

**Conclusión:** El manejo de múltiples hilos **está correcto**.

---

## Verificación de Relaciones Matemáticas

Con voltaje **ambiguo** en tests actuales:

### ¿Son correctas las relaciones?

**Solo SI todos usan el mismo tipo de voltaje (Vfn o Vff):**

```
MONOFASICO = 2 × BIFASICO  ✅ (si ambos usan Vfn)
DELTA = √3 × BIFASICO      ❌ (DELTA usa Vff, BIFASICO usa Vfn)
ESTRELLA = BIFASICO        ✅ (si ambos usan Vfn)
```

**Relación correcta DELTA vs BIFASICO:**

Si ambos usan el **mismo voltaje de línea** (220V):
- BIFASICO: %e = (1 × I × Z × L) / (220/√3) × 100 = (I·Z·L / 127) × 100
- DELTA: %e = (√3 × I × Z × L) / 220 × 100

Relación: DELTA / BIFASICO = (√3·I·Z·L / 220) / (I·Z·L / 127) = (√3 / 220) × 127 = **1.0**

**Si se compara con mismo voltaje de línea, DELTA y BIFASICO dan igual resultado** (esto es INCORRECTO según NOM).

---

## Recomendación

### Implementar **Opción 1** (Campo `tipo_voltaje`)

1. **Agregar campo obligatorio** `tipo_voltaje` al request
2. **Convertir** voltaje internamente según sistema eléctrico
3. **Actualizar tests** con voltajes reales NOM (127V/220V, 277V/480V)
4. **Documentar** breaking change y migración

### Valores NOM comunes para tests:

| Voltaje nominal | Vfn | Vff |
|-----------------|-----|-----|
| 220V (México) | 127V | 220V |
| 440V (industrial) | 254V | 440V |
| 480V (USA) | 277V | 480V |

---

## Próximos Pasos

1. ✅ Documentar problema (este archivo)
2. ⏳ Validar con el usuario si entendió el problema
3. ⏳ Aprobar solución (Opción 1 recomendada)
4. ⏳ Implementar corrección
5. ⏳ Actualizar tests con valores NOM reales
6. ⏳ Actualizar documentación de API

---

## Referencias

- NOM-001-SEDE-2012: Sección 210-19(a) — Caída de tensión permitida
- IEEE Std 141 (Red Book): Voltage Drop Calculations
- Relación Vff = √3 × Vfn para sistemas trifásicos balanceados
