# Pruebas Manuales - Endpoint Caída de Tensión

## Endpoint

```
POST /api/v1/calculos/caida-tension
```

## Fórmula NOM Implementada

El endpoint utiliza la **fórmula NOM simplificada** para calcular la caída de tensión:

```
e = factor × I × Z × L
%e = (e / V) × 100

Donde:
- Z = √(R² + X²)  [impedancia del conductor]
- R = resistencia Ω/km (Tabla 9 NOM)
- X = reactancia Ω/km (Tabla 9 NOM)
- I = corriente ajustada (A)
- L = longitud del circuito (km)
- V = tensión del sistema (V)
```

**Factores por sistema eléctrico:**
- `MONOFASICO` (1F-2H): **2**
- `BIFASICO` (2F-3H): **1**
- `DELTA` (3F-3H): **√3 ≈ 1.732**
- `ESTRELLA` (3F-4H): **1**

---

## Casos de Prueba por Sistema Eléctrico

### Caso 1: Sistema MONOFASICO (factor = 2)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
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
    "Porcentaje": 1.22,
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
- %e = (2.677 / 220) × 100 = 1.22%
- `Cumple` = true (< 3%)

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
    "tension": 220,
    "sistema_electrico": "MONOFASICO",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

### Caso 2: Sistema BIFASICO (factor = 1)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
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
    "Porcentaje": 0.61,
    "CaidaVolts": 1.34,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
- %e = (1.338 / 220) × 100 = 0.61%
- Exactamente **la mitad** que sistema monofásico

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
    "tension": 220,
    "sistema_electrico": "BIFASICO",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

### Caso 3: Sistema DELTA (factor = √3)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
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
- %e = (2.318 / 220) × 100 = 1.05%
- Factor √3 ≈ 1.732 para trifásico delta

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
    "tension": 220,
    "sistema_electrico": "DELTA",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

### Caso 4: Sistema ESTRELLA (factor = 1)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
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
    "Porcentaje": 0.61,
    "CaidaVolts": 1.34,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.637
  }
}
```

**Verificación:**
- e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
- %e = (1.338 / 220) × 100 = 0.61%
- Mismo resultado que BIFASICO (ambos usan factor 1)

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
    "tension": 220,
    "sistema_electrico": "ESTRELLA",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

### Caso 5: 2 Hilos por Fase (reduce impedancia a la mitad)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
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
    "Porcentaje": 0.61,
    "CaidaVolts": 1.34,
    "Cumple": true,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 0.319
  }
}
```

**Verificación:**
- Z_efectiva = 0.6374 / 2 = 0.3187 Ω/km
- e = 2 × 70 × 0.3187 × 0.030 = 1.338 V
- %e = (1.338 / 220) × 100 = 0.61%
- Exactamente **la mitad** que Caso 1 (1 hilo por fase)

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
    "tension": 220,
    "sistema_electrico": "MONOFASICO",
    "hilos_por_fase": 2,
    "limite_caida": 3.0
  }'
```

---

### Caso 6: Material Aluminio (Al)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Al",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Verificación:**
- Tabla 9 NOM para 2 AWG Al PVC: R=1.02 Ω/km, X=0.144 Ω/km
- Z = √(1.02² + 0.144²) = 1.030 Ω/km
- e = 2 × 70 × 1.030 × 0.030 = 4.326 V
- %e = (4.326 / 220) × 100 = 1.97%
- Mayor caída que cobre (Al tiene mayor resistencia)

**Comando curl:**
```bash
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{
    "calibre": "2 AWG",
    "material": "Al",
    "tipo_canalizacion": "TUBERIA_PVC",
    "corriente_ajustada": 70,
    "longitud_circuito": 30,
    "tension": 220,
    "sistema_electrico": "MONOFASICO",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

### Caso 7: Excede límite NOM (NO cumple)

**Request:**
```json
{
  "calibre": "14 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 25.0,
  "longitud_circuito": 100.0,
  "tension": 220,
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
    "Porcentaje": 11.84,
    "CaidaVolts": 26.05,
    "Cumple": false,
    "LimitePorcentaje": 3.0,
    "ResistenciaEfectiva": 5.21
  }
}
```

**Verificación:**
- Tabla 9 NOM para 14 AWG Cu PVC: R=5.21 Ω/km, X=0.157 Ω/km
- Z = √(5.21² + 0.157²) = 5.212 Ω/km
- e = 2 × 25 × 5.212 × 0.100 = 26.06 V
- %e = (26.06 / 220) × 100 = 11.84%
- `Cumple` = **false** (excede límite de 3%)

**Comando curl:**
```bash
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{
    "calibre": "14 AWG",
    "material": "Cu",
    "tipo_canalizacion": "TUBERIA_PVC",
    "corriente_ajustada": 25,
    "longitud_circuito": 100,
    "tension": 220,
    "sistema_electrico": "MONOFASICO",
    "hilos_por_fase": 1,
    "limite_caida": 3.0
  }'
```

---

## Casos de Error

### Error 1: Material inválido

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Ag",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Material del conductor inválido",
  "code": "MATERIAL_INVALIDO"
}
```

---

### Error 2: Sistema eléctrico inválido

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
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

### Error 3: Tipo de canalización inválido

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "DUCTO_MAGICO",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Tipo de canalización inválido",
  "code": "TIPO_CANALIZACION_INVALIDO"
}
```

