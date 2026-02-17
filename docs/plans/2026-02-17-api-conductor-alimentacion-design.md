# API Conductor Alimentacion - Design Document

**Fecha:** 2026-02-17  
**Estado:** Aprobado  
**Autor:** Orquestador

## Resumen

Crear endpoint independiente para seleccion de conductor de alimentacion, separando la funcionalidad del use case combinado existente. Esto permite:

1. API rapida para interfaces que solo necesitan seleccion de conductor
2. Verificar funcionalidad de forma aislada
3. Responsabilidades separadas (Single Responsibility)

## Contexto

### Estado Actual

| Componente | Estado |
|------------|--------|
| `SeleccionarConductorAlimentacion` (domain service) | Existe |
| `SeleccionarTemperatura` (domain service) | Existe (regla NOM <=100A->60C) |
| `SeleccionarConductorUseCase` (combinado) | Mezcla alimentacion + tierra |
| `SeleccionarConductorTierraUseCase` | Ya separado |
| API `/conductor-tierra` | Ya existe |
| API `/conductor-alimentacion` | NO existe |

### Problema

El use case `SeleccionarConductorUseCase` viola Single Responsibility al manejar tanto alimentacion como tierra. Ademas, no hay forma de obtener solo el conductor de alimentacion via API sin ejecutar toda la memoria de calculo.

## Solucion

### Enfoque Seleccionado

**Enfoque A:** Crear UseCase y Handler independientes para alimentacion, eliminar el use case combinado, y refactorizar el orquestador para usar los dos use cases separados.

### Archivos Afectados

| Archivo | Accion |
|---------|--------|
| `application/dto/conductor_alimentacion.go` | CREAR |
| `application/usecase/seleccionar_conductor_alimentacion.go` | CREAR |
| `infrastructure/adapter/driver/http/conductor_alimentacion_handler.go` | CREAR |
| `application/usecase/seleccionar_conductor.go` | ELIMINAR |
| `application/dto/memoria_output.go` | MODIFICAR (eliminar ResultadoConductores) |
| `application/usecase/orquestador_memoria.go` | MODIFICAR |
| `infrastructure/router.go` | MODIFICAR |
| `cmd/api/main.go` | MODIFICAR |
| Tests afectados | MODIFICAR |

## API Contract

### Endpoint

```
POST /api/v1/calculos/conductor-alimentacion
```

### Request

```json
{
  "corriente_ajustada": 85.5,
  "tipo_canalizacion": "TUBERIA_PVC",
  "material": "Cu",
  "temperatura": 75,
  "hilos_por_fase": 1
}
```

| Campo | Tipo | Obligatorio | Default | Descripcion |
|-------|------|-------------|---------|-------------|
| `corriente_ajustada` | float | Si | - | Amperaje ya ajustado por factores NOM |
| `tipo_canalizacion` | string | Si | - | TUBERIA_PVC, TUBERIA_EMT, CHAROLA_ESPACIADO, CHAROLA_TRIANGULAR |
| `material` | string | No | "Cu" | "Cu" o "Al" |
| `temperatura` | int | No | Regla NOM | 60, 75, o 90 grados C |
| `hilos_por_fase` | int | No | 1 | Divide corriente entre hilos |

### Response (exito)

```json
{
  "success": true,
  "data": {
    "calibre": "4 AWG",
    "material": "Cu",
    "seccion_mm2": 21.2,
    "tipo_aislamiento": "THHN/THHW",
    "capacidad_nominal": 85,
    "tabla_usada": "NOM-001 310-15(b)(16) Cu 75C"
  }
}
```

### Response (error)

```json
{
  "success": false,
  "error": "No se encontro conductor adecuado",
  "code": "CONDUCTOR_NO_ENCONTRADO",
  "details": "corriente por hilo 650.00 A excede maxima capacidad de tabla 615.00 A"
}
```

## Componentes

### DTO Input (primitivos)

```go
// application/dto/conductor_alimentacion.go

type ConductorAlimentacionInput struct {
    CorrienteAjustada float64  // Obligatorio, > 0
    TipoCanalizacion  string   // Obligatorio
    Material          string   // Opcional, default "Cu"
    Temperatura       *int     // Opcional, nil -> regla NOM
    HilosPorFase      int      // Opcional, default 1
}

func (i ConductorAlimentacionInput) Validate() error
func (i ConductorAlimentacionInput) ToDomainMaterial() valueobject.MaterialConductor
```

