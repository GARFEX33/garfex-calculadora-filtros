# Plan: Capa de Application - Cálculo de Tubería

## Objetivo
Implementar la capa de aplicación para el cálculo de tamaño de tubería según normativa NOM.

## Tareas

### Tarea 3.1: Extender port TablaNOMRepository
- **Archivo:** `internal/calculos/application/port/tabla_nom_repository.go`
- **Agregar métodos:**
  - `ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error)` - Retorna el área con aislamiento (area_tw_thw)
  - `ObtenerTablaOcupacionTuberia(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]service.EntradaTablaOcupacion, error)` - Retorna la tabla de ocupación para 40% fill

### Tarea 3.2: Crear DTOs para tubería
- **TuberiaInput:** `internal/calculos/application/dto/tuberia_input.go`
  - NumFases int
  - CalibreFase string
  - NumNeutros int
  - CalibreNeutro string
  - CalibreTierra string
  - TipoCanalizacion string
  - NumTuberias int
  
- **TuberiaOutput:** `internal/calculos/application/dto/tuberia_output.go`
  - AreaPorTuboMM2 float64
  - TuberiaRecomendada string
  - DesignacionMetrica string
  - TipoCanalizacion string
  - NumTuberias int

### Tarea 3.3: Crear Use Case de cálculo de tubería
- **Archivo:** `internal/calculos/application/usecase/calcular_tamanio_tuberia.go`
- **Responsabilidades:**
  1. Validar input DTO
  2. Parsear tipo canalización
  3. Obtener áreas de conductores por calibre (usando repository) - área_tw_thw
  4. Obtener tabla de ocupación (usando repository)
  5. Llamar al servicio de dominio `CalcularTamanioTuberiaWithMultiplePipes`
  6. Retornar output DTO

## Dependencias
- domain/service/calcular_tamanio_tuberia.go
- domain/entity/tamanio_tuberia.go
- application/port/tabla_nom_repository.go existente
- DTOs existentes como referencia