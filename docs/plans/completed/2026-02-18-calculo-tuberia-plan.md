# Plan: Implementación de Cálculo de Tubería

**Fecha:** 2026-02-18
**Feature:** Endpoint de cálculo de tubería
**Basado en:** `docs/plans/2026-02-18-calculo-tuberia-design.md`

---

## Resumen del Plan

Este plan cubre la implementación completa del endpoint de cálculo de tubería según normativa NOM, siguiendo la arquitectura hexagonal del proyecto.

---

## Fase 1: Datos (CSV)

### Tarea 1.1: Crear tabla PVC ocupación 40%

**Archivo:** `data/tablas_nom/tubo-ocupacion-pvc-40.csv`

```
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

### Tarea 1.2: Crear tabla Acero PG ocupación 40%

**Archivo:** `data/tablas_nom/tubo-ocupacion-acero-pg-40.csv`

```
tamano,area_ocupacion_mm2,designacion_metrica,pulgadas
1/2,81,16,1/2
3/4,141,21,3/4
1,229,27,1
1 1/4,394,35,1 1/4
1 1/2,533,41,1 1/2
2,879,53,2
2 1/2,1255,63,2 1/2
3,1936,78,3
3 1/2,2584,91,3 1/2
4,3326,103,4
```

### Tarea 1.3: Crear tabla Acero PD ocupación 40%

**Archivo:** `data/tablas_nom/tubo-ocupacion-acero-pd-40.csv`

```
tamano,area_ocupacion_mm2,designacion_metrica,pulgadas
1/2,78,16,1/2
3/4,137,21,3/4
1,222,27,1
1 1/4,387,35,1 1/4
1 1/2,526,41,1 1/2
2,866,53,2
2 1/2,1513,63,2 1/2
3,2280,78,3
3 1/2,2980,91,3 1/2
4,3808,103,4
```

---

## Fase 2: Domain Layer

### Tarea 2.1: Crear Value Object para resultado de tubería

**Archivo:** `internal/calculos/domain/entity/tamanio_tuberia.go`

```go
type ResultadoTamanioTuberia struct {
    areaPorTuboMM2       float64
    tuberiaRecomendada   string
    designacionMetrica   string
    tipoCanalizacion     entity.TipoCanalizacion
    numTuberias          int
}
```

### Tarea 2.2: Crear servicio de cálculo de tubería

**Archivo:** `internal/calculos/domain/service/calcular_tamanio_tuberia.go`

**Responsabilidades:**
- Calcular área total por tubo según distribución
- Las tierras van completas en cada tubo (no se dividen)
- Buscar el tamaño de tubería donde área_ocupacion >= área_requerida
- Si no cabe → tomar el inmediato superior

**Métodos:**
- `CalcularAreaPorTubo(fases, neutros, tierras int, areaFase, areaNeutral, areaTierra float64, numTuberias int) float64`
- `BuscarTamanioTuberia(areaRequerida float64, tipoCanalizacion TipoCanalizacion, repo TablaOcupacionRepository) (ResultadoTamanioTuberia, error)`

---

## Fase 3: Application Layer

### Tarea 3.1: Extender port TablaNOMRepository

**Archivo:** `internal/calculos/application/port/tabla_nom_repository.go`

**Agregar métodos:**

```go
// ObtenerAreaConductor returns the area with insulation for a given calibre.
ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error)