### DTO Output (primitivos)

```go
type ConductorAlimentacionOutput struct {
    Calibre          string  
    Material         string  
    SeccionMM2       float64 
    TipoAislamiento  string  
    CapacidadNominal float64
    TablaUsada       string  
}
```

### Use Case

```go
// application/usecase/seleccionar_conductor_alimentacion.go

type SeleccionarConductorAlimentacionUseCase struct {
    tablaRepo port.TablaNOMRepository
}

func (uc *SeleccionarConductorAlimentacionUseCase) Execute(
    ctx context.Context,
    input dto.ConductorAlimentacionInput,
) (dto.ConductorAlimentacionOutput, error) {
    // 1. Validar DTO
    // 2. Convertir primitivos -> value objects (Corriente, TipoCanalizacion, Material)
    // 3. Determinar temperatura (input o regla NOM via service.SeleccionarTemperatura)
    // 4. Obtener tabla ampacidad del repo
    // 5. Llamar service.SeleccionarConductorAlimentacion
    // 6. Obtener capacidad del conductor seleccionado
    // 7. Generar nombre de tabla usada (via helpers)
    // 8. Convertir domain -> DTO output
}
```

### Handler HTTP

```go
// infrastructure/adapter/driver/http/conductor_alimentacion_handler.go

type ConductorAlimentacionHandler struct {
    useCase *usecase.SeleccionarConductorAlimentacionUseCase
}

func (h *ConductorAlimentacionHandler) SeleccionarConductorAlimentacion(c *gin.Context)
```

## Refactorizacion del Orquestador

### Antes

```go
type OrquestadorMemoriaCalculo struct {
    seleccionarConductor *SeleccionarConductorUseCase  // combinado
}

// En Execute():
resultadoConductores, err := o.seleccionarConductor.Execute(...)
output.ConductorAlimentacion = resultadoConductores.Alimentacion
output.ConductorTierra = resultadoConductores.Tierra
```

### Despues

```go
type OrquestadorMemoriaCalculo struct {
    seleccionarConductorAlimentacion *SeleccionarConductorAlimentacionUseCase
    seleccionarConductorTierra       *SeleccionarConductorTierraUseCase
}

// En Execute():
resultadoAlimentacion, err := o.seleccionarConductorAlimentacion.Execute(ctx, inputAlim)
output.ConductorAlimentacion = mapToResultadoConductor(resultadoAlimentacion)

resultadoTierra, err := o.seleccionarConductorTierra.Execute(ctx, input.ITM, input.Material.String())
output.ConductorTierra = mapToResultadoConductorFromTierra(resultadoTierra)
```

## Mapeo de Errores HTTP

| Error Domain/Application | HTTP Status |
|--------------------------|-------------|
| ErrEquipoInputInvalido | 400 |
| ErrTipoCanalizacionInvalido | 400 |
| valueobject.ErrCorrienteInvalida | 400 |
| service.ErrConductorNoEncontrado | 422 |
| Error interno | 500 |

## Eliminaciones

- `application/usecase/seleccionar_conductor.go` - use case combinado
- `dto.ResultadoConductores` de `memoria_output.go` - ya no se usa

## Dependencias de Dominio (ya existen)

- `service.SeleccionarConductorAlimentacion` - selecciona conductor por ampacidad
- `service.SeleccionarTemperatura` - aplica regla NOM (<=100A->60C, >100A->75C)
- `valueobject.Corriente`, `valueobject.MaterialConductor`, `valueobject.Temperatura`
- `entity.TipoCanalizacion`

## Testing

- Unit tests para DTO validation
- Unit tests para use case con mock de TablaNOMRepository
- Integration test para handler HTTP

## Criterios de Aceptacion

1. Endpoint `/conductor-alimentacion` responde correctamente
2. Use case combinado eliminado
3. Orquestador usa los dos use cases separados
4. Todos los tests pasan
5. `go test ./...` verde
