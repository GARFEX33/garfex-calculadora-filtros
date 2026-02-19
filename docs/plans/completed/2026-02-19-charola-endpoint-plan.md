# Plan: Endpoint de Cálculo de Charolas

## Fecha: 2026-02-19

## Objetivo

Crear endpoint API para cálculo de charolas (espaciado y triangular).

## Proyecto

- **Rama**: `feature/charola-endpoint`
- **Módulo Go**: `github.com/garfex/calculadora-filtros`

---

## Pasos de Implementación

### Paso 1: Application Layer - DTOs

**Archivos a crear:**

1. `internal/calculos/application/dto/charola_espaciado.go`
2. `internal/calculos/application/dto/charola_triangular.go`

**Contenido:**
- `CharolaEspaciadoInput` struct con campos: `HilosPorFase`, `SistemaElectrico`, `DiametroFaseMM`, `DiametroTierraMM`, `DiametroControlMM` (opcional)
- `CharolaEspaciadoOutput` struct con campos: `Tipo`, `Tamano`, `TamanoPulgadas`, `AnchoRequerido`
- Métodos `Validate()` en cada input

---

### Paso 2: Application Layer - Use Cases

**Archivos a crear:**

1. `internal/calculos/application/usecase/calcular_charola_espaciado.go`
2. `internal/calculos/application/usecase/calcular_charola_triangular.go`

**Contenido:**
- `CalcularCharolaEspaciadoUseCase` struct con campo `repo port.TablaNOMRepository`
- `Execute()` que:
  1. Valida input
  2. Crea `ConductorCharola` para fase y tierra
  3. Crea `CableControl` si se proporciona
  4. Obtiene tabla de charolas del repo
  5. Llama `service.CalcularCharolaEspaciado`
  6. Convierte resultado a DTO output
  7. Retorna

- Mismo patrón para `CalcularCharolaTriangularUseCase`

---

### Paso 3: Application Layer - Tests

**Archivos a crear:**

1. `internal/calculos/application/usecase/calcular_charola_espaciado_test.go`
2. `internal/calculos/application/usecase/calcular_charola_triangular_test.go`

**Contenido:**
- Tests con mock de `TablaNOMRepository`
- Casos: happy path, errores de validación, tabla vacía

---

### Paso 4: Infrastructure Layer - Handler

**Archivo a crear:**

`internal/calculos/infrastructure/adapter/driver/http/charola_handler.go`

**Contenido:**
- `CharolaHandler` struct con campos `calcularEspaciadoUseCase`, `calcularTriangularUseCase`
- `PostCharolaEspaciado(c *gin.Context)` - endpoint `/charola/espaciado`
- `PostCharolaTriangular(c *gin.Context)` - endpoint `/charola/triangular`
- Manejo de errores con `c.JSON(400/500, ...)`

---

### Paso 5: Infrastructure Layer - Router

**Archivo a modificar:**

`internal/calculos/infrastructure/router.go`

**Cambios:**
```go
charolaHandler := driver.NewCharolaHandler(...)
api.POST("/charola/espaciado", charolaHandler.PostCharolaEspaciado)
api.POST("/charola/triangular", charolaHandler.PostCharolaTriangular)
```

---

### Paso 6: Wiring - main.go

**Archivo a modificar:**

`cmd/api/main.go`

**Cambios:**
- Crear instancias de use cases
- Crear handler y registrar en router

---

### Paso 7: Verificación

```bash
go build ./...
go test ./...
```

---

### Paso 8: Pruebas Manuales

```bash
# Iniciar servidor
go run cmd/api/main.go

# Probar endpoint espaciado
curl -X POST http://localhost:8080/api/v1/calculos/charola/espaciado \
  -H "Content-Type: application/json" \
  -d '{"hilos_por_fase":1,"sistema_electrico":"DELTA","diametro_fase_mm":25.48,"diametro_tierra_mm":8.5}'

# Probar endpoint triangular
curl -X POST http://localhost:8080/api/v1/calculos/charola/triangular \
  -H "Content-Type: application/json" \
  -d '{"hilos_por_fase":2,"diametro_fase_mm":25.48,"diametro_tierra_mm":7.42}'
```

---

## Dependencias Entre Pasos

```
Paso 1 (DTOs) → Paso 2 (Use Cases) → Paso 3 (Tests)
                                            ↓
                  Paso 6 (main.go) ← Paso 5 (Router) ← Paso 4 (Handler)
                                            ↓
                                    Paso 7 (Verificación)
                                            ↓
                                    Paso 8 (Pruebas Manuales)
```

---

## Responsables

- **Domain**: No requiere cambios (ya existe)
- **Application**: `application-agent`
- **Infrastructure**: `infrastructure-agent`
- **Wiring**: Orquestador

---

## Completion Criteria

- [ ] go build ./... pasa
- [ ] go test ./... pasa
- [ ] Endpoints responden correctamente
- [ ] AGENTS.md actualizado si es necesario
