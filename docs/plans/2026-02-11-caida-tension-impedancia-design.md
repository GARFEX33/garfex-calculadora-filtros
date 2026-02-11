# Diseño: CalcularCaidaTension con Método de Impedancia (NOM)

**Fecha:** 2026-02-11
**Estado:** Validado
**Contexto:** Fase 1 — Domain Layer — Reemplazo del servicio simplificado

## Problema

El servicio actual `CalcularCaidaTension` usa la fórmula simplificada de resistividad:

```
VD% = (√3 × ρ × L × I) / (S × V) × 100
```

Donde ρ es una constante fija (Cu: 0.01724, Al: 0.02826 Ω·mm²/m).

Esto **no es el método NOM correcto**. La normativa usa impedancia (R + jX) que toma en cuenta:
- Resistencia AC real (varía por material del conduit: PVC, Aluminio, Acero)
- Reactancia inductiva (depende de la geometría de instalación)

## Método NOM: Impedancia

La fórmula correcta es:

```
VD = √3 × I × Z × L_km
%VD = (VD / V) × 100
```

Donde **Z = √(R² + X²)** y:
- **R** = Resistencia AC del conductor (Ω/km) — de Tabla 9, varía por material conductor + material conduit
- **X** = Reactancia inductiva (Ω/km) — calculada geométricamente con DMG/RMG

### Por qué DMG/RMG en vez de Tabla 9 para reactancia

La Tabla 9 tiene columnas de reactancia solo para tubería (Aluminio y Acero). No cubre charola. Para mantener un método unificado y consistente para TODOS los tipos de canalización, se calcula X geométricamente:

```
X = (1/n) × 2π × 60 × 2 × 10⁻⁷ × ln(DMG/RMG) × 1000  [Ω/km]
```

Donde:
- **n** = hilos por fase (conductores en paralelo)
- **60** = frecuencia Hz (México)
- **DMG** = Distancia Media Geométrica entre fases (mm)
- **RMG** = Radio Medio Geométrico del conductor (mm)

## Cálculo de RMG

```
RMG = (diametro_conductor_mm / 2) × factor_hilos
```

**Datos de entrada:** Tabla 8 (`diametro_mm`, `numero_hilos`)

| Número de hilos | Factor |
|----------------|--------|
| 1 (sólido) | 0.7788 |
| 7 | 0.726 |
| 19 | 0.758 |
| 37 | 0.768 |
| 61 | 0.772 |

Estos factores son constantes electromagnéticas estándar (no NOM), se definen en el domain service.

## Cálculo de DMG

```
DMG = diametro_exterior_mm × factor_canalizacion
```

**Datos de entrada:** Tabla 5 (`diam_tw_thw` — siempre THW) para diámetro exterior con aislamiento.

| TipoCanalizacion | Factor DMG | Razón |
|-----------------|-----------|-------|
| TUBERIA_PVC | 1.0 | Cables adyacentes en tubería |
| TUBERIA_ALUMINIO | 1.0 | Cables adyacentes en tubería |
| TUBERIA_ACERO_PG | 1.0 | Cables adyacentes en tubería |
| TUBERIA_ACERO_PD | 1.0 | Cables adyacentes en tubería |
| CHAROLA_CABLE_ESPACIADO | 2.0 | Cables separados un diámetro |
| CHAROLA_CABLE_TRIANGULAR | 1.0 | Cables tocándose en triángulo |

## TipoCanalizacion — Enum Expandido (6 valores)

El enum actual tiene 3 valores genéricos. Se expande a 6 para capturar el material del conduit (necesario para seleccionar columna de resistencia en Tabla 9):

```go
const (
    TipoCanalizacionTuberiaPVC             TipoCanalizacion = "TUBERIA_PVC"
    TipoCanalizacionTuberiaAluminio        TipoCanalizacion = "TUBERIA_ALUMINIO"
    TipoCanalizacionTuberiaAceroPG         TipoCanalizacion = "TUBERIA_ACERO_PG"
    TipoCanalizacionTuberiaAceroPD         TipoCanalizacion = "TUBERIA_ACERO_PD"
    TipoCanalizacionCharolaCableEspaciado  TipoCanalizacion = "CHAROLA_CABLE_ESPACIADO"
    TipoCanalizacionCharolaCableTriangular TipoCanalizacion = "CHAROLA_CABLE_TRIANGULAR"
)
```

### Mapeo a Tabla NOM de Ampacidad

Los 4 tipos de tubería comparten la misma tabla de ampacidad:

