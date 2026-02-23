# Tasks: Agregar unidades de voltaje (V/kV)

## Phase 1: Value Object Tension - Modificaciones Core

- [x] 1.1 Agregar tipo `UnidadTension` con constantes `V` y `kV` en `internal/shared/kernel/valueobject/tension.go`
- [x] 1.2 Agregar función `ParseUnidadTension(s string) (UnidadTension, error)` en `internal/shared/kernel/valueobject/tension.go`
- [x] 1.3 Agregar función `normalizarAVolts(valor float64, unidad UnidadTension) int` en `internal/shared/kernel/valueobject/tension.go`
- [x] 1.4 Modificar `NewTension` para aceptar `(valor float64, unidad string)` y normalizar internamente
- [x] 1.5 Agregar error `ErrUnidadTensionInvalida` para unidades no reconocidas

## Phase 2: DTO EquipoInput - Actualización

- [x] 2.1 Agregar campo `TensionUnidad string` a `EquipoInput` en `internal/calculos/application/dto/equipo_input.go`
- [x] 2.2 Modificar `ApplyDefaults()` para agregar default "V" si TensionUnidad está vacío
- [x] 2.3 Modificar `ToDomainTension()` para pasar la unidad al value object
- [x] 2.4 Actualizar validación en `Validate()` si es necesario

## Phase 3: Testing - Unit Tests

- [x] 3.1 Agregar test: `NewTension(480, "V")` retorna Tension con valor 480
- [x] 3.2 Agregar test: `NewTension(0.48, "kV")` retorna Tension con valor 480
- [x] 3.3 Agregar test: `NewTension(0.5, "kV")` retorna error (500V no válido NOM)
- [x] 3.4 Agregar test: `NewTension(230, "V")` retorna error (no válido NOM)
- [x] 3.5 Agregar test: `NewTension(480, "")` usa default "V" (compatibilidad)
- [x] 3.6 Agregar test: `ParseUnidadTension("kV")` y `ParseUnidadTension("KV")` funcionan
- [x] 3.7 Agregar test: `Unidad()` retorna la unidad original

## Phase 4: Verificación e Integración

- [x] 4.1 Ejecutar tests existentes para verificar que no hay regression: `go test ./...`
- [x] 4.2 Verificar que el endpoint memoria acepta el nuevo campo `tension_unidad: "kV"`
- [x] 4.3 Verificar que el endpoint memoria funciona igual que antes sin especificar `tension_unidad`

## Phase 5: Documentación (opcional)

- [x] 5.1 Agregar comentarios en el código sobre las nuevas funcionalidades
