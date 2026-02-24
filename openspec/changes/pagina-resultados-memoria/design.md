# Design: Página de Resultados de Memoria de Cálculo

## Technical Approach

El cambio implementa una nueva página de resultados dedicada que muestra la memoria de cálculo eléctrico en formato técnico-profesional. Se modifica el flujo existente para redirigir a la nueva página en lugar de mostrar resultados inline.

La implementación sigue las convenciones existentes del proyecto:
- Svelte 5 con runes (`$state`, `$derived`, `$props`)
- Tailwind CSS v4 con design tokens
- TypeScript estricto
- SvelteKit para routing

## Architecture Decisions

### Decision: Pass Data via URL Query Params

**Choice**: Serializar los datos del cálculo como query param codificado en base64
**Alternatives considered**: 
- SvelteKit page data (requiere reload)
- Svelte store (no persiste en navegación)
- LocalStorage (más complejo)

**Rationale**: Mantiene la URL compartible y permite recarga de página. Los datos de MemoriaOutput no son excesivamente grandes (~2KB JSON).

### Decision: Nuevo Componente Separado (no reutilizar ResultadosMemoria)

**Choice**: Crear `MemoriaTecnica.svelte` nuevo en lugar de modificar `ResultadosMemoria.svelte`
**Alternatives considered**: 
- Reutilizar el componente existente con prop "modo"
- Mostrar ambos formatos con condicional

**Rationale**: El formato actual (tarjetas UI) es diferente al técnico (fórmulas paso a paso). Mejor separación de responsabilidades y facilita mantenimiento.

### Decision: Ruta en /calculos/resultado

**Choice**: Crear `frontend/web/src/routes/calculos/resultado/+page.svelte`
**Alternatives considered**:
- `/resultados` (raíz)
- `/memoria/resultado`
- Sub-ruta de la página actual

**Rationale**: Coherente con estructura existente (`/calculos/`), sigue convenciones del proyecto.

## Data Flow

```
+page.svelte (formulario)
       │
       ▼ submit
API: POST /api/v1/calculos/memoria
       │
       ▼ response (MemoriaOutput)
handleSubmit()
       │
       ▼ encodeBase64(JSON.stringify(data))
goto('/calculos/resultado?data=...')
       │
       ▼
+page.svelte (resultado)
       │
       ▼ parse query param
MemoriaTecnica.svelte
       │
       ▼ render
HTML con fórmulas y pasos técnicos
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/web/src/routes/+page.svelte` | Modify | Cambiar `handleSubmit` para redirigir con `goto()` |
| `frontend/web/src/routes/calculos/resultado/+page.svelte` | Create | Nueva página de resultados |
| `frontend/web/src/routes/calculos/resultado/+page.ts` | Create | Load function para parsear query param |
| `frontend/web/src/lib/components/calculos/MemoriaTecnica.svelte` | Create | Componente de memoria técnica con formato profesional |
| `frontend/web/src/lib/components/calculos/SeccionCorriente.svelte` | Create | Sub-componente: cálculo de corriente |
| `frontend/web/src/lib/components/calculos/SeccionConductor.svelte` | Create | Sub-componente: dimensionamiento conductor |
| `frontend/web/src/lib/components/calculos/SeccionCaidaTension.svelte` | Create | Sub-componente: caída de tensión |
| `frontend/web/src/lib/components/calculos/SeccionCanalizacion.svelte` | Create | Sub-componente: canalización |
| `frontend/web/src/lib/components/calculos/SeccionConclusion.svelte` | Create | Sub-componente: conclusión técnica |

## Interfaces / Contracts

### Nueva Ruta: /calculos/resultado

**Query Parameter**: `?data=<base64-encoded-MemoriaOutput>`

```typescript
// frontend/web/src/routes/calculos/resultado/+page.ts
import type { PageLoad } from './$types';
import type { MemoriaOutput } from '$lib/types/calculos.types';

export const load: PageLoad = ({ url }) => {
  const data = url.searchParams.get('data');
  if (!data) {
    throw error(400, 'Faltan datos del cálculo');
  }
  
  try {
    const decoded = atob(data);
    const memoria: MemoriaOutput = JSON.parse(decoded);
    return { memoria };
  } catch {
    throw error(400, 'Datos inválidos');
  }
};
```

### Componentes de Sección

Cada sección será un sub-componente independiente para mejor organización:

```typescript
// Example: Sección de Corriente
interface Props {
  corrienteNominal: number;
  potencia?: number;
  voltaje: number;
  sistemaElectrico: SistemaElectrico;
  factorPotencia?: number;
}
```

### Tipos Existentes a Utilizar

- `MemoriaOutput` - datos completos del cálculo
- `ResultadoConductor` - datos del conductor
- `ResultadoCaidaTension` - datos de caída
- `ResultadoCanalizacion` - datos de tubería/charola

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Componentes individuales | Verificar renderizado con datos mock |
| Integration | Ruta completa | Test de navegación con datos |
| E2E | Flujo usuario | Playwright/Cypress (si existe) |

### Tests Prioritarios

1. `SeccionCorriente.svelte` - renderiza fórmula correcta según sistema
2. `SeccionConductor.svelte` - muestra todos los campos del conductor
3. `SeccionCaidaTension.svelte` - calcula porcentaje correctamente
4. Flujo: formulario → cálculo → resultados → regresar

## Migration / Rollback

**No se requiere migración de datos.** 

El cambio es completamente向后 compatible:
- El endpoint del backend no cambia
- Los tipos de datos no cambian
- Solo modifica la capa de presentación

**Rollback**:
1. Revertir cambios en `+page.svelte` (quitar `goto()`)
2. Restaurar ` ResultadosMemoria.svelte` como estaba
3. Eliminar carpeta `routes/calculos/resultado/`
4. Eliminar componentes creados

## Open Questions

- [ ] ¿El encode en base64 puede tener problemas con caracteres especiales? (Considerar URL-safe base64)
- [ ] ¿Cuántos datos es razonable pasar por URL antes de considerar alternativas?
- [ ] ¿Debe guardarse la memoria en BD para referencia futura? (fuera de scope)

---

**Decisiones técnicas confirmadas:**
- ✅ Redirección via query param con base64
- ✅ Nueva ruta `/calculos/resultado`
- ✅ Componentes de sección separados
- ✅ Usar tipos existentes de MemoriaOutput
- ✅ Sin cambio en backend
