# Plan: Implementación Endpoint Amperaje

## Objetivo
Agregar endpoint `POST /api/v1/calculos/amperaje` que calcula amperaje nominal desde potencia.

## Tareas

### Tarea 1: Agregar structs de request/response al handler
- Ubicación: `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`
- Agregar struct `CalcularAmperajeRequest` con campos: potencia_watts, tension, tipo_carga, sistema_electrico, factor_potencia
- Agregar struct `CalcularAmperajeResponse` con campos: amperaje, unidad

### Tarea 2: Agregar método CalcularAmperaje al handler
- Ubicación: `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`
- Agregar método `CalcularAmperaje(c *gin.Context)` al `CalculoHandler`
- Incluir binding JSON con tags de validación
- Llamar al use case `CalcularAmperajeNominalUseCase`
- Mapear respuesta exitosa

### Tarea 3: Agregar método de mapeo de errores
- Ubicación: `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`
- Agregar método `mapAmperajeErrorToResponse(err error)` que retorne (int, error struct)

### Tarea 4: Modificar constructor del handler
- Ubicación: `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`
- Agregar campo `calcularAmperajeUC *usecase.CalcularAmperajeNominalUseCase` al struct
- Actualizar constructor `NewCalculoHandler`

### Tarea 5: Actualizar router
- Ubicación: `internal/calculos/infrastructure/router.go`
- Agregar parámetro `calcularAmperajeUC *usecase.CalcularAmperajeNominalUseCase` a NewRouter
- Agregar route: `calculos.POST("/amperaje", calculoHandler.CalcularAmperaje)`

### Tarea 6: Verificar con tests
- Comando: `go test ./internal/calculos/infrastructure/...`
- Verificar que compila y los tests pasan
