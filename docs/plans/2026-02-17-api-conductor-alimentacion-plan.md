# API Conductor Alimentacion - Plan de Implementacion

**Fecha:** 2026-02-17  
**Design Doc:** `2026-02-17-api-conductor-alimentacion-design.md`  
**Modulo Go:** `github.com/garfex/calculadora-filtros`

## Resumen

Implementar endpoint independiente para seleccion de conductor de alimentacion, eliminar use case combinado, y refactorizar orquestador.

## Orden de Tareas

El orden es critico: primero crear lo nuevo, luego refactorizar, finalmente eliminar.

---

## Paso 1: Crear DTO conductor_alimentacion.go

**Archivo:** `internal/calculos/application/dto/conductor_alimentacion.go`

**Crear:**

```go
// internal/calculos/application/dto/conductor_alimentacion.go
package dto

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ConductorAlimentacionInput contiene los datos de entrada para seleccionar
// el conductor de alimentacion segun tablas NOM 310-15.
type ConductorAlimentacionInput struct {
	// CorrienteAjustada es el amperaje ya ajustado por factores NOM.
	CorrienteAjustada float64 `json:"corriente_ajustada"`
	// TipoCanalizacion es el tipo de canalizacion (TUBERIA_PVC, CHAROLA_ESPACIADO, etc).
	TipoCanalizacion string `json:"tipo_canalizacion"`
	// Material es el material del conductor ("Cu" o "Al").
	// Si esta vacio, se usa "Cu" por defecto.
	Material string `json:"material"`
	// Temperatura es la temperatura de operacion (60, 75, 90).
	// Si es nil, se aplica la regla NOM automatica.
	Temperatura *int `json:"temperatura"`
	// HilosPorFase es el numero de conductores por fase.
	// Si es 0 o menor, se usa 1 por defecto.
	HilosPorFase int `json:"hilos_por_fase"`
}

// Validate verifica que los campos requeridos sean validos.
func (i ConductorAlimentacionInput) Validate() error {
	if i.CorrienteAjustada <= 0 {
		return fmt.Errorf("%w: corriente_ajustada debe ser mayor que cero", ErrEquipoInputInvalido)
	}
	if i.TipoCanalizacion == "" {
		return fmt.Errorf("%w: tipo_canalizacion es requerido", ErrEquipoInputInvalido)
	}
	return nil
}

// ToDomainMaterial convierte el material del DTO a value object.
// Si esta vacio o es invalido, retorna cobre por defecto.
func (i ConductorAlimentacionInput) ToDomainMaterial() valueobject.MaterialConductor {
	if i.Material == "Al" {
		return valueobject.MaterialAluminio
	}
	return valueobject.MaterialCobre
}

// ConductorAlimentacionOutput contiene el resultado de seleccionar
// el conductor de alimentacion.
type ConductorAlimentacionOutput struct {
	// Calibre es el calibre del conductor seleccionado (ej: "4 AWG", "250 MCM").
	Calibre string `json:"calibre"`
	// Material es el material del conductor ("Cu" o "Al").
	Material string `json:"material"`
	// SeccionMM2 es la seccion transversal en mm2.
	SeccionMM2 float64 `json:"seccion_mm2"`
	// TipoAislamiento es el tipo de aislamiento (ej: "THHN/THHW").
	TipoAislamiento string `json:"tipo_aislamiento"`
	// CapacidadNominal es la ampacidad del conductor segun la tabla.
	CapacidadNominal float64 `json:"capacidad_nominal"`
	// TablaUsada es el nombre descriptivo de la tabla NOM utilizada.
	TablaUsada string `json:"tabla_usada"`
}
```

**Verificacion:**
```bash
go build ./internal/calculos/application/dto/...
```

---

## Paso 2: Crear Use Case seleccionar_conductor_alimentacion.go

**Archivo:** `internal/calculos/application/usecase/seleccionar_conductor_alimentacion.go`

**Crear:**