| TipoCanalizacion | Tabla ampacidad | Columna R (Tabla 9) |
|-----------------|-----------------|---------------------|
| TUBERIA_PVC | 310-15-b-16.csv | `res_{material}_pvc` |
| TUBERIA_ALUMINIO | 310-15-b-16.csv | `res_{material}_al` |
| TUBERIA_ACERO_PG | 310-15-b-16.csv | `res_{material}_acero` |
| TUBERIA_ACERO_PD | 310-15-b-16.csv | `res_{material}_acero` |
| CHAROLA_CABLE_ESPACIADO | 310-15-b-17.csv | `res_{material}_pvc` (nota 1) |
| CHAROLA_CABLE_TRIANGULAR | 310-15-b-20.csv | `res_{material}_pvc` (nota 1) |

**Nota 1:** Charola no tiene conduit metálico, por lo que la resistencia AC es la misma que en PVC (sin efecto de proximidad del conduit).

### Mapeo PG vs PD

- **PG (Pared Gruesa):** Tubería de acero con pared gruesa — usa columna `acero` de Tabla 9
- **PD (Pared Delgada):** Tubería de acero con pared delgada — usa misma columna `acero` de Tabla 9

Ambas usan la misma resistencia AC. La diferencia entre PG y PD afecta las dimensiones de canalización (Tabla de tuberías), no la impedancia eléctrica.

## Firma del Servicio (Domain)

```go
// EntradaCalculoCaidaTension contiene los datos pre-resueltos
// necesarios para calcular caída de tensión por método de impedancia.
// La capa application resuelve los valores desde las tablas NOM.
type EntradaCalculoCaidaTension struct {
    ResistenciaOhmPorKm float64           // Tabla 9 → columna según material + canalización
    DiametroExteriorMM  float64           // Tabla 5 → diam_tw_thw (siempre THW)
    DiametroConductorMM float64           // Tabla 8 → diametro_mm (conductor desnudo)
    NumeroHilos         int               // Tabla 8 → numero_hilos
    TipoCanalizacion    TipoCanalizacion  // Determina factor DMG
    HilosPorFase        int               // Conductores en paralelo (≥1)
}

func CalcularCaidaTension(
    entrada EntradaCalculoCaidaTension,
    corriente valueobject.Corriente,
    distancia float64,               // metros
    tension valueobject.Tension,
    limiteNOM float64,               // porcentaje (3.0 o 5.0)
) (resultado ResultadoCaidaTension, err error)
```

### ResultadoCaidaTension (nuevo struct de retorno)

```go
type ResultadoCaidaTension struct {
    Porcentaje       float64  // %VD
    CaidaVolts       float64  // VD en volts
    Cumple           bool     // %VD <= limiteNOM
    Impedancia       float64  // Z (Ω/km) — para reporte
    Resistencia      float64  // R (Ω/km) — para reporte
    Reactancia       float64  // X (Ω/km) — para reporte
}
```

**Razón del struct:** El reporte de memoria de cálculo necesita mostrar R, X, Z como valores intermedios, no solo el resultado final.

## Flujo Interno del Servicio (7 pasos)

```
1. RMG  = (DiametroConductorMM / 2) × factorHilos[NumeroHilos]
2. DMG  = DiametroExteriorMM × factorDMG[TipoCanalizacion]
3. X    = (1/HilosPorFase) × 2π×60 × 2×10⁻⁷ × ln(DMG/RMG) × 1000  [Ω/km]
4. R    = ResistenciaOhmPorKm / HilosPorFase
5. Z    = √(R² + X²)
6. VD   = √3 × corriente × Z × (distancia / 1000)
7. %VD  = (VD / tensión) × 100
```

## Tablas NOM de Referencia (CSVs validados)

| CSV | Tabla NOM | Dato que aporta | Validado |
|-----|-----------|-----------------|----------|
| `tabla-9-resistencia-reactancia.csv` | Tabla 9 | R (Ω/km) por calibre + material conduit | ✅ 100% |
| `tabla-5-dimensiones-aislamiento.csv` | Tabla 5 | Diámetro exterior THW (mm) para DMG | ✅ 100% |
| `tabla-8-conductor-desnudo.csv` | Tabla 8 | Diámetro desnudo (mm) + número hilos para RMG | ✅ 100% |

## Impacto en Código Existente

### Cambios en domain (Fase 1)

