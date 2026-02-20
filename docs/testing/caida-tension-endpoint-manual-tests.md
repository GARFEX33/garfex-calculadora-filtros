# Pruebas Manuales - Endpoint Caída de Tensión

## Endpoint

```
POST /api/v1/calculos/caida-tension
```

## Casos de Prueba

### Caso 1: Factor de Potencia = 1.0 (Filtro Activo/Transformador)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 1.0,
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "porcentaje": 0.80445,
    "caida_volts": 3.861,
    "cumple": true,
    "limite_porcentaje": 3.0,
    "resistencia_efectiva": 0.62
  }
}
```

**Verificación:**
- `porcentaje` ≈ 0.8045% (±0.01)
- `caida_volts` ≈ 3.86V (±0.01)
- `cumple` = true (menor que límite 3%)
- `resistencia_efectiva` = 0.62 Ω/km (solo resistencia, X no contribuye con FP=1)

---

### Caso 2: Factor de Potencia = 0.85 (Carga)

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 0.85,
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "porcentaje": 0.719,
    "caida_volts": 3.451,
    "cumple": true,
    "limite_porcentaje": 3.0,
    "resistencia_efectiva": 0.5539
  }
}
```

**Verificación:**
- `porcentaje` ≈ 0.719% (±0.01)
- `caida_volts` ≈ 3.45V (±0.01)
- `cumple` = true
- `resistencia_efectiva` ≈ 0.554 Ω/km (R·cosθ + X·senθ)

---

### Caso 3: 2 Hilos por Fase

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 1.0,
  "hilos_por_fase": 2,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "porcentaje": 0.4035,
    "caida_volts": 1.937,
    "cumple": true,
    "limite_porcentaje": 3.0,
    "resistencia_efectiva": 0.31
  }
}
```

**Verificación:**
- `porcentaje` ≈ 0.40% (±0.01) — mitad del caso 1
- `caida_volts` ≈ 1.94V (±0.01)
- `resistencia_efectiva` = 0.31 Ω/km (R/2)

---

### Caso 4: Charola Espaciado

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "CHAROLA_CABLE_ESPACIADO",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 1.0,
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "porcentaje": 0.80445,
    "caida_volts": 3.861,
    "cumple": true,
    "limite_porcentaje": 3.0,
    "resistencia_efectiva": 0.62
  }
}
```

**Verificación:**
- Mismo resultado que Caso 1 (charola usa `reactancia_al`)

---

### Caso 5: Excede límite NOM

**Request:**
```json
{
  "calibre": "14 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 25.0,
  "longitud_circuito": 100.0,
  "tension": 220,
  "factor_potencia": 1.0,
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": true,
  "data": {
    "porcentaje": 5.82,
    "caida_volts": 12.80,
    "cumple": false,
    "limite_porcentaje": 3.0,
    "resistencia_efectiva": 10.2
  }
}
```

**Verificación:**
- `cumple` = false (excede 3%)
- `porcentaje` > 3.0

---

## Casos de Error

### Error 1: Material inválido

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Ag",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 1.0,
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

### Error 2: Factor de potencia fuera de rango

**Request:**
```json
{
  "calibre": "2 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 1.5,
  "hilos_por_fase": 1,
  "limite_caida": 3.0
}
```

**Response esperada:**
```json
{
  "success": false,
  "error": "Error de validación",
  "code": "VALIDATION_ERROR"
}
```

---

### Error 3: Calibre no encontrado

**Request:**
```json
{
  "calibre": "99 AWG",
  "material": "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_ajustada": 120.0,
  "longitud_circuito": 30.0,
  "tension": 480,
  "factor_potencia": 1.0,
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

## Comandos curl

```bash
# Caso 1 - FP=1.0
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":120.0,"longitud_circuito":30.0,"tension":480,"factor_potencia":1.0,"hilos_por_fase":1,"limite_caida":3.0}'

# Caso 2 - FP=0.85
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":120.0,"longitud_circuito":30.0,"tension":480,"factor_potencia":0.85,"hilos_por_fase":1,"limite_caida":3.0}'

# Caso 3 - 2 hilos
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","tipo_canalizacion":"TUBERIA_PVC","corriente_ajustada":120.0,"longitud_circuito":30.0,"tension":480,"factor_potencia":1.0,"hilos_por_fase":2,"limite_caida":3.0}'
```

## Notas

- Todos los valores numéricos tienen tolerancia de ±0.01 para decimales
- El campo `resistencia_efectiva` representa el término efectivo IEEE-141: R·cosθ + X·senθ
- Para FP=1.0, la reactancia no contribuye (senθ=0), solo la resistencia
- Para FP<1.0, ambos R y X contribuyen según el ángulo de la carga
