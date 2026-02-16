# Diseño: API Corriente Ajustada

**Fecha:** 2026-02-15  
**Feature:** calculos  
**Endpoint:** POST /api/v1/calculos/corriente-ajustada

## Resumen

Endpoint dedicado para calcular la corriente ajustada aplicando factores de corrección según NOM:
- Factor de temperatura (tabla NOM 310-15(b)(2)(a))
- Factor de agrupamiento (distribución en tuberías)
- Factor de uso (según tipo de equipo: filtro 1.35, otros 1.25)

## Alcance

**SÍ incluye:**
- Cálculo de corriente ajustada con factores NOM
- Soporte para tubería y charola (espaciada/triangular)
- Distribución de conductores en múltiples tubos
- Factor de uso por tipo de equipo

**NO incluye:**
- Selección de calibre de conductor
- Dimensionamiento de canalización (diámetros)
- Cálculo de caída de tensión

## Contrato API

### Request

```json
{
  "corriente_nominal": 50.0,
  "estado": "Sonora",
  "tipo_canalizacion": "TUBERIA_PVC",
  "sistema_electrico": "DELTA",
  "tipo_equipo": "FILTRO_ACTIVO",
  "hilos_por_fase": 2,
  "num_tuberias": 2
}
```

| Campo | Tipo | Obligatorio | Descripción |
|-------|------|-------------|-------------|
| `corriente_nominal` | float64 | Sí | Corriente nominal del circuito (A) |
| `estado` | string | Sí | Estado mexicano para temperatura ambiente |
| `tipo_canalizacion` | string | Sí | `TUBERIA_PVC`, `TUBERIA_METALICA`, `CHAROLA_ESPACIADA`, `CHAROLA_TRIANGULAR` |
| `sistema_electrico` | string | Sí | `MONOFASICO`, `ESTRELLA`, `DELTA` |
| `tipo_equipo` | string | Sí | `FILTRO_ACTIVO`, `FILTRO_RECHAZO`, `TRANSFORMADOR`, `CARGA` |
| `hilos_por_fase` | int | No | Default: 1. Conductores en paralelo por fase |
| `num_tuberias` | int | No | Default: 1. Solo aplica para tubería |

### Response (200 OK)

```json
{
  "success": true,
  "data": {
    "corriente_nominal": 50.0,
    "corriente_ajustada": 91.406,
    "factor_total": 1.828,
    "factor_temperatura": 0.91,
    "factor_agrupamiento": 0.8,
    "factor_uso": 1.35,
    "conductores_por_tubo": 3,
    "cantidad_conductores_total": 6,
    "temperatura_conductor": 75,
    "temperatura_ambiente": 40
  }
}
```

### Errores

| HTTP | Código | Escenario |
|------|--------|-----------|
| 400 | `VALIDATION_ERROR` | Campos inválidos o faltantes |
| 400 | `TIPO_CANALIZACION_INVALIDO` | Canalización no soportada |
| 400 | `SISTEMA_ELECTRICO_INVALIDO` | Sistema eléctrico inválido |
| 400 | `TIPO_EQUIPO_INVALIDO` | Tipo de equipo no soportado |
| 422 | `FACTOR_NO_ENCONTRADO` | No se encontró factor en tablas NOM |
| 500 | `INTERNAL_ERROR` | Error interno del servidor |

## Lógica de Cálculo

### 1. Factor de Temperatura
```
1. Obtener temperatura ambiente por estado (repositorio)
2. Seleccionar temperatura de conductor (60°C, 75°C, 90°C) según canalización y corriente
3. Buscar factor en tabla NOM 310-15(b)(2)(a)
```

### 2. Factor de Agrupamiento
```
Para TUBERIA:
  fases = sistemaElectrico.CantidadConductores()  // 2, 3 o 4
  cantidad_total = fases × hilos_por_fase
  conductores_por_tubo = cantidad_total / num_tuberias
  factor_agrupamiento = lookup_tabla_NOM(conductores_por_tubo)

Para CHAROLA_ESPACIADA:
  factor_agrupamiento = 1.0  // No aplica

Para CHAROLA_TRIANGULAR:
  // Usar cálculo específico de charola triangular
  factor_agrupamiento = calcularFactorCharolaTriangular(sistema)
```

### 3. Factor de Uso
```
SWITCH tipo_equipo:
  FILTRO_ACTIVO, FILTRO_RECHAZO → 1.35
  TRANSFORMADOR, CARGA → 1.25
```

### 4. Cálculo Final
```
factor_total = factor_temperatura × factor_agrupamiento × factor_uso
corriente_ajustada = corriente_nominal × factor_total
```

