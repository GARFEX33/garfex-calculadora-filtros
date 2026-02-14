# Diseño: Tablas NOM de Ampacidad y Selección por Canalización

**Fecha:** 2026-02-11
**Estado:** Validado
**Contexto:** Fase 1 — Domain Layer

## Problema

El diseño original asumía una sola tabla de conductores para selección de alimentación. En realidad, la **tabla NOM de ampacidad depende del tipo de canalización**:

| Tabla NOM | Tipo de Canalización |
|-----------|---------------------|
| 310-15(b)(16) | Tubería conduit |
| 310-15(b)(17) | Charola — cables con espaciado |
| 310-15(b)(20) | Charola — arreglo triangular |

Cada tabla tiene los mismos calibres pero **diferentes valores de ampacidad** para cada uno.

## Flujo Revisado de Cálculo

El flujo correcto es:

```
1. Calcular corriente nominal (In)
2. Ajustar corriente (factores)
3. *** Seleccionar tipo de canalización ***     ← NUEVO: primero la canalización
4. Cargar tabla NOM correcta según canalización
5. Auto-seleccionar columna de temperatura
6. Seleccionar conductor de alimentación
7. Seleccionar conductor de tierra (tabla 250-122, independiente)
8. Calcular canalización (dimensiones)
9. Calcular caída de tensión
```

**Cambio clave:** La canalización se decide ANTES de seleccionar el conductor, porque determina qué tabla usar.

## Formato CSV Universal

Todas las tablas de ampacidad usan el mismo formato CSV:

```csv
seccion_mm2,calibre,cu_60c,cu_75c,cu_90c,al_60c,al_75c,al_90c
```

- Valores vacíos donde la columna no aplica (ej: 310-15(b)(20) no tiene 60°C)
- Valores vacíos donde el calibre no tiene dato de aluminio (calibres muy pequeños)
- Calibres en formato estándar: `14 AWG`, `4/0 AWG`, `250 MCM`, etc.

### Archivos CSV — Fase 1

| Archivo | Tabla NOM | Rango calibres | Columnas temperatura |
|---------|-----------|----------------|---------------------|
| `310-15-b-16.csv` | 310-15(b)(16) | 14 AWG – 2000 MCM | Cu: 60/75/90°C, Al: 60/75/90°C |
| `310-15-b-17.csv` | 310-15(b)(17) | 14 AWG – 2000 MCM | Cu: 60/75/90°C, Al: 60/75/90°C |
| `310-15-b-20.csv` | 310-15(b)(20) | 8 AWG – 1000 MCM | Cu: 75/90°C, Al: 75/90°C (sin 60°C) |
| `250-122.csv` | 250-122 | (tierra) | Independiente |

## TipoCanalizacion — Enum de Dominio

Nuevo value object en `internal/domain/entity/tipo_canalizacion.go`:

```go
type TipoCanalizacion string

const (
    TipoCanalizacionTuberiaPVC             TipoCanalizacion = "TUBERIA_PVC"
    TipoCanalizacionTuberiaAluminio        TipoCanalizacion = "TUBERIA_ALUMINIO"
    TipoCanalizacionTuberiaAceroPG         TipoCanalizacion = "TUBERIA_ACERO_PG"
    TipoCanalizacionTuberiaAceroPD         TipoCanalizacion = "TUBERIA_ACERO_PD"
    TipoCanalizacionCharolaCableEspaciado  TipoCanalizacion = "CHAROLA_CABLE_ESPACIADO"
    TipoCanalizacionCharolaCableTriangular TipoCanalizacion = "CHAROLA_CABLE_TRIANGULAR"
)
```

**Naming rationale:** Nombres descriptivos (no códigos NOM) para escalabilidad futura.
**Expansión (2026-02-11):** Se expandió de 3 a 6 valores para capturar el material del conduit, necesario para seleccionar la columna correcta de resistencia AC en Tabla 9 para caída de tensión por impedancia.

### Mapeo a tabla de ampacidad CSV

```go
// Los 4 tipos de tubería comparten la misma tabla de ampacidad
var tablaAmpacidad = map[TipoCanalizacion]string{
    TipoCanalizacionTuberiaPVC:             "310-15-b-16.csv",
    TipoCanalizacionTuberiaAluminio:        "310-15-b-16.csv",
    TipoCanalizacionTuberiaAceroPG:         "310-15-b-16.csv",
    TipoCanalizacionTuberiaAceroPD:         "310-15-b-16.csv",
    TipoCanalizacionCharolaCableEspaciado:  "310-15-b-17.csv",
    TipoCanalizacionCharolaCableTriangular: "310-15-b-20.csv",
}
```

