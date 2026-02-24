# Tasks: Equipo Selection Auto-Fill

## Phase 1: Type Definitions (Foundation)

- [ ] 1.1 Add `tipo_voltaje?: string | null` field to `EquipoFiltro` interface in `frontend/web/src/lib/types/equipos.types.ts`
- [ ] 1.2 Run `npm run check` to verify type definitions compile correctly

## Phase 2: CamposInstalacion Component (Core)

- [ ] 2.1 Add `soloLectura?: boolean` prop to `CamposInstalacion.svelte` interface
- [ ] 2.2 Add `disabled` attribute to "Tensión" input field using conditional: `disabled={soloLectura}`
- [ ] 2.3 Add `disabled` attribute to "Sistema Eléctrico" select using conditional: `disabled={soloLectura}`
- [ ] 2.4 Add `disabled` attribute to "Tipo de Voltaje" radio buttons using conditional: `disabled={soloLectura}`
- [ ] 2.5 Test that CamposInstalacion renders correctly with soloLectura=true

## Phase 3: Main Page Integration (Wiring)

- [ ] 3.1 Add `mapearConexionASistemaElectrico(conexion: string | null)` function in `+page.svelte`
  - Maps: DELTA→DELTA, ESTRELLA→ESTRELLA, MONOFASICO→MONOFASICO, BIFASICO→BIFASICO
  - Returns empty string for null/unknown values
- [ ] 3.2 Add `mapearTipoVoltaje(tipoVoltaje: string | null)` function in `+page.svelte`
  - Maps: FF→FASE_FASE, FN→FASE_NEUTRO
  - Returns empty string for null/unknown values
- [ ] 3.3 Modify `handleEquipoChange` function to:
  - Update `instalacion.tension` with `equipo.voltaje`
  - Call mapping functions and update `instalacion.sistema_electrico` and `instalacion.tipo_voltaje`
- [ ] 3.4 Pass `soloLectura={!!equipoSeleccionado}` prop to `CamposInstalacion` component
- [ ] 3.5 Handle mode switch: when switching from LISTADO to MANUAL, ensure fields become editable (already handled by the prop logic)
- [ ] 3.6 Run `npm run check` to verify no TypeScript errors

## Phase 4: Testing & Verification

- [ ] 4.1 Manual test: Select equipment with all fields populated, verify auto-fill works
- [ ] 4.2 Manual test: Select equipment with null conexion/tipo_voltaje, verify fields remain editable
- [ ] 4.3 Manual test: Switch from LISTADO to MANUAL mode, verify fields become editable
- [ ] 4.4 Manual test: Submit calculation with selected equipment, verify correct data sent to backend
- [ ] 4.5 Run `npm run qa` to ensure lint, format, and type checks pass

## Phase 5: Optional Enhancement (Display Details)

- [ ] 5.1 Display `conexion` and `tipo_voltaje` in the equipment selection card (in `FormularioListado.svelte`)
- [ ] 5.2 Run `npm run check` after changes
