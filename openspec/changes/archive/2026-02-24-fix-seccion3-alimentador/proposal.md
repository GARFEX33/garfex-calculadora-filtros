# Proposal: fix-seccion3-alimentador

## Change Name
`fix-seccion3-alimentador` — Correcciones en Sección 3 "Dimensionamiento del Alimentador"

## Context
El proyecto tiene una calculadora de instalaciones eléctricas. La Sección 3 ("Dimensionamiento del Alimentador") ya fue mejorada recientemente para mostrar factores de corrección. Sin embargo, se identificaron dos problemas y una limpieza necesaria.

## Scope

### In Scope
- **Frontend only**: `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`
- No se requieren cambios en el backend

### Out of Scope
- Cambios en la lógica de negocio del backend
- Cambios en otros componentes de la memoria de cálculo
- Nuevas funcionalidades

## Problems to Solve

### Problem 1: Eliminar "Factor Total" redundante
**Descripción**: El "Factor Total de Ajuste" aparece tanto en una sección separada (líneas 138-145) como expresado en la sección de "Desarrollo" (líneas 157-183). La expresión en Desarrollo es más completa y clara.

**Solución**: Eliminar la sección "Factor Total de Ajuste"独立性 (líneas 138-145) ya que es redundante.

### Problem 2: Bug - Número de Hilos muestra 0 cuando es 1
**Descripción**: Cuando hay 1 hilo por fase, actualmente muestra "Número de Hilos = 0" en lugar de "1".

**Causa probable**: El valor de `numHilosAlimentacion` podría ser 0 o undefined cuando `memoria.conductor_alimentacion.NumHilos` es null/undefined y `memoria.hilos_por_fase` también es 0 o undefined.

**Solución**: Asegurar que cuando hay 1 hilo, el valor mostrado sea 1.

### Problem 3: Mejora terminológica
**Descripción**: El label "Número de Hilos" no es técnicamente correcto.

**Solución**: Cambiar a "Número de Hilos por Fase" para aclarar que cada fase tiene ese número de conductores.

## Approach

1. **Análisis**: Revisar el código actual y entender el flujo de datos
2. **Correcciones**:
   - Eliminar la sección "Factor Total de Ajuste"独立性
   - Corregir el cálculo/derivación de `numHilosAlimentacion` para manejar el caso de 1 hilo
   - Actualizar el label de "Número de Hilos" a "Número de Hilos por Fase"
3. **Verificación**: Verificar visualmente que los cambios son correctos

## Rollback Plan
Si los cambios causan problemas visuales o regressions:
- Revertir los cambios en `SeccionAlimentador.svelte`
- No hay cambios en backend que requieran rollback

## Affected Modules
- `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`

## Risk Level
**Bajo** — Solo cambios visuales en un componente de presentación