### Mapeo a columna de resistencia (Tabla 9)

```go
// Cada tipo de canalización mapea a una columna de resistencia AC en Tabla 9
// Charola no tiene conduit metálico → usa columna PVC (sin efecto de proximidad)
var columnaResistencia = map[TipoCanalizacion]string{
    TipoCanalizacionTuberiaPVC:             "res_{material}_pvc",
    TipoCanalizacionTuberiaAluminio:        "res_{material}_al",
    TipoCanalizacionTuberiaAceroPG:         "res_{material}_acero",
    TipoCanalizacionTuberiaAceroPD:         "res_{material}_acero",
    TipoCanalizacionCharolaCableEspaciado:  "res_{material}_pvc",
    TipoCanalizacionCharolaCableTriangular: "res_{material}_pvc",
}
```

Estos mapeos viven en **infrastructure** (CSV repository), no en domain.

### Relación con Caída de Tensión

El TipoCanalizacion expandido también determina el factor DMG para el cálculo de reactancia inductiva (método impedancia). Ver diseño completo: `docs/plans/2026-02-11-caida-tension-impedancia-design.md`

## Selección Automática de Temperatura

Según normativa NOM:

| Condición | Columna |
|-----------|---------|
| Circuitos ≤ 100 A o calibres 14–1 AWG | 60°C |
| Circuitos > 100 A o calibres > 1 AWG | 75°C |
| Override explícito (todos los equipos rated 90°C) | 90°C |

### Reglas de implementación

1. **Auto-selección por defecto:** El sistema decide 60°C o 75°C según la corriente ajustada
2. **90°C nunca se auto-selecciona:** Requiere `temperatura_override: 90` explícito en el input del usuario
3. **Charola triangular (310-15(b)(20)):** No tiene columna 60°C → si la regla dice 60°C, se usa 75°C automáticamente
4. **TemperaturaUsada:** Se registra en `MemoriaCalculo` para el reporte (qué temperatura se usó realmente)

### Lógica (pseudocódigo)

```
func seleccionarTemperatura(corrienteAjustada, tipoCanalizacion, override):
    if override == 90:
        return 90
    if corrienteAjustada <= 100:
        if tipoCanalizacion tiene columna 60°C:
            return 60
        else:
            return 75  // fallback para charola triangular
    return 75
```

## Impacto en Código Existente

### Sin cambios (domain services)

`SeleccionarConductorAlimentacion` **no cambia**. Ya recibe `[]EntradaTablaConductor` con un campo `Capacidad` resuelto. La resolución de qué tabla/columna usar ocurre ANTES de llamar al servicio.

### Cambios en domain (Fase 1)

| Cambio | Archivo | Descripción |
|--------|---------|-------------|
| Expandir enum | `tipo_canalizacion.go` | `TipoCanalizacion` 3→6 valores (incluye material conduit) |
| Actualizar struct | `canalizacion.go` | Usar `TipoCanalizacion` en vez de `string` |
| Nuevo campo | `memoria_calculo.go` | Agregar `TemperaturaUsada int` + `ResultadoCaidaTension` |

### Cambios en infrastructure/application (Fase 2)

| Cambio | Capa | Descripción |
|--------|------|-------------|
| CSV reader | infrastructure | Leer tabla correcta según TipoCanalizacion |
| Mapeo tabla→CSV | infrastructure | `map[TipoCanalizacion]string` |
| Lógica temperatura | application/usecase | Auto-selección de columna |
| Input ampliado | application/dto | `TipoCanalizacion` + `TemperaturaOverride` opcionales |

## Decisiones Clave

1. **Nombres descriptivos** para TipoCanalizacion (no códigos NOM) — escalabilidad
2. **6 valores** para capturar material del conduit (necesario para Tabla 9 de resistencia)
3. **Auto-selección de temperatura** con override explícito para 90°C — seguridad NOM
4. **CSV universal** con mismas columnas — simplifica el parser
5. **Domain services sin cambio** — la resolución tabla/columna es responsabilidad de capas superiores
6. **Charola triangular sin 60°C** — fallback automático a 75°C, documentado en MemoriaCalculo
7. **Método impedancia para caída de tensión** — usa TipoCanalizacion para factor DMG y columna R (ver diseño: `2026-02-11-caida-tension-impedancia-design.md`)
