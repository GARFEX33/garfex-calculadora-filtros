# Proposal: Equipo Selection Auto-Fill

## Intent

El problema es que cuando el usuario selecciona un equipo del catálogo (modo LISTADO), los campos de "Sistema Eléctrico", "Voltaje" y "Tipo de Voltaje" en la sección de instalación siguen siendo editables y vacíos.

La solución debe ser **100% desde el frontend** (sin modificar el backend):
- El backend espera `SistemaElectrico` y `TipoVoltaje` como campos separados en `EquipoInput`
- El catálogo de equipos tiene `conexion` (DELTA, ESTRELLA, etc.) y `tipo_voltaje` (FF, FN)
- El frontend debe mapear: `conexion` → `sistema_electrico`, `tipo_voltaje` → `tipo_voltaje`

Estos valores deben prellenarse automáticamente y los campos deben deshabilitarse.

## Scope

### In Scope
- Actualizar tipo `EquipoFiltro` en frontend para incluir `conexion` y `tipo_voltaje` (vienen del backend)
- Modificar `CamposInstalacion.svelte` para recibir prop de "solo lectura" (disabled)
- En `+page.svelte`:
  - Mapear `equipo.conexion` (DELTA/ESTRELLA/MONOFASICO/BIFASICO) → `sistema_electrico`
  - Mapear `equipo.tipo_voltaje` (FF/FN) → `tipo_voltaje` (FASE_FASE/FASE_NEUTRO)
  - Mapear `equipo.voltaje` → `tension`
  - Cuando hay equipo seleccionado, deshabilitar esos campos

### Out of Scope
- Modificar el backend
- Cambiar la estructura de la API

## Approach

1. **Actualizar tipos**: `equipos.types.ts` - agregar `conexion?: string` y `tipo_voltaje?: string`

2. **Modificar CamposInstalacion.svelte**:
   - Agregar prop `soloLectura: boolean`
   - Agregar `disabled={soloLectura}` a los campos: tensión, sistema eléctrico, tipo voltaje

3. **En +page.svelte**:
   - Agregar función `mapearConexionAProgramaElectrico(conexion)`:
     - DELTA → DELTA, ESTRELLA → ESTRELLA, MONOFASICO → MONOFASICO, BIFASICO → BIFASICO
   - Agregar función `mapearTipoVoltaje(tipoVoltaje)`:
     - FF → FASE_FASE, FN → FASE_NEUTRO
   - En `handleEquipoChange`:
     - Si equipo tiene `conexion`, actualizar `instalacion.sistema_electrico`
     - Si equipo tiene `tipo_voltaje`, actualizar `instalacion.tipo_voltaje`
     - Siempre actualizar `instalacion.tension` con `equipo.voltaje`
   - Pasar `soloLectura={!!equipoSeleccionado}` a `CamposInstalacion`

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/web/src/lib/types/equipos.types.ts` | Modified | Agregar campos `conexion` y `tipo_voltaje` |
| `frontend/web/src/lib/components/calculos/CamposInstalacion.svelte` | Modified | Agregar prop `soloLectura` y deshabilitar campos |
| `frontend/web/src/routes/+page.svelte` | Modified | Lógica de mapping y auto-fill |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Equipo sin conexion/tipo_voltaje (null) | Medium | Mantener campos editables si el valor es null |
| Mapeo de valores incorrecto | Low | FF→FASE_FASE, FN→FASE_NEUTRO son mapeos directos |

## Rollback Plan

1. Revertir cambios en los 3 archivos mencionados
2. El código volvería al estado actual con campos siempre editables

## Dependencies

- Ninguno - el backend ya acepta los campos correctamente, solo el frontend no los estaba enviando

## Success Criteria

- [ ] Al seleccionar un equipo, los campos de tensión, sistema eléctrico y tipo de voltaje se autocompletan con valores del equipo
- [ ] Los campos autocompletados están deshabilitados (no editables)
- [ ] Si se cambia a modo MANUAL, los campos снова son editables
- [ ] Si el equipo no tiene conexion/tipo_voltaje, los campos permanecen editables
- [ ] El mapping es correcto: conexion(DELTA/ESTRELLA/etc) → sistema_electrico, tipo_voltaje(FF/FN) → tipo_voltaje(FASE_FASE/FASE_NEUTRO)
- [ ] `npm run check` pasa sin errores de tipo
