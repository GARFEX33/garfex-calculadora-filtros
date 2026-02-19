# Diseño: Endpoint de Cálculo de Tubería

**Fecha:** 2026-02-18
**Feature:** Cálculo de tamaño de tubería según normativa NOM
**Arquitectura:** Clean + DDD + Hexagonal + Vertical Slices

---

## 1. Resumen

Nuevo endpoint REST para dimensionar tubería conduit según normativa NOM-001-SEDE. El cálculo se basa en la ocupación máxima del 40% del área interior de la tubería.

---

## 2. Input del Endpoint

```json
{
  "num_fases": 6,
  "calibre_fase": "2 AWG",
  "num_neutros": 2,
  "calibre_neutral": "2 AWG",
  "calibre_tierra": "6 AWG",
  "tipo_canalizacion": "TUBERIA_PVC",
  "num_tuberias": 1
}
```

### Campos

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|--------------|
| num_fases | int | Sí | Número de conductores de fase |
| calibre_fase | string | Sí | Calibre AWG (ej: "2 AWG") |
| num_neutros | int | Sí | Número de conductores neutrales |
| calibre_neutral | string | Sí | Calibre AWG |
| calibre_tierra | string | Sí | Calibre AWG |
| tipo_canalizacion | string | Sí | Tipo: TUBERIA_PVC, TUBERIA_ACERO_PG, TUBERIA_ACERO_PD |
| num_tuberias | int | Sí | Número de tuberías a dimensionar |

---

## 3. Output del Endpoint

```json
{
  "success": true,
  "data": {
    "area_por_tubo_mm2": 361.09,
    "tuberia_recomendada": "1 1/4",
    "designacion_metrica": "35",
    "tipo_canalizacion": "TUBERIA_PVC",
    "num_tuberias": 2
  }
}
```

---

## 4. Flujo de Cálculo

### Paso 1: Obtener áreas unitarias
- Consultar `tabla-5-dimensiones-aislamiento.csv` para obtener el área con aislamiento (`area_tw_thw`) de cada calibre

### Paso 2: Calcular áreas por grupo
- **Fases por tubo:** (num_fases ÷ num_tuberias) × área_fase
- **Neutros por tubo:** (num_neutros ÷ num_tuberias) × área_neutral
- **Tierra por tubo:** NO se divide → área_tierra × num_tuberias (va completa en cada tubo)

### Paso 3: Sumar áreas
- área_por_tubo = fases_por_tubo + neutros_por_tubo + tierra_por_tubo

### Paso 4: Buscar tubería
- En tabla de ocupación 40% (por tipo de tubería), buscar el tamaño donde área_ocupacion >= área_por_tubo
- Si no cabe en el primer tamaño → tomar el inmediato superior (siguiente tamaño)

---

## 5. Regla de Distribución de Tierras

Según NOM-001-SEDE, cada tubería debe contener su propio conductor de tierra. Las tierras **NO se dividen** entre tuberías - van completas en cada una.

---

## 6. Nuevas Tablas CSV

### Archivos a crear en `data/tablas_nom/`

| Archivo | Tipo Tubería | Artículo NOM |
|---------|--------------|--------------|
| `tubo-ocupacion-pvc-40.csv` | PVC CED 40 | Art 352-353 |
| `tubo-ocupacion-acero-pg-40.csv` | Acero PG (RMC/CED 40) | Art 344 |
| `tubo-ocupacion-acero-pd-40.csv` | Acero PD (EMT) | Art 358 |

### Formato CSV

```csv
tamano,area_ocupacion_mm2,designacion_metrica,pulgadas
1/2,74,16,1/2
3/4,131,21,3/4
1,214,27,1
1 1/4,374,35,1 1/4
1 1/2,513,41,1 1/2
2,849,53,2
2 1/2,1212,63,2 1/2
3,1877,78,3
3 1/2,2511,91,3 1/2
4,3237,103,4
```

### Datos (provenientes del usuario)