| Cambio | Archivo | Descripción |
|--------|---------|-------------|
| **Reescribir** servicio | `calculo_caida_tension.go` | Nueva firma + método impedancia |
| **Reescribir** tests | `calculo_caida_tension_test.go` | Tests con datos reales NOM |
| **Nuevo** struct retorno | `calculo_caida_tension.go` | `ResultadoCaidaTension` |
| **Nuevo** struct entrada | `calculo_caida_tension.go` | `EntradaCalculoCaidaTension` |
| **Expandir** enum | `tipo_canalizacion.go` | 3 → 6 valores |
| **Actualizar** tests | `tipo_canalizacion_test.go` | Tests para 6 valores |
| **Agregar** campo | `memoria_calculo.go` | `TemperaturaUsada int` + `ResultadoCaidaTension` |

### Sin cambios

- `SeleccionarConductorAlimentacion` — no cambia (ya recibe datos pre-resueltos)
- `SeleccionarConductorTierra` — independiente
- `CalcularCanalizacion` — independiente (usa dimensiones, no impedancia)
- Tabla 250-122 — tabla de tierra, sin relación

### Cambios en infrastructure/application (Fase 2)

| Cambio | Capa | Descripción |
|--------|------|-------------|
| CSV reader Tabla 9 | infrastructure | Leer R según calibre + canalización |
| CSV reader Tabla 5 | infrastructure | Leer diámetro THW por calibre |
| CSV reader Tabla 8 | infrastructure | Leer diámetro + hilos por calibre |
| Resolver `EntradaCalculoCaidaTension` | application | Mapear canalización → columna R correcta |

## Constantes del Domain Service

```go
// factorHilos mapea número de hilos → factor RMG (constantes electromagnéticas).
var factorHilos = map[int]float64{
    1:  0.7788,
    7:  0.726,
    19: 0.758,
    37: 0.768,
    61: 0.772,
}

// factorDMG mapea tipo de canalización → multiplicador de diámetro exterior para DMG.
var factorDMG = map[TipoCanalizacion]float64{
    TipoCanalizacionTuberiaPVC:             1.0,
    TipoCanalizacionTuberiaAluminio:        1.0,
    TipoCanalizacionTuberiaAceroPG:         1.0,
    TipoCanalizacionTuberiaAceroPD:         1.0,
    TipoCanalizacionCharolaCableEspaciado:  2.0,
    TipoCanalizacionCharolaCableTriangular: 1.0,
}
```

## Ejemplo de Validación

Datos de ejemplo (de imagen de cálculo NOM proporcionada por usuario):
- Conductor: 2 AWG Cu, THW, en tubería PVC
- Corriente: 120 A, distancia: 30 m, tensión: 480 V

Lookup de tablas:
- Tabla 9: `res_cu_pvc` para 2 AWG = **0.62** Ω/km
- Tabla 5: `diam_tw_thw` para 2 AWG = **10.46** mm
- Tabla 8: `diametro_mm` = **7.42** mm, `numero_hilos` = **7**

Cálculo:
```
1. RMG = (7.42/2) × 0.726 = 3.71 × 0.726 = 2.6935 mm
2. DMG = 10.46 × 1.0 = 10.46 mm  (tubería → factor 1.0)
3. X   = (1/1) × 2π×60 × 2×10⁻⁷ × ln(10.46/2.6935) × 1000
       = 376.99 × 2×10⁻⁷ × 1.3567 × 1000
       = 376.99 × 0.0002 × 1.3567 × 1000
       = 376.99 × 0.00027134
       = 0.1023 Ω/km
4. R   = 0.62 / 1 = 0.62 Ω/km
5. Z   = √(0.62² + 0.1023²) = √(0.3844 + 0.01047) = √0.3949 = 0.6284 Ω/km
6. VD  = √3 × 120 × 0.6284 × (30/1000) = 207.85 × 0.6284 × 0.03 = 3.917 V
7. %VD = (3.917 / 480) × 100 = 0.816%
```

## Decisiones Clave

1. **DMG/RMG para todos los tipos** — método unificado, no depende de columnas de reactancia de Tabla 9
2. **R desde Tabla 9** — resistencia AC real, no resistividad teórica
3. **6 tipos de canalización** — capturan material del conduit para seleccionar columna R correcta
4. **Siempre THW** — para diámetro exterior (usuario siempre usa este tipo)
5. **Struct de retorno** — expone R, X, Z para el reporte de memoria de cálculo
6. **Domain recibe datos pre-resueltos** — no conoce CSV, mantiene arquitectura hexagonal
7. **Charola usa `res_pvc`** — sin conduit metálico, misma resistencia AC que en PVC