---

### Error 4: Calibre no encontrado en Tabla 9

**Request:**
```json
{
  "calibre": "99 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 70.0,
  "longitud_circuito": 30.0,
  "tension": 220,
  "sistema_electrico": "MONOFASICO",
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "No se encontró la impedancia para el calibre y canalización",
  "code": "IMPEDANCIA_NO_ENCONTRADA"
}
```

---

### Error 5: Validación de campos obligatorios

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC"
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Error de validación",
  "code": "VALIDATION_ERROR",
  "details": "Key: 'CaidaTensionRequest.corriente_ajustada' Error:Field validation for 'corriente_ajustada' failed on the 'required' tag"
}
```

---

## Tabla Comparativa de Resultados

Para facilitar la verificación, aquí está la tabla con los mismos parámetros base (2 AWG Cu, 70A, 30m, 220V, 1 hilo):

| Sistema Eléctrico | Factor | % Caída | Caída (V) | Cumple (< 3%) |
|-------------------|--------|---------|-----------|---------------|
| MONOFASICO        | 2      | 1.22%   | 2.68 V    | ✅ Sí         |
| BIFASICO          | 1      | 0.61%   | 1.34 V    | ✅ Sí         |
| DELTA             | √3     | 1.05%   | 2.32 V    | ✅ Sí         |
| ESTRELLA          | 1      | 0.61%   | 1.34 V    | ✅ Sí         |

**Relación de factores:**
- MONOFASICO es **2x** BIFASICO/ESTRELLA
- DELTA es **√3x** BIFASICO/ESTRELLA (≈1.73x)

---

## Notas Técnicas

1. **Impedancia Z**: Se calcula con `√(R² + X²)` según Tabla 9 NOM
2. **Valores R y X**: Dependen de:
   - Calibre del conductor
   - Material (Cu o Al)
   - Tipo de canalización (PVC, metálica, charola)
3. **Hilos por fase**: Divide la impedancia efectiva (conexión en paralelo)
4. **ResistenciaEfectiva**: El campo reporta `Z` (impedancia), no la resistencia pura R
5. **Tolerancia**: ±0.01 en valores decimales por redondeo de punto flotante

---

## Validación Rápida (Smoke Test)

Script bash para validar los 4 sistemas eléctricos:

```bash
#!/bin/bash

echo "=== Testing Voltage Drop Endpoint ==="

# Test MONOFASICO
echo -e "\n1. MONOFASICO (factor 2):"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":220,"sistema_electrico":"MONOFASICO","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

# Test BIFASICO
echo -e "\n2. BIFASICO (factor 1):"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":220,"sistema_electrico":"BIFASICO","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

# Test DELTA
echo -e "\n3. DELTA (factor √3):"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":220,"sistema_electrico":"DELTA","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

# Test ESTRELLA
echo -e "\n4. ESTRELLA (factor 1):"
curl -s -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":70,"longitud_circuito":30,"tension":220,"sistema_electrico":"ESTRELLA","hilos_por_fase":1,"limite_caida":3.0}' \
  | python -m json.tool

echo -e "\n=== Tests Complete ==="
```

---

## Changelog

### 2026-02-19
- **BREAKING CHANGE**: Reemplazado campo `factor_potencia` con `sistema_electrico`
- Implementada fórmula NOM simplificada: `e = factor × I × Z × L`
- Agregados factores correctos por sistema: MONOFASICO=2, BIFASICO=1, DELTA=√3, ESTRELLA=1
- Actualizado campo `ResistenciaEfectiva` para reportar impedancia Z
- Eliminados casos de prueba con factor de potencia variable
- Agregados casos de prueba para los 4 sistemas eléctricos