## Arquitectura

### Capa Domain (domain-agent)

**Archivo:** `internal/calculos/domain/service/calcular_factor_uso.go`

```go
// CalcularFactorUso retorna el factor de uso según tipo de equipo
// FILTRO_* → 1.35
// TRANSFORMADOR, CARGA → 1.25
func CalcularFactorUso(tipoEquipo entity.TipoEquipo) (float64, error)
```

**Tests:** `calcular_factor_uso_test.go`

### Capa Application (application-agent)

**Modificar:** `internal/calculos/application/dto/resultado_ajuste.go` (nuevo archivo)

```go
// Extender ResultadoAjusteCorriente
type ResultadoAjusteCorriente struct {
    CorrienteAjustada        float64
    FactorTemperatura        float64
    FactorAgrupamiento       float64
    FactorUso                float64          // NUEVO
    FactorTotal              float64
    Temperatura              int
    ConductoresPorTubo       int              // NUEVO
    CantidadConductoresTotal int              // NUEVO
    TemperaturaAmbiente      int              // NUEVO
}
```

**Modificar:** `internal/calculos/application/usecase/ajustar_corriente.go`

```go
// Extender Execute para aceptar nuevos parámetros
func (uc *AjustarCorrienteUseCase) Execute(
    ctx context.Context,
    corrienteNominal valueobject.Corriente,
    estado string,
    tipoCanalizacion entity.TipoCanalizacion,
    sistemaElectrico entity.SistemaElectrico,
    tipoEquipo entity.TipoEquipo,        // NUEVO
    hilosPorFase int,                     // NUEVO
    numTuberias int,                      // NUEVO
) (dto.ResultadoAjusteCorriente, error)
```

**Nota:** Eliminar parámetro `temperaturaOverride` ya que no se usa en este endpoint.

### Capa Infrastructure (infrastructure-agent)

**Crear:** `internal/calculos/infrastructure/adapter/driver/http/corriente_ajustada_handler.go` (o extender calculo_handler.go)

```go
// CorrienteAjustadaRequest body del POST
type CorrienteAjustadaRequest struct {
    CorrienteNominal  float64 `json:"corriente_nominal" binding:"required,gt=0"`
    Estado            string  `json:"estado" binding:"required"`
    TipoCanalizacion  string  `json:"tipo_canalizacion" binding:"required"`
    SistemaElectrico  string  `json:"sistema_electrico" binding:"required"`
    TipoEquipo        string  `json:"tipo_equipo" binding:"required"`
    HilosPorFase      int     `json:"hilos_por_fase" binding:"gte=1"`
    NumTuberias       int     `json:"num_tuberias" binding:"gte=1"`
}

// CorrienteAjustadaResponse respuesta exitosa
type CorrienteAjustadaResponse struct {
    Success bool                         `json:"success"`
    Data    dto.ResultadoAjusteCorriente `json:"data"`
}

// CalcularCorrienteAjustada handler HTTP
func (h *CalculoHandler) CalcularCorrienteAjustada(c *gin.Context)
```

**Modificar:** `internal/calculos/infrastructure/router.go`

```go
// Agregar ruta
calculos.POST("/corriente-ajustada", calculoHandler.CalcularCorrienteAjustada)
```

## Dependencias Entre Agentes

```
domain-agent (primero)
  ↓ crea CalcualarFactorUso
  
aplication-agent (segundo)
  ↓ usa CalcualarFactorUso
  ↓ extiende AjustarCorrienteUseCase
  
infrastructure-agent (tercero)
  ↓ usa AjustarCorrienteUseCase extendido
  ↓ crea handler y ruta
```

## Notas de Implementación

1. **Validaciones:**
   - `hilos_por_fase` ≥ 1
   - `num_tuberias` ≥ 1
   - (fases × hilos_por_fase) debe ser divisible por `num_tuberias`

2. **Charola:**
   - Espaciada: factor_agrupamiento siempre 1.0
   - Triangular: usar servicio existente `calcular_charola_triangular.go`

3. **Reutilización:**
   - Usar `entity.TipoEquipo` existente
   - Usar `AjustarCorriente` de `service/ajuste_corriente.go`
   - Usar tablas NOM ya cargadas en repositorio

## QA Checklist

- [ ] Tests unitarios para `CalcularFactorUso`
- [ ] Tests unitarios para `AjustarCorrienteUseCase` extendido
- [ ] Tests de integración para endpoint HTTP
- [ ] Validación de errores (400, 422, 500)
- [ ] Documentación de campos en request/response
- [ ] `go test ./...` pasa
