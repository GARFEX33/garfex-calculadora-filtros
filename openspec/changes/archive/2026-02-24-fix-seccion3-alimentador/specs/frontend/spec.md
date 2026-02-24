# Specs: fix-seccion3-alimentador (Frontend)

## Overview
Especificaciones delta para las correcciones en la Sección 3 "Dimensionamiento del Alimentador".

## References
- Componente: `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`

---

## SPEC-001: Eliminar Factor Total de Ajuste redundante

### Description
La sección "Factor Total de Ajuste"独立性 aparece antes de la sección "Desarrollo". Ya existe una expresión equivalente del factor total en la sección "Desarrollo" que es más completa. Esta sección separada es redundante y debe eliminarse.

### Requirement
- La sección "Factor Total de Ajuste" (actualmente líneas 138-145) debe ser eliminada del componente.

### Validation Criteria
- [ ] La sección "Factor Total de Ajuste" ya no aparece en el componente
- [ ] La sección "Desarrollo" sigue mostrando correctamente la fórmula con el factor total

---

## SPEC-002: Corregir número de hilos cuando es 1

### Description
Actualmente, cuando hay 1 hilo por fase, el componente muestra "Número de Hilos = 0" en lugar de "1".

### Requirement
- El valor de "Número de Hilos" debe ser 1 cuando hay un solo hilo por fase.
- No debe mostrar 0 en ningún caso válido.

### Validation Criteria
- [ ] Cuando `numHilosAlimentacion` es 1, el componente muestra "1"
- [ ] No se muestra "0" para ningún caso válido
- [ ] El cálculo de `$derived` de `numHilosAlimentacion` maneja correctamente el caso de 1 hilo

### Technical Detail
El código actual es:
```typescript
let numHilosAlimentacion = $derived(
  memoria.conductor_alimentacion.NumHilos ?? memoria.hilos_por_fase
);
```
Debe asegurar que cuando ambos valores son 0 o undefined, se use 1 como默认值.

---

## SPEC-003: Renombrar "Número de Hilos" a "Número de Hilos por Fase"

### Description
El label "Número de Hilos" no es técnicamente preciso. El término correcto es "hilos por fase" ya que cada fase puede tener múltiples conductores en paralelo.

### Requirement
- El label de la fila en la tabla "Conductor Seleccionado" debe cambiar de "Número de Hilos" a "Número de Hilos por Fase".

### Validation Criteria
- [ ] El label "Número de Hilos por Fase" aparece en la tabla de Conductor Seleccionado
- [ ] El valor correcto (1 o más) se muestra junto al label

---

## Scenarios de Prueba

### Scenario 1: Conductor con 1 hilo por fase
**Given**: Un cálculo con `NumHilos = 1` o `hilos_por_fase = 1`
**Then**:
- "Número de Hilos por Fase" muestra "1"
- No hay mensaje de "0 hilos"

### Scenario 2: Conductor con múltiples hilos en paralelo
**Given**: Un cálculo con `NumHilos = 3` o `hilos_por_fase = 3`
**Then**:
- "Número de Hilos por Fase" muestra "3"
- La capacidad total se calcula correctamente (Capacidad × 3)

### Scenario 3: Verificar ausencia de Factor Total redundante
**Given**: La sección 3 "Dimensionamiento del Alimentador"
**Then**:
- La sección独立性 "Factor Total de Ajuste" no aparece
- La fórmula en "Desarrollo" sigue mostrando el cálculo del factor total