**TUBO ACERO TIPO RMC (CED 40) - Art 344, Pesada**
| Tamano | Area (mm²) | Metrica | Pulgadas |
|--------|------------|---------|----------|
| 1/2 | 81 | 16 | 1/2 |
| 3/4 | 141 | 21 | 3/4 |
| 1 | 229 | 27 | 1 |
| 1 1/4 | 394 | 35 | 1 1/4 |
| 1 1/2 | 533 | 41 | 1 1/2 |
| 2 | 879 | 53 | 2 |
| 2 1/2 | 1255 | 63 | 2 1/2 |
| 3 | 1936 | 78 | 3 |
| 3 1/2 | 2584 | 91 | 3 1/2 |
| 4 | 3326 | 103 | 4 |

**TUBO ACERO TIPO EMT - Art 358, Pared Delgada**
| Tamano | Area (mm²) | Metrica | Pulgadas |
|--------|------------|---------|----------|
| 1/2 | 78 | 16 | 1/2 |
| 3/4 | 137 | 21 | 3/4 |
| 1 | 222 | 27 | 1 |
| 1 1/4 | 387 | 35 | 1 1/4 |
| 1 1/2 | 526 | 41 | 1 1/2 |
| 2 | 866 | 53 | 2 |
| 2 1/2 | 1513 | 63 | 2 1/2 |
| 3 | 2280 | 78 | 3 |
| 3 1/2 | 2980 | 91 | 3 1/2 |
| 4 | 3808 | 103 | 4 |

**TUBO PVC CED 40 - Art 352 y 353**
| Tamano | Area (mm²) | Metrica | Pulgadas |
|--------|------------|---------|----------|
| 1/2 | 74 | 16 | 1/2 |
| 3/4 | 131 | 21 | 3/4 |
| 1 | 214 | 27 | 1 |
| 1 1/4 | 374 | 35 | 1 1/4 |
| 1 1/2 | 513 | 41 | 1 1/2 |
| 2 | 849 | 53 | 2 |
| 2 1/2 | 1212 | 63 | 2 1/2 |
| 3 | 1877 | 78 | 3 |
| 3 1/2 | 2511 | 91 | 3 1/2 |
| 4 | 3237 | 103 | 4 |

---

## 7. Estructura en Capas

### Domain Layer (`internal/calculos/domain/service/`)

- **Nuevo servicio:** `CalcularTamanioTuberia`
  - Recibe: áreas por grupo, tipo canalización
  - Retorna: tamaño de tubería recomendado
  - Lógica pura, sin I/O

### Application Layer (`internal/calculos/application/`)

- **Extender port:** `TablaNOMRepository`
  - Agregar método: `ObtenerAreaConductor(ctx, calibre, material, tipoAislamiento) float64`
  - Agregar método: `ObtenerTamanioTuberia(ctx, areaRequerida, tipoCanalizacion) (TamanioTuberia, error)`

- **Nuevo use case:** `CalcularTamanioTuberiaUseCase`
  - Orchestras: obtener áreas → calcular distribución → buscar tubería
  - Retorna DTO con resultado

### Infrastructure Layer (`internal/calculos/infrastructure/`)

- **Implementar métodos del port** en `CSVTablaNOMRepository`
- **Crear archivos CSV** de tablas de ocupación 40%
- **Nuevo handler:** `TuberiaHandler` con endpoint `POST /api/v1/calculos/tuberia`

---

## 8. Mapping de Tipos de Canalización

| Input API | Tabla CSV | Descripción |
|-----------|-----------|-------------|
| TUBERIA_PVC | tubo-ocupacion-pvc-40.csv | PVC CED 40 |
| TUBERIA_ACERO_PG | tubo-ocupacion-acero-pg-40.csv | Acero RMC/CED 40 |
| TUBERIA_ACERO_PD | tubo-ocupacion-acero-pd-40.csv | Acero EMT |

---

## 9. Consideraciones

- El área de aislamiento se toma de la columna `area_tw_thw` de `tabla-5-dimensiones-aislamiento.csv`
- Las tierras van completas en cada tubería (no se dividen)
- El cálculo es por tubo, no el total de todos los tubos
- Si el área requerida excede la máxima tabla → retornar error o la máxima disponible

---

## 10. Referencias

- NOM-001-SEDE-2012, Artículo 352 (PVC), 344 (RMC), 358 (EMT)
- Capítulo 9, Tabla 4 (dimensiones conduit)
- Tabla 5 NOM (dimensiones aislamiento conductores)