```go
// internal/calculos/application/usecase/seleccionar_conductor_alimentacion.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase/helpers"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// SeleccionarConductorAlimentacionUseCase ejecuta la seleccion de conductor
// de alimentacion segun tablas NOM 310-15.
type SeleccionarConductorAlimentacionUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewSeleccionarConductorAlimentacionUseCase crea una nueva instancia.
func NewSeleccionarConductorAlimentacionUseCase(
	tablaRepo port.TablaNOMRepository,
) *SeleccionarConductorAlimentacionUseCase {
	return &SeleccionarConductorAlimentacionUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute selecciona el conductor de alimentacion apropiado.
func (uc *SeleccionarConductorAlimentacionUseCase) Execute(
	ctx context.Context,
	input dto.ConductorAlimentacionInput,
) (dto.ConductorAlimentacionOutput, error) {
	// 1. Validar DTO
	if err := input.Validate(); err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("validar input: %w", err)
	}

	// 2. Convertir primitivos a value objects
	corrienteAjustada, err := valueobject.NewCorriente(input.CorrienteAjustada)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("corriente invalida: %w", err)
	}

	tipoCanalizacion, err := entity.ParseTipoCanalizacion(input.TipoCanalizacion)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("tipo canalizacion invalido: %w", err)
	}

	material := input.ToDomainMaterial()

	hilosPorFase := input.HilosPorFase
	if hilosPorFase < 1 {
		hilosPorFase = 1
	}

	// 3. Determinar temperatura (input o regla NOM)
	var temperatura valueobject.Temperatura
	if input.Temperatura != nil {
		temperatura = valueobject.Temperatura(*input.Temperatura)
	} else {
		temperatura = service.SeleccionarTemperatura(corrienteAjustada, tipoCanalizacion, nil)
	}

	// 4. Obtener tabla de ampacidad
	tablaAmpacidad, err := uc.tablaRepo.ObtenerTablaAmpacidad(ctx, tipoCanalizacion, material, temperatura)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("obtener tabla ampacidad: %w", err)
	}

	// 5. Llamar servicio de dominio
	conductor, err := service.SeleccionarConductorAlimentacion(corrienteAjustada, hilosPorFase, tablaAmpacidad)
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("seleccionar conductor: %w", err)
	}

	// 6. Obtener capacidad del conductor seleccionado
	capacidad, err := uc.tablaRepo.ObtenerCapacidadConductor(ctx, tipoCanalizacion, material, temperatura, conductor.Calibre())
	if err != nil {
		return dto.ConductorAlimentacionOutput{}, fmt.Errorf("obtener capacidad: %w", err)
	}

	// 7. Generar nombre de tabla usada
	tablaUsada := helpers.NombreTablaAmpacidad(string(tipoCanalizacion), material, temperatura)

	// 8. Retornar DTO output
	return dto.ConductorAlimentacionOutput{
		Calibre:          conductor.Calibre(),
		Material:         conductor.Material().String(),
		SeccionMM2:       conductor.SeccionMM2(),
		TipoAislamiento:  conductor.TipoAislamiento(),
		CapacidadNominal: capacidad,
		TablaUsada:       tablaUsada,
	}, nil
}
```

**Verificacion:**
```bash
go build ./internal/calculos/application/usecase/...
```

---

## Paso 3: Crear Handler conductor_alimentacion_handler.go

**Archivo:** `internal/calculos/infrastructure/adapter/driver/http/conductor_alimentacion_handler.go`

**Crear:**

