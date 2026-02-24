# Proposal: Cálculo Automático de Número de Hilos de Tierra

## Intent

Eliminar el valor hardcodeado de `1` hilo de tierra en el cálculo de tuberías y charolas, implementando lógica automática según la normativa NOM. El usuario NO debe enviar este dato — el sistema debe calcularlo internamente basándose en el número de tubos.

**Problema actual**: En `calcular_tamanio_tuberia.go` (línea 76) y en el orquestador, el número de hilos de tierra está hardcodeado a `1`, lo cual no cumple con las reglas de dimensionamiento cuando hay más de 2 tubos en paralelo.

## Scope

### In Scope
- Implementar cálculo automático de número de hilos de tierra en la lógica de dominio
- Modificar `CalcularTamanioTuberiaUseCase` para calcular y pasar el valor correcto
- Modificar el orquestador `OrquestadorMemoriaCalculoUseCase` para calcular el valor antes de llamar al use case de tubería
- Asegurar que para charolas siempre se use 1 hilo de tierra
- Actualizar el DTO `TuberiaInput` para NO exponer el campo `NumTierras` (el cálculo es interno)

### Out of Scope
- Cambios en el endpoint HTTP de tubería independiente (solo se modifica el flujo de la memoria de cálculo)
- Modificación de tablas NOM o datos estáticos
- Tests automatizados (serán definidos en la fase de tareas)

## Approach

**Estrategia**: Calcular el número de hilos de tierra en la capa Application (orquestador/use case) antes de llamar al servicio de dominio.

1. **Regla de negocio a implementar**:
   - **Charola**: Siempre 1 hilo de tierra
   - **Tubería**: 
     - Si `numTuberias <= 2` → 1 hilo de tierra
     - Si `numTuberias > 2` → 2 hilos de tierra

2. **Flujo de implementación**:
   - En `OrquestadorMemoriaCalculoUseCase`, calcular `numTierras` antes de construir `TuberiaInput`
   - Pasar el valor calculado al use case de tubería
   - El DTO `TuberiaInput` NO expondrá este campo al usuario (cálculo interno)

3. **Archivos a modificar**:
   - `calcular_tamanio_tuberia.go`: Eliminar el hardcode y usar el parámetro
   - `orquestador_memoria_calculo.go`: Calcular `numTierras` según la regla y pasarlo

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/calculos/application/usecase/calcular_tamanio_tuberia.go` | Modified | Eliminar hardcode `1` y usar parámetro `tierras` existente |
| `internal/calculos/application/usecase/orquestador_memoria_calculo.go` | Modified | Calcular `numTierras` según regla (tubería) o usar 1 (charola) |
| `internal/calculos/application/dto/tuberia_input.go` | No change | No exponer campo; el cálculo es interno |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Compatibilidad con endpoint independiente de tubería | Low | El endpoint usa el mismo use case; si requiere el campo, se adicionará después |
| Cambios en la regla de negocio | Low | La regla está documentada; cualquier cambio requerirá actualización del código |

## Rollback Plan

1. Revertir cambios en `orquestador_memoria_calculo.go` — restaurar hardcode a `1`
2. Revertir cambios en `calcular_tamanio_tuberia.go` — restaurar parámetro a `1`
3. No se requieren cambios en DTOs

## Dependencies

- No hay dependencias externas
- El dominio ya soporta el parámetro `tierras` en `CalcularTamanioTuberiaWithMultiplePipes`

## Success Criteria

- [ ] El cálculo de tuberías con ≤2 tubos usa 1 hilo de tierra
- [ ] El cálculo de tuberías con >2 tubos usa 2 hilos de tierra
- [ ] El cálculo de charolas siempre usa 1 hilo de tierra
- [ ] El usuario NO puede especificar el número de hilos de tierra (campo no expuesto)
- [ ] La memoria de cálculo completa funciona correctamente
