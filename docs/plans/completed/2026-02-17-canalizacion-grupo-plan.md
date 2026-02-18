# Plan: Endpoint Canalización de Grupo de Alimentadores

**Fecha:** 2026-02-17  
**Diseño:** `docs/plans/2026-02-17-canalizacion-grupo-design.md`  
**Estado:** Listo para ejecutar

---

## Contexto

Implementar endpoint `POST /api/v1/calculos/canalizacion-grupo` que calcula la tubería para un grupo de alimentadores con número de tubos configurable. Reutiliza el servicio de dominio existente `CalcularCanalizacion`.

---

## Tareas por Fase

### Fase 1: Domain (Modificación)

#### Tarea 1: Agregar campo FactorRelleno a entidad Canalizacion

**Archivo:** `internal/calculos/domain/entity/canalizacion.go`  
**Acción:** Agregar campo `FactorRelleno float64` al struct `Canalizacion`  
**Cambios:**
- Agregar campo en struct
- Agregar parámetro en constructor `NewCanalizacion`
- No requiere validación adicional (el servicio garantiza valores válidos)

**Verificación:** `go build ./internal/calculos/domain/...`

---

#### Tarea 2: Actualizar tests de Canalizacion

**Archivo:** `internal/calculos/domain/entity/canalizacion_test.go`  
**Acción:** Actualizar tests existentes para incluir FactorRelleno  
**Verificación:** `go test ./internal/calculos/domain/entity/...`

---

#### Tarea 3: Modificar servicio CalcularCanalizacion para asignar FactorRelleno

**Archivo:** `internal/calculos/domain/service/calculo_canalizacion.go`  
**Acción:** Asignar `FactorRelleno` al construir resultado  
**Cambios:**
- En el return exitoso, agregar `FactorRelleno: factorRelleno`
- La variable `factorRelleno` ya existe (línea 68)

**Verificación:** `go test ./internal/calculos/domain/service/...`

---

#### Tarea 4: Actualizar tests del servicio CalcularCanalizacion

**Archivo:** `internal/calculos/domain/service/calculo_canalizacion_test.go`  
**Acción:** Verificar que FactorRelleno se asigna correctamente  
**Verificación:** `go test ./internal/calculos/domain/service/...`

---

### Fase 2: Application (DTOs y Use Case)

#### Tarea 5: Crear DTOs para CanalizacionGrupo

**Archivo:** `internal/calculos/application/dto/canalizacion_grupo.go` (NUEVO)  
**Acción:** Crear structs:
- `ConductorGrupoInput` con campos `Cantidad int`, `SeccionMM2 float64`
- `CanalizacionGrupoInput` con campos `Conductores []ConductorGrupoInput`, `SeccionTierraMM2 float64`, `TipoCanalizacion string`, `NumeroDeTubos int`
- `CanalizacionGrupoOutput` con campos `Tamano string`, `AreaTotalMM2 float64`, `AreaPorTuboMM2 float64`, `NumeroDeTubos int`, `FactorRelleno float64`
- Método `Validate()` en input

**Verificación:** `go build ./internal/calculos/application/dto/...`

---

#### Tarea 6: Crear CalcularCanalizacionGrupoUseCase

**Archivo:** `internal/calculos/application/usecase/calcular_canalizacion_grupo.go` (NUEVO)  
**Acción:** Crear use case con:
- Struct con dependencia `tablaRepo port.TablaNOMRepository`
- Constructor `NewCalcularCanalizacionGrupoUseCase`
- Método `Execute(ctx, input) (output, error)`:
  1. Validar input
  2. Convertir conductores DTO → `[]service.ConductorParaCanalizacion`
  3. Agregar tierras: `NumeroDeTubos` conductores de tierra
  4. Parsear tipo canalización
  5. Obtener tabla NOM
  6. Llamar `service.CalcularCanalizacion`
  7. Mapear resultado a DTO output

**Verificación:** `go build ./internal/calculos/application/usecase/...`

---

#### Tarea 7: Crear tests para CalcularCanalizacionGrupoUseCase

**Archivo:** `internal/calculos/application/usecase/calcular_canalizacion_grupo_test.go` (NUEVO)  
**Acción:** Tests unitarios con mock de TablaNOMRepository:
- Test caso exitoso con 1 tubo
- Test caso exitoso con 2 tubos (verifica 2 tierras)
- Test error validación input
- Test error tabla no disponible

**Verificación:** `go test ./internal/calculos/application/usecase/...`

---

### Fase 3: Infrastructure (Handler y Wiring)

#### Tarea 8: Agregar use case al CalculoHandler

