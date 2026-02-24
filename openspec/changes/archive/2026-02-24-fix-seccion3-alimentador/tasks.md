# Tasks: fix-seccion3-alimentador

## Task Breakdown

### Phase 1: Implementación

#### Task 1.1: Eliminar sección "Factor Total de Ajuste" independencia
- **File**: `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`
- **Lines**: 138-145
- **Action**: Eliminar el bloque div completo de "Factor Total de Ajuste"
- **Verification**: Confirmar que la sección desaparece del render

#### Task 1.2: Corregir derivación de numHilosAlimentacion
- **File**: `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`
- **Line**: 33
- **Action**: Cambiar la derivación para manejar el caso de 1 hilo correctamente
- **Before**:
  ```typescript
  let numHilosAlimentacion = $derived(
    memoria.conductor_alimentacion.NumHilos ?? memoria.hilos_por_fase
  );
  ```
- **After**:
  ```typescript
  let numHilosAlimentacion = $derived(
    (memoria.conductor_alimentacion.NumHilos ?? memoria.hilos_por_fase) || 1
  );
  ```

#### Task 1.3: Renombrar label "Número de Hilos" a "Número de Hilos por Fase"
- **File**: `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`
- **Line**: 228
- **Action**: Actualizar el texto del label en la tabla
- **Before**: `Número de Hilos`
- **After**: `Número de Hilos por Fase`

---

### Phase 2: Verificación

#### Task 2.1: Verificar que la sección Factor Total fue eliminada
- **Check**: La sección独立性 "Factor Total de Ajuste" ya no aparece en el componente

#### Task 2.2: Verificar que "Desarrollo" sigue mostrando el cálculo
- **Check**: La fórmula con el factor total sigue visible en la sección "Desarrollo"

#### Task 2.3: Verificar que número de hilos muestra 1 correctamente
- **Check**: Con datos de prueba de 1 hilo, muestra "1" y no "0"

#### Task 2.4: Verificar label correcto
- **Check**: El label dice "Número de Hilos por Fase"

#### Task 2.5: Ejecutar QA checks
- **Command**: `cd frontend/web && npm run qa`
- **Expected**: Sin errores ni warnings

---

## Dependencies
- Ninguna — tareas independientes que pueden ejecutarse en cualquier orden dentro de cada fase

## Estimated Effort
- Implementación: ~10 minutos
- Verificación: ~10 minutos
- Total: ~20 minutos