```go
// internal/calculos/infrastructure/adapter/driver/http/conductor_alimentacion_handler.go
package http

import (
	"errors"
	"net/http"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/gin-gonic/gin"
)

// ConductorAlimentacionHandler maneja el endpoint de conductor de alimentacion.
type ConductorAlimentacionHandler struct {
	useCase *usecase.SeleccionarConductorAlimentacionUseCase
}

// NewConductorAlimentacionHandler crea un nuevo handler.
func NewConductorAlimentacionHandler(uc *usecase.SeleccionarConductorAlimentacionUseCase) *ConductorAlimentacionHandler {
	return &ConductorAlimentacionHandler{
		useCase: uc,
	}
}

// ConductorAlimentacionRequest representa el body de la peticion POST.
type ConductorAlimentacionRequest struct {
	CorrienteAjustada float64 `json:"corriente_ajustada" binding:"required,gt=0"`
	TipoCanalizacion  string  `json:"tipo_canalizacion" binding:"required"`
	Material          string  `json:"material"`
	Temperatura       *int    `json:"temperatura"`
	HilosPorFase      int     `json:"hilos_por_fase"`
}

// ConductorAlimentacionResponse representa la respuesta exitosa.
type ConductorAlimentacionResponse struct {
	Success bool                           `json:"success"`
	Data    dto.ConductorAlimentacionOutput `json:"data"`
}

// ConductorAlimentacionResponseError representa la respuesta de error.
type ConductorAlimentacionResponseError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SeleccionarConductorAlimentacion POST /api/v1/calculos/conductor-alimentacion
func (h *ConductorAlimentacionHandler) SeleccionarConductorAlimentacion(c *gin.Context) {
	var req ConductorAlimentacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Error de validacion",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Convertir request a DTO
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: req.CorrienteAjustada,
		TipoCanalizacion:  req.TipoCanalizacion,
		Material:          req.Material,
		Temperatura:       req.Temperatura,
		HilosPorFase:      req.HilosPorFase,
	}

	// Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		status, response := h.mapErrorToResponse(err)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, ConductorAlimentacionResponse{
		Success: true,
		Data:    output,
	})
}

// mapErrorToResponse mapea errores del dominio a respuestas HTTP.
func (h *ConductorAlimentacionHandler) mapErrorToResponse(err error) (int, ConductorAlimentacionResponseError) {
	// Errores 400 - Bad Request
	if errors.Is(err, dto.ErrEquipoInputInvalido) {
		return http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Datos de entrada invalidos",
			Code:    "INPUT_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, entity.ErrTipoCanalizacionInvalido) {
		return http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Tipo de canalizacion invalido",
			Code:    "TIPO_CANALIZACION_INVALIDO",
			Details: err.Error(),
		}
	}

	if errors.Is(err, valueobject.ErrCorrienteInvalida) {
		return http.StatusBadRequest, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "Corriente invalida",
			Code:    "CORRIENTE_INVALIDA",
			Details: err.Error(),
		}
	}

	// Errores 422 - Unprocessable Entity
	if errors.Is(err, service.ErrConductorNoEncontrado) {
		return http.StatusUnprocessableEntity, ConductorAlimentacionResponseError{
			Success: false,
			Error:   "No se encontro conductor adecuado",
			Code:    "CONDUCTOR_NO_ENCONTRADO",
			Details: err.Error(),
		}
	}

	// Por defecto: error interno 500
	return http.StatusInternalServerError, ConductorAlimentacionResponseError{
		Success: false,
		Error:   "Error interno del servidor",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	}
}
```

**Verificacion:**
```bash
go build ./internal/calculos/infrastructure/...
```

---

## Paso 4: Actualizar router.go

**Archivo:** `internal/calculos/infrastructure/router.go`

**Modificar:** Agregar parametro del nuevo use case y registrar ruta.

**Verificacion:**
```bash
go build ./internal/calculos/infrastructure/...
```

---

## Paso 5: Refactorizar orquestador_memoria.go

**Archivo:** `internal/calculos/application/usecase/orquestador_memoria.go`

**Cambios:**
1. Reemplazar `seleccionarConductor *SeleccionarConductorUseCase` por los dos use cases separados
2. Actualizar constructor
3. Actualizar Execute() para llamar a ambos use cases

**Verificacion:**
```bash
go build ./internal/calculos/application/usecase/...
```

---

## Paso 6: Actualizar cmd/api/main.go

**Archivo:** `cmd/api/main.go`

**Cambios:**
1. Crear instancia de SeleccionarConductorAlimentacionUseCase
2. Actualizar creacion del orquestador
3. Pasar nuevo use case al router

**Verificacion:**
```bash
go build ./cmd/api/...
```

---

## Paso 7: Eliminar use case combinado

**Archivo a eliminar:** `internal/calculos/application/usecase/seleccionar_conductor.go`

**Verificacion:**
```bash
go build ./...
```

---

## Paso 8: Limpiar ResultadoConductores de memoria_output.go

**Archivo:** `internal/calculos/application/dto/memoria_output.go`

**Cambio:** Eliminar struct `ResultadoConductores` si ya no se usa.

**Verificacion:**
```bash
go build ./...
```

---

## Paso 9: Actualizar tests

**Archivos afectados:**
- `internal/calculos/application/usecase/orquestador_memoria_test.go`
- `internal/calculos/infrastructure/adapter/driver/http/calculo_handler_test.go`
- `tests/integration/fase2_calculo_test.go`

**Verificacion:**
```bash
go test ./...
```

---

## Paso 10: Test de integracion del nuevo endpoint

**Verificacion manual:**
```bash
go run cmd/api/main.go &
curl -X POST http://localhost:8080/api/v1/calculos/conductor-alimentacion \
  -H "Content-Type: application/json" \
  -d '{"corriente_ajustada":85.5,"tipo_canalizacion":"TUBERIA_PVC","material":"Cu","hilos_por_fase":1}'
```

---

## Verificacion Final

```bash
go test ./...
go build ./...
go vet ./...
```

## Criterios de Exito

- [ ] Endpoint `/conductor-alimentacion` responde correctamente
- [ ] Use case combinado eliminado
- [ ] Orquestador usa los dos use cases separados  
- [ ] Todos los tests pasan
- [ ] `go test ./...` verde
- [ ] `go build ./...` verde
