# Tasks: Número de Hilos de Tierra

## Overview

Este cambio elimina el valor hardcodeado de `1` hilo de tierra en el cálculo de dimensionamiento de canalización, implementando cálculo automático según las reglas de la normativa NOM:
- Charola → 1 hilo
- Tubería con ≤2 tubes → 1 hilo
- Tubería con >2 tubes → 2 hilos

## Phase 1: Foundation (DTO Changes)

- [x] 1.1 Agregar campo `NumTierras int` con tag `json:"-"` al struct `TuberiaInput` en `internal/calculos/application/dto/tuberia_input.go`
  - Ubicación: después del campo `NumTuberias` (línea 19)
  - El campo no se expone al usuario, solo uso interno

- [x] 1.2 Agregar validación para `NumTierras` en el método `Validate()` de `TuberiaInput`
  - Validar que sea mayor a 0
  - Si es 0 o negativo, usar default de 1

## Phase 2: Core Implementation (Helper Function)

- [x] 2.1 Crear función helper `calcularNumHilosTierra()` en `internal/calculos/application/usecase/orquestador_memoria_calculo.go`
  -签名: `func calcularNumHilosTierra(tipoCanalizacion entity.TipoCanalizacion, numTuberias int) int`
  - Reglas:
    - Si tipo es Charola → return 1
    - Si numTuberias ≤ 2 → return 1
    - Si numTuberias > 2 → return 2
    - Si numTuberias ≤ 0 → default a 1 antes de calcular
  - Ubicación: antes de la función `Execute()` del orquestador

- [x] 2.2 Agregar documentación GoDoc a la función helper
  - Descripción de la regla de negocio según NOM
  - Ejemplos de uso

## Phase 3: Integration (Use Cases Modification)

- [x] 3.1 Modificar el orquestador `orquestador_memoria_calculo.go` para usar la función helper
  - En la sección donde se construye `tuberiaInput` (alrededor de línea 287)
  - Llamar `calcularNumHilosTierra(tipoCanalizacion, input.NumTuberias)` antes de crear el struct
  - Asignar el resultado al campo `NumTierras` del `TuberiaInput`

- [x] 3.2 Modificar `calcular_tamanio_tuberia.go` para usar `input.NumTierras` en lugar del hardcode
  - En la línea 76, cambiar `1, // 1 tierra por tubo` por `input.NumTierras,`
  - Verificar que el valor se pase correctamente al servicio de dominio

- [x] 3.3 Verificar que no haya otros lugares con hardcode de hilos de tierra
  - Buscar en `internal/calculos/application/usecase/` patrones como `1, // tierra` o similar
  - Resultado: No se encontraron otros hardcodes

## Phase 4: Testing (Unit Tests)

- [x] 4.1 Crear archivo de tests para la función helper `calcularNumHilosTierra`
  - Ubicación: `internal/calculos/application/usecase/calcular_num_hilos_tierra_test.go`

- [x] 4.2 Escribir unit tests cubriendo los siguientes escenarios:
  - Charola cable espaciado → 1 hilo
  - Charola cable triangular → 1 hilo
  - Tubería PVC con 1 tubo → 1 hilo
  - Tubería PVC con 2 tubos → 1 hilo
  - Tubería PVC con 3 tubos → 2 hilos
  - Tubería PVC con 4+ tubos → 2 hilos
  - Tubería PVC con 100 tubos → 2 hilos
  - Tubería PVC/Aluminio/AceroPG/AceroPD con valores límite probados
  - Tubería con 0 tubos → default 1 hilo
  - Tubería con valor negativo → default 1 hilo

- [x] 4.3 Ejecutar tests existentes para verificar compatibilidad
  - Tests de usecase: ✅ PASS (16 subtests)
  - Tests de domain/service: ✅ PASS
  - Tests de domain/entity: ✅ PASS
  - NOTA: Test preexistente `TestEquipoInput_JSONTensionUnidad` falla pero no está relacionado con este cambio

## Implementation Order

1. **Phase 1** 首先执行，因为 DTO 字段是其他任务的基础依赖
2. **Phase 2** 在 Phase 1 之后执行，因为需要在管道输入中使用 helper 函数
3. **Phase 3** 在 Phase 1 和 2 之后执行，因为需要同时修改 DTO 和使用 helper 函数
4. **Phase 4** 在所有实现完成后执行

## Dependencies

- `entity.TipoCanalizacion.EsCharola()` 已存在于 `internal/calculos/domain/entity/tipo_canalizacion.go`
- `service.CalcularTamanioTuberiaWithMultiplePipes` 已支持 `tierras` 参数（第3个参数）
- `input.NumTuberias` 在 `dto.TuberiaInput` 中已存在

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| 修改 DTO 后与现有 API 端点不兼容 | Low | Medium | `json:"-"` 标签确保字段不暴露给用户 |
| 计算逻辑错误 | Low | High | Phase 4 的单元测试覆盖所有场景 |
| 管道计算 use case 中存在其他硬编码值 | Medium | Medium | 在 Phase 3.3 中进行代码搜索检查 |
