# Pruebas Manuales - Endpoint Caída de Tensión

## Endpoint

```
POST /api/v1/calculos/caida-tension
```

## Fórmula NOM Implementada

El endpoint utiliza la **fórmula NOM simplificada** para calcular la caída de tensión:

```
e = factor × I × Z × L
%e = (e / V_referencia) × 100

Donde:
- Z = √(R² + X²)  [impedancia del conductor]
- R = resistencia Ω/km (Tabla 9 NOM)
- X = reactancia Ω/km (Tabla 9 NOM)
- I = corriente ajustada (A)
- L = longitud del circuito (km)
- V_referencia = voltaje de referencia según sistema eléctrico
```

**Factores por sistema eléctrico:**
- `MONOFASICO` (1F-2H): **2** → usa **Vfn** (fase-neutro)
- `BIFASICO` (2F-3H): **1** → usa **Vfn** (fase-neutro)
- `DELTA` (3F-3H): **√3 ≈ 1.732** → usa **Vff** (fase-fase)
- `ESTRELLA` (3F-4H): **1** → usa **Vfn** (fase-neutro)

**Relación de voltajes:**
```
Vff = √3 × Vfn
Vfn = Vff / √3

Ejemplos comunes (México):
- 127V = Vfn (fase-neutro)
- 220V = Vff (fase-fase)
- 220V / √3 = 127V
```

---

## ⚠️ Campo `tipo_voltaje` (OBLIGATORIO)

A partir de la versión actual, el endpoint requiere especificar el **tipo de voltaje** ingresado:

| Campo `tipo_voltaje` | Descripción | Ejemplo |
|----------------------|-------------|---------|
| `"FASE_NEUTRO"` o `"FN"` | Voltaje entre fase y neutro | 127V, 277V |
| `"FASE_FASE"` o `"FF"` | Voltaje entre fases (línea a línea) | 220V, 480V |

**El sistema convierte automáticamente** al voltaje de referencia correcto según el sistema eléctrico.

---

## Casos de Prueba con Voltajes NOM Reales

### Caso 1: MONOFASICO con Vfn 127V

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 2.11,
    "CaidaVolts": 2.68,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- Tabla 9 NOM para 2 AWG Cu PVC: R=0.62 Ω/km, X=0.148 Ω/km
- Z = √(0.62² + 0.148²) = 0.6374 Ω/km
- e = 2 × 70 × 0.6374 × 0.030 = 2.677 V
- V_referencia = 127V (ya es Vfn, no se convierte)
- %e = (2.677 / 127) × 100 = **2.11%**

**Comando curl:**
```bash
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{
    "calibre": "2 AWG",
    "material": "Cu",
    "tipo_canalizacion": "TUBERIA_PVC",
    "corriente_ajustada": 70,
    "longitud_circuito": 30,
    "tension": 127,
    "tipo_voltaje": "FASE_NEUTRO",
    "sistema_electrico": "MONOFASICO",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

### Caso 2: MONOFASICO con Vff 220V (convertido a Vfn)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
  "tipo_voltaje": "FASE_FASE",
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 2.11,
    "CaidaVolts": 2.68,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- Usuario ingresa 220V como Vff
- Sistema convierte: V_referencia = 220 / √3 = 127V (Vfn)
- e = 2 × 70 × 0.6374 × 0.030 = 2.677 V
- %e = (2.677 / 127) × 100 = **2.11%**
- **Mismo resultado que Caso 1** ✅

---

### Caso 3: BIFASICO con Vfn 127V

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "BIFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 1.05,
    "CaidaVolts": 1.34,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
- V_referencia = 127V (Vfn)
- %e = (1.338 / 127) × 100 = **1.05%**
- **Exactamente la mitad que MONOFASICO** (factor 1 vs factor 2) ✅

---

### Caso 4: DELTA con Vff 220V

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
  "tipo_voltaje": "FASE_FASE",
  "sistema_electrico": "DELTA",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 1.05,
    "CaidaVolts": 2.32,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- e = √3 × 70 × 0.6374 × 0.030 = 2.318 V
- V_referencia = 220V (ya es Vff, no se convierte)
- %e = (2.318 / 220) × 100 = **1.05%**
- **Mismo porcentaje que BIFASICO** ✅

---

### Caso 5: DELTA con Vfn 127V (convertido a Vff)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "DELTA",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 1.05,
    "CaidaVolts": 2.32,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- Usuario ingresa 127V como Vfn
- Sistema convierte: V_referencia = 127 × √3 = 220V (Vff)
- e = √3 × 70 × 0.6374 × 0.030 = 2.318 V
- %e = (2.318 / 220) × 100 = **1.05%**
- **Mismo resultado que Caso 4** ✅

---

### Caso 6: ESTRELLA con Vfn 127V

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "ESTRELLA",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 1.05,
    "CaidaVolts": 1.34,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
- V_referencia = 127V (Vfn)
- %e = (1.338 / 127) × 100 = **1.05%**
- **Mismo resultado que BIFASICO** ✅

---

### Caso 7: 2 Hilos por Fase (impedancia reducida)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 2,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "Porcentaje": 1.05,
    "CaidaVolts": 1.34,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.319
  }
}
```

**Verificación:**
- R_ef = 0.62 / 2 = 0.31 Ω/km
- X_ef = 0.148 / 2 = 0.074 Ω/km
- Z_ef = √(0.31² + 0.074²) = 0.319 Ω/km
- e = 2 × 70 × 0.319 × 0.030 = 1.338 V
- %e = (1.338 / 127) × 100 = **1.05%**
- **Exactamente la mitad que con 1 hilo** ✅

