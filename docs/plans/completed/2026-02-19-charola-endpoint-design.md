# Diseño: Endpoint de Cálculo de Charolas

## Fecha: 2026-02-19

## Objetivo

Crear un endpoint API para calcular el tamaño de charolas (bandejas portacables) usando los servicios de dominio existentes.

## Contexto

El proyecto ya cuenta con:
- **Domain services**: `CalcularCharolaEspaciado` y `CalcularCharolaTriangular`
- **Value objects**: `ConductorCharola`, `CableControl`, `EntradaTablaCanalizacion`
- **Tablas NOM**: Tablas de charolas ya existentes en el repositorio CSV

## Arquitectura

### Flujo de Datos

```
HTTP Request (JSON)
       ↓
    Handler (infrastructure)
       ↓ parsea JSON
    DTO Input (primitivos)
       ↓
    Use Case (application)
       ↓ convierte DTO → VOs
    Domain Service (servicios existentes)
       ↓ retorna entity.Canalizacion
    Use Case
       ↓ convierte Entity → DTO
    DTO Output
       ↓
    Handler
       ↓
HTTP Response (JSON)
```

### Componentes a Crear

#### 1. DTOs (`internal/calculos/application/dto/`)

**CharolaInput**:
```go
type CharolaInput struct {
    TipoCharola        string  // "ESPACIADO" o "TRIANGULAR"
    HilosPorFase       int
    SistemaElectrico   string  // "MONOFASICO", "BIFASICO", "DELTA", "ESTRELLA"
    DiametroFaseMM     float64
    DiametroTierraMM   float64
    DiametroControlMM  *float64 // opcional
}
```

**CharolaOutput**:
```go
type CharolaOutput struct {
    Tipo            string  // "CHAROLA_CABLE_ESPACIADO" o "CHAROLA_CABLE_TRIANGULAR"
    Tamano          string  // ej: "300mm"
    TamanoPulgadas  string  // ej: "12"
    AnchoRequerido  float64 // mm
}
```

#### 2. Use Cases (`internal/calculos/application/usecase/`)

- `CalcularCharolaEspaciadoUseCase` — usa `service.CalcularCharolaEspaciado`
- `CalcularCharolaTriangularUseCase` — usa `service.CalcularCharolaTriangular`

Cada use case:
1. Valida el DTO de entrada
2. Convierte primitivos → value objects (`ConductorCharola`)
3. Llama al servicio de dominio
4. Convierte resultado → DTO de salida
5. Retorna

#### 3. Handler HTTP (`internal/calculos/infrastructure/adapter/driver/http/`)

- `charola_handler.go` — nuevo archivo con endpoints

**Endpoints**:
```
POST /api/v1/calculos/charola/espaciado
POST /api/v1/calculos/charola/triangular
```

#### 4. Router (`internal/calculos/infrastructure/router.go`)

- Agregar rutas de charola

#### 5. Wiring (`cmd/api/main.go`)

- Registrar handler de charola

## Decisiones de Diseño

1. **Separación de endpoints**: Dos endpoints distintos (espaciado vs triangular) en lugar de un endpoint con parámetro — más claro y seguír fumes de API REST.

2. **DTOs con primitivos**: Siguiendo las reglas del proyecto, los DTOs usan solo tipos primitivos (`string`, `int`, `float64`).

3. **Use cases separados**: Cada tipo de charola tiene su propio use case — sigue el principio de Single Responsibility.

4. **Ports existentes**: Se reutiliza `TablaNOMRepository` existente para obtener las tablas de charolas.

## Manejo de Errores

| Error | Código HTTP | Mensaje |
|-------|-------------|----------|
| Validación DTO | 400 | "datos de entrada inválidos" |
| Tabla no encontrada | 500 | "error interno" |
| Charola no suficiente | 400 | "no se encontró charola suficiente para el ancho requerido" |

## Pruebas

- Tests unitarios de use cases
- Tests de integración del endpoint
- Pruebas manuales con curl

### Ejemplo de Request/Response

**Request**:
```json
{
  "hilos_por_fase": 1,
  "sistema_electrico": "DELTA",
  "diametro_fase_mm": 25.48,
  "diametro_tierra_mm": 8.5
}
```

**Response**:
```json
{
  "tipo": "CHAROLA_CABLE_ESPACIADO",
  "tamano": "300mm",
  "tamano_pulgadas": "12",
  "ancho_requerido": 162.38
}
```

## Task ID

- [ ] Crear DTOs
- [ ] Crear use cases
- [ ] Crear handler HTTP
- [ ] Actualizar router
- [ ] Actualizar main.go (wiring)
- [ ] Tests
- [ ] Pruebas manuales