// ObtenerTamanioTuberiaOcupacion returns the conduit sizing table entries for 40% fill.
ObtenerTamanioTuberiaOcupacion(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]TuboOcupacionEntry, error)
```

### Tarea 3.2: Agregar tipo para entrada de tabla tubería

**Archivo:** `internal/shared/kernel/valueobject/tubo_ocupacion.go` (nuevo)

```go
type TuboOcupacionEntry struct {
    Tamano              string
    AreaOcupacionMM2    float64
    DesignacionMetrica string
    Pulgadas           string
}
```

### Tarea 3.3: Crear DTO para input del use case

**Archivo:** `internal/calculos/application/dto/tuberia_input.go`

```go
type TuberiaInput struct {
    NumFases         int     `json:"num_fases" binding:"required,gt=0"`
    CalibreFase     string  `json:"calibre_fase" binding:"required"`
    NumNeutros      int     `json:"num_neutros" binding:"required,gte=0"`
    CalibreNeutro   string  `json:"calibre_neutral" binding:"required"`
    CalibreTierra   string  `json:"calibre_tierra" binding:"required"`
    TipoCanalizacion string  `json:"tipo_canalizacion" binding:"required"`
    NumTuberias     int     `json:"num_tuberias" binding:"required,gt=0"`
}
```

### Tarea 3.4: Crear DTO para output del use case

**Archivo:** `internal/calculos/application/dto/tuberia_output.go`

```go
type TuberiaOutput struct {
    AreaPorTuboMM2       float64 `json:"area_por_tubo_mm2"`
    TuberiaRecomendada  string  `json:"tuberia_recomendada"`
    DesignacionMetrica  string  `json:"designacion_metrica"`
    TipoCanalizacion    string  `json:"tipo_canalizacion"`
    NumTuberias         int     `json:"num_tuberias"`
}
```

### Tarea 3.5: Crear use case de cálculo de tubería

**Archivo:** `internal/calculos/application/usecase/calcular_tamanio_tuberia.go`

**Responsabilidades:**
1. Validar input DTO
2. Obtener áreas de conductores por calibre (usando repository)
3. Calcular distribución por tubo
4. Buscar tubería en tabla de ocupación
5. Retornar output DTO

---

## Fase 4: Infrastructure Layer

### Tarea 4.1: Implementar métodos del port en CSVTablaNOMRepository

**Archivo:** `internal/calculos/infrastructure/adapter/driven/csv/csv_tabla_nom_repository.go`

**Agregar:**
- `ObtenerAreaConductor()` - leer de `tabla-5-dimensiones-aislamiento.csv`, columna `area_tw_thw`
- `ObtenerTamanioTuberiaOcupacion()` - leer de las nuevas tablas CSV

### Tarea 4.2: Crear handler HTTP para tubería

**Archivo:** `internal/calculos/infrastructure/adapter/driver/http/tuberia_handler.go`

**Endpoint:** `POST /api/v1/calculos/tuberia`

**Request/Response:** Usar los DTOs de application

### Tarea 4.3: Registrar ruta en router

**Archivo:** `internal/calculos/infrastructure/router.go`

Agregar ruta para el nuevo handler.

---

## Fase 5: Wiring en main.go

### Tarea 5.1: Registrar dependencias

**Archivo:** `cmd/api/main.go`

- Crear instancia del nuevo use case
- Crear handler y registrar en router

---

## Fase 6: Verificación

### Tarea 6.1: Tests unitarios

```bash
go test ./internal/calculos/domain/service/... -run Tuberia
go test ./internal/calculos/application/usecase/... -run Tuberia
```

### Tarea 6.2: Build

```bash
go build ./...
```

### Tarea 6.3: Prueba manual del endpoint

```bash
curl -X POST http://localhost:8080/api/v1/calculos/tuberia \
  -H "Content-Type: application/json" \
  -d '{
    "num_fases": 3,
    "calibre_fase": "2 AWG",
    "num_neutros": 1,
    "calibre_neutral": "2 AWG",
    "calibre_tierra": "6 AWG",
    "tipo_canalizacion": "TUBERIA_PVC",
    "num_tuberias": 1
  }'
```

---

## Orden de Ejecución

| Fase | Responsable | Dependencias |
|------|-------------|--------------|
| 1. Datos CSV | Infraestructura | Ninguna |
| 2. Domain | Domain Agent | Ninguna |
| 3. Application | Application Agent | Domain |
| 4. Infrastructure | Infrastructure Agent | Application |
| 5. Wiring | Orquestador | Infrastructure |
| 6. Verificación | Todos | Todo |

---

## Notas

- Las tierras NO se dividen entre tubos - van completas en cada uno
- El área calculada es POR TUBO, no el total
- Tablas de ocupación ya incorporan el 40% de ocupación máxima