---

### Caso 8: Material Aluminio (Al)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Al",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Verificación:**
- Tabla 9 NOM para 2 AWG Al PVC: R=1.02 Ω/km, X=0.144 Ω/km
- Z = √(1.02² + 0.144²) = 1.030 Ω/km
- e = 2 × 70 × 1.030 × 0.030 = 4.326 V
- %e = (4.326 / 127) × 100 = **3.41%**
- `Cumple` = **false** (excede límite 3%)
- **Mayor caída que cobre** (Al tiene mayor resistencia)

---

### Caso 9: Sistema 480V (USA industrial)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "tipo_voltaje": "FASE_FASE",
  "sistema_electrico": "ESTRELLA",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Verificación:**
- Usuario ingresa 480V como Vff
- Sistema convierte: V_referencia = 480 / √3 = 277V (Vfn)
- e = 1 × 120 × 0.6374 × 0.030 = 2.295 V
- %e = (2.295 / 277) × 100 = **0.83%**
- `Cumple` = true

---

## Casos de Error

### Error 1: tipo_voltaje inválido

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70,
  "longitud_circuito": 30,
  "tension": 127,
  "tipo_voltaje": "TRIFASICO",
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Tipo de voltaje inválido",
  "code": "TIPO_VOLTAJE_INVALIDO",
  "details": "tipo de voltaje inválido: debe ser 'FASE_NEUTRO' o 'FASE_FASE'"
}
```

---

### Error 2: tipo_voltaje faltante

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70,
  "longitud_circuito": 30,
  "tension": 127,
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Error de validación",
  "code": "VALIDATION_ERROR",
  "details": "Key: 'CaidaTensionRequest.TipoVoltaje' Error:Field validation for 'TipoVoltaje' failed on the 'required' tag"
}
```

---

### Error 3: Sistema eléctrico inválido

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70,
  "longitud_circuito": 30,
  "tension": 127,
  "tipo_voltaje": "FASE_NEUTRO",
  "sistema_electrico": "HEXAFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Sistema eléctrico inválido",
  "code": "SISTEMA_ELECTRICO_INVALIDO"
}
```

---

## Tabla Comparativa de Resultados

Para facilitar la verificación, aquí está la tabla con los mismos parámetros base (2 AWG Cu, 70A, 30m, 1 hilo):

| Sistema Eléctrico | Voltaje Ingresado | Tipo Voltaje | V_ref (interno) | Factor | % Caída | Cumple |
|-------------------|-------------------|--------------|-----------------|--------|---------|--------|
| MONOFASICO | 127V | FASE_NEUTRO | 127V | 2 | 2.11% | ✅ |
| MONOFASICO | 220V | FASE_FASE | 127V (conv) | 2 | 2.11% | ✅ |
| BIFASICO | 127V | FASE_NEUTRO | 127V | 1 | 1.05% | ✅ |
| DELTA | 220V | FASE_FASE | 220V | √3 | 1.05% | ✅ |
| DELTA | 127V | FASE_NEUTRO | 220V (conv) | √3 | 1.05% | ✅ |
| ESTRELLA | 127V | FASE_NEUTRO | 127V | 1 | 1.05% | ✅ |
| ESTRELLA | 220V | FASE_FASE | 127V (conv) | 1 | 1.05% | ✅ |

**Relaciones matemáticas verificadas:**
- MONOFASICO = 2 × BIFASICO ✅ (ambos con Vfn 127V)
- DELTA = BIFASICO ✅ (en % caída, con sus respectivas V_ref)
- ESTRELLA = BIFASICO ✅ (ambos usan factor 1 con Vfn 127V)

---

## Script de Validación Rápida

```bash
#!/bin/bash

echo "=== Testing Voltage Drop Endpoint with tipo_voltaje ===" echo ""

# Test 1: MONOFASICO con Vfn
echo "1. MONOFASICO con Vfn 127V:"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":127,"tipo_voltaje":"FASE_NEUTRO","sistema_electrico":"MONOFASICO","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

echo ""
echo "2. BIFASICO con Vfn 127V:"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":127,"tipo_voltaje":"FASE_NEUTRO","sistema_electrico":"BIFASICO","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

echo ""
echo "3. DELTA con Vff 220V:"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":220,"tipo_voltaje":"FASE_FASE","sistema_electrico":"DELTA","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

echo ""
echo "4. ESTRELLA con Vfn 127V:"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":127,"tipo_voltaje":"FASE_NEUTRO","sistema_electrico":"ESTRELLA","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

echo ""
echo "=== Tests Complete ==="
```

---

## Changelog

### 2026-02-19 v2.0.0
- **BREAKING CHANGE**: Agregado campo obligatorio `tipo_voltaje`
- Implementada conversión automática Vff ↔ Vfn según sistema eléctrico
- Agregados valores correctos de voltaje de referencia:
  - MONOFASICO, BIFASICO, ESTRELLA → usan Vfn
  - DELTA → usa Vff
- Todos los tests actualizados con voltajes NOM reales (127V/220V)
- Documentado problema de error del 73% en versión anterior
- Agregada tabla comparativa con conversiones automáticas

### 2026-02-19 v1.0.0
- Implementación inicial con `sistema_electrico`
- Fórmula NOM simplificada
- Sin distinción de tipo de voltaje (bug detectado)
