# Design: Equipo Selection Auto-Fill

## Technical Approach

The change implements auto-fill and disable functionality for installation fields when an equipment is selected in LISTADO mode. The approach is purely frontend, using Svelte 5 runes for reactivity.

## Architecture Decisions

### Decision: Prop-based disable mechanism for CamposInstalacion

**Choice**: Add a `soloLectura` boolean prop to `CamposInstalacion.svelte` that controls the `disabled` attribute on specific fields.

**Alternatives considered**: 
- Use a separate "view mode" component
- Pass a signal/store with the equipment data

**Rationale**: This is the simplest approach that follows the existing pattern of passing data via props. It keeps the component composable and doesn't require creating a new component or store.

### Decision: Mapping functions for connection type conversion

**Choice**: Create pure mapping functions in `+page.svelte` to convert equipment connection values to the expected format.

**Alternatives considered**:
- Add mapping logic directly in the component
- Create a shared utility module

**Rationale**: Keeping mapping functions near where they're used (in the page) is simpler for this use case. If the mapping logic grows or is reused elsewhere, it can be extracted later.

## Data Flow

```
User selects equipment
        │
        ▼
handleEquipoChange(equipo)
        │
        ├─► equipo.conexion ──► mapearConexionASistemaElectrico() ──► instalacion.sistema_electrico
        │
        ├─► equipo.tipo_voltaje ──► mapearTipoVoltaje() ──► instalacion.tipo_voltaje
        │
        └─► equipo.voltaje ──► instalacion.tension
        │
        ▼
onDatosChange(instalacion) ──► CamposInstalacion.svelte
        │
        ▼
Render with disabled={soloLectura}
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/web/src/lib/types/equipos.types.ts` | Modify | Add `tipo_voltaje?: string \| null` to `EquipoFiltro` interface |
| `frontend/web/src/lib/components/calculos/CamposInstalacion.svelte` | Modify | Add `soloLectura?: boolean` prop and apply `disabled` attribute to relevant fields |
| `frontend/web/src/routes/+page.svelte` | Modify | Add mapping functions and call them in `handleEquipoChange`; pass `soloLectura` prop |

## Interfaces / Contracts

### Modified: EquipoFiltro interface

```typescript
export interface EquipoFiltro {
  id: string;
  clave: string;
  tipo: TipoFiltroEquipo;
  voltaje: number;
  amperaje: number;
  itm: number;
  bornes?: number | null;
  conexion?: string | null;        // Already exists
  tipo_voltaje?: string | null;    // NEW: FF or FN
  created_at: string;
}
```

### New: Mapping Functions (in +page.svelte)

```typescript
function mapearConexionASistemaElectrico(conexion: string | null): SistemaElectrico | '' {
  if (!conexion) return '';
  
  const mapa: Record<string, SistemaElectrico> = {
    'DELTA': 'DELTA',
    'ESTRELLA': 'ESTRELLA',
    'MONOFASICO': 'MONOFASICO',
    'BIFASICO': 'BIFASICO'
  };
  return mapa[conexion] || '';
}

function mapearTipoVoltaje(tipoVoltaje: string | null): TipoVoltaje | '' {
  if (!tipoVoltaje) return '';
  
  const mapa: Record<string, TipoVoltaje> = {
    'FF': 'FASE_FASE',
    'FN': 'FASE_NEUTRO'
  };
  return mapa[tipoVoltaje] || '';
}
```

### Modified: CamposInstalacion Props

```typescript
interface Props {
  datos: CamposInstalacionData;
  onDatosChange: (datos: CamposInstalacionData) => void;
  soloLectura?: boolean;  // NEW: defaults to false
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Mapping functions | Test conversion of all valid values and null handling |
| Component | CamposInstalacion soloLectura prop | Verify disabled attribute is applied correctly |
| Integration | Full flow from equipment selection to form submission | Manual testing with browser |

## Migration / Rollout

No migration required. This is a pure frontend enhancement with no database or API changes.

## Open Questions

- None. The implementation is straightforward and follows existing patterns.
