# Design: fix-seccion3-alimentador

## Technical Design

### Component Affected
- `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`

### Changes Required

#### CHANGE 1: Eliminar sección "Factor Total de Ajuste"

**Location**: Líneas 138-145 del componente

**Current Code**:
```svelte
<!-- Factor Total -->
<div class="mb-4 rounded border border-primary/20 bg-primary/5 p-4">
	<h3 class="mb-2 font-semibold text-foreground">Factor Total de Ajuste</h3>
	<p class="font-mono text-sm text-foreground">
		F<sub>total</sub> = {factorUso.toFixed(2)} × {memoria.factor_temperatura.toFixed(2)} ×
		{memoria.factor_agrupamiento.toFixed(2)} = {memoria.factor_total_ajuste.toFixed(3)}
	</p>
</div>
```

**Action**: Eliminar estas líneas completamente (138-145)

**Rationale**: Esta información ya está expresada en la sección "Desarrollo" (líneas 157-183) donde se muestra la fórmula completa:
```
Iajustada = Inominal × Fuso / (Ftemp × Fagr)
```

---

#### CHANGE 2: Corregir derivación de numHilosAlimentacion

**Location**: Línea 33 del componente

**Current Code**:
```typescript
let numHilosAlimentacion = $derived(
  memoria.conductor_alimentacion.NumHilos ?? memoria.hilos_por_fase
);
```

**Issue**: El operador `??` solo maneja `null` y `undefined`, pero no maneja el caso donde el valor es `0`.

**Proposed Fix**:
```typescript
let numHilosAlimentacion = $derived(
  (memoria.conductor_alimentacion.NumHilos ?? memoria.hilos_por_fase) || 1
);
```

**Rationale**: Si ambos valores son 0, null, o undefined, el valor predeterminado debe ser 1 (un hilo por fase es el caso más común).

---

#### CHANGE 3: Renombrar label de Número de Hilos

**Location**: Línea 228 del componente

**Current Code**:
```svelte
<tr>
	<td class="px-4 py-2 text-muted-foreground">Número de Hilos</td>
	<td class="px-4 py-2 text-foreground">
		{numHilosAlimentacion}
	</td>
</tr>
```

**Proposed Change**:
```svelte
<tr>
	<td class="px-4 py-2 text-muted-foreground">Número de Hilos por Fase</td>
	<td class="px-4 py-2 text-foreground">
		{numHilosAlimentacion}
	</td>
</tr>
```

**Rationale**: Técnicamente es más preciso ya que cada fase puede tener múltiples conductores en paralelo.

---

## Implementation Order

1. Apply CHANGE 1 (eliminar sección)
2. Apply CHANGE 2 (corregir derivación)
3. Apply CHANGE 3 (renombrar label)

## Testing Approach

### Manual Testing
1. Abrir la calculadora con datos de prueba
2. Verificar que la sección "Factor Total de Ajuste"独立性 no aparezca
3. Verificar que la sección "Desarrollo" muestre correctamente el cálculo
4. Verificar que cuando hay 1 hilo, muestre "1" y no "0"
5. Verificar el label dice "Número de Hilos por Fase"

### Code Review
- Verificar que los cambios no rompan el layout
- Verificar que el diseño responsive se mantenga

## Dependencies
- No hay dependencias externas
- Solo cambios en el componente Svelte afectado

## No Backend Changes Required
Este cambio es puramente de frontend y no requiere modificaciones en el backend de Go.
