# Plan de Implementación de Infrastructure - Cálculo de Tubería

## Tareas

### Tarea 1: Crear tablas CSV de ocupación de tubería (40%)
**Archivos a crear en `data/tablas_nom/`:**
- [x] `tubo-ocupacion-pvc-40.csv`
- [x] `tubo-ocupacion-acero-pg-40.csv`  
- [x] `tubo-ocupacion-acero-pd-40.csv`

### Tarea 2: Implementar métodos del port en CSVTablaNOMRepository
**Archivo:** `internal/calculos/infrastructure/adapter/driven/csv/csv_tabla_nom_repository.go`
- [x] Agregar estructura `tuboOcupacionEntry` para cache
- [x] Agregar campo `tablasOcupacionTuberia` al struct
- [x] Cargar tablas de ocupación en `NewCSVTablaNOMRepository`
- [x] Implementar `ObtenerAreaConductor(ctx, calibre)` - leer de `tabla-5-dimensiones-aislamiento.csv`, columna `area_tw_thw`
- [x] Implementar `ObtenerTablaOcupacionTuberia(ctx, canalizacion)` - leer de las nuevas tablas

### Tarea 3: Crear handler HTTP para tubería
**Archivo:** `internal/calculos/infrastructure/adapter/driver/http/tuberia_handler.go`
- [x] Estructura TuberiaHandler con use case
- [x] NewTuberiaHandler constructor
- [x] Request/Response types
- [x] CalcularTuberia endpoint

### Tarea 4: Registrar ruta en router
**Archivo:** `internal/calculos/infrastructure/router.go`
- [x] Agregar use case al router
- [x] Agregar POST /api/v1/calculos/tuberia

### Tarea 5: Verificar
- [x] `go build ./internal/calculos/...` ✅

## Estado: COMPLETADO ✅