**Archivo:** `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`  
**Acción:**
- Agregar campo `calcularCanalizacionGrupoUC *usecase.CalcularCanalizacionGrupoUseCase` al struct
- Agregar parámetro al constructor `NewCalculoHandler`

**Verificación:** `go build ./internal/calculos/infrastructure/...`

---

#### Tarea 9: Crear método handler CalcularCanalizacionGrupo

**Archivo:** `internal/calculos/infrastructure/adapter/driver/http/calculo_handler.go`  
**Acción:** Agregar método:
```go
func (h *CalculoHandler) CalcularCanalizacionGrupo(c *gin.Context) {
    // 1. Bind JSON a CanalizacionGrupoInput
    // 2. Ejecutar use case
    // 3. Mapear errores
    // 4. Retornar JSON
}
```
También agregar structs de request/response y función de mapeo de errores.

**Verificación:** `go build ./internal/calculos/infrastructure/...`

---

#### Tarea 10: Registrar ruta en router

**Archivo:** `internal/calculos/infrastructure/router.go`  
**Acción:**
- Agregar parámetro `calcularCanalizacionGrupoUC *usecase.CalcularCanalizacionGrupoUseCase` a `NewRouter`
- Pasar al constructor de `CalculoHandler`
- Agregar línea: `calculos.POST("/canalizacion-grupo", calculoHandler.CalcularCanalizacionGrupo)`

**Verificación:** `go build ./internal/calculos/infrastructure/...`

---

#### Tarea 11: Actualizar wiring en main.go

**Archivo:** `cmd/api/main.go`  
**Acción:**
- Crear instancia: `calcularCanalizacionGrupoUC := usecase.NewCalcularCanalizacionGrupoUseCase(tablaRepo)`
- Pasar al router: `infrastructure.NewRouter(..., calcularCanalizacionGrupoUC)`

**Verificación:** `go build ./cmd/api/...`

---

#### Tarea 12: Crear test de integración para el endpoint

**Archivo:** `internal/calculos/infrastructure/adapter/driver/http/calculo_handler_test.go`  
**Acción:** Agregar tests:
- Test POST /canalizacion-grupo exitoso
- Test validación de input
- Test con 2 tubos (verifica factor relleno y área)

**Verificación:** `go test ./internal/calculos/infrastructure/...`

---

### Fase 4: Actualizar código existente afectado

#### Tarea 13: Actualizar DimensionarCanalizacionUseCase para usar FactorRelleno

**Archivo:** `internal/calculos/application/usecase/dimensionar_canalizacion.go`  
**Acción:** Actualizar mapeo a DTO para incluir `FactorRelleno` del resultado del servicio

**Verificación:** `go test ./internal/calculos/application/usecase/...`

---

#### Tarea 14: Actualizar ResultadoCanalizacion DTO

**Archivo:** `internal/calculos/application/dto/memoria_output.go`  
**Acción:** Agregar campo `FactorRelleno float64` a `ResultadoCanalizacion` si no existe

**Verificación:** `go build ./internal/calculos/application/dto/...`

---

## Verificación Final

- [ ] `go build ./...` — compila sin errores
- [ ] `go test ./...` — todos los tests pasan
- [ ] `go vet ./...` — sin warnings
- [ ] Probar endpoint manualmente:
  ```bash
  curl -X POST http://localhost:8080/api/v1/calculos/canalizacion-grupo \
    -H "Content-Type: application/json" \
    -d '{
      "conductores": [{"cantidad": 6, "seccion_mm2": 53.49}],
      "seccion_tierra_mm2": 21.15,
      "tipo_canalizacion": "TUBERIA_PVC",
      "numero_de_tubos": 2
    }'
  ```

---

## Orden de Ejecución

```
Fase 1 (Domain)     →  Tareas 1-4
Fase 2 (Application) →  Tareas 5-7
Fase 3 (Infrastructure) →  Tareas 8-12
Fase 4 (Actualizar existente) →  Tareas 13-14
Verificación Final
```

---

## Agentes Responsables

| Fase | Agente | Tareas |
|------|--------|--------|
| Domain | `domain-agent` | 1, 2, 3, 4 |
| Application | `application-agent` | 5, 6, 7, 13, 14 |
| Infrastructure | `infrastructure-agent` | 8, 9, 10, 11, 12 |
| Wiring | Orquestador | Verificación final |

---

## Riesgos y Mitigación

| Riesgo | Mitigación |
|--------|------------|
| Cambio en entidad Canalizacion rompe código existente | Tareas 13-14 actualizan código afectado |
| Tests existentes fallan por nuevo campo | Tarea 2 y 4 actualizan tests |
