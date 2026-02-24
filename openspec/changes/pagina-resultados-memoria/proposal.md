# Proposal: Página de Resultados de Memoria de Cálculo

## Intent

El usuario necesita una página dedicada de resultados que muestre la memoria de cálculo eléctrico en formato técnico-profesional, siguiendo la estructura de la NOM-001-SEDE-2012. Actualmente los resultados se muestran inline en la página del formulario; se necesita separar el flujo y mostrar los cálculos paso a paso con fórmulas, sustituciones y justificaciones técnicas.

## Scope

### In Scope
- Crear nueva ruta `/calculos/resultado` para mostrar memoria de cálculo
- Modificar flujo del botón "Calcular Memoria" para redirigir en lugar de mostrar inline
- Crear componente `MemoriaTecnica.svelte` con formato profesional:
  - Encabezado (empresa, proyecto, equipo, capacidad, voltaje, longitud)
  - Cálculo de corriente nominal (fórmula + sustitución numérica)
  - Dimensionamiento del alimentador (Artículo 310-15, factor 125%/135%)
  - Conductor de puesta a tierra (Tabla 250-122)
  - Cálculo de caída de tensión (fórmula, DMG, RMG, impedancia)
  - Cálculo de charola/tubería
  - Conclusión técnica
- Pasar datos del cálculo via URL (query params) o state
- Mantener compatibilidad con datos existentes del endpoint

### Out of Scope
- Modificación del backend (el endpoint ya existe y funciona)
- Guardado permanente de memorias (solo visualización)
- Exportación a PDF (futuro)
- Autenticación/autorización

## Approach

1. **Nueva ruta**: Crear `frontend/web/src/routes/calculos/resultado/+page.svelte`
2. **Modificar flujo**:
   - En `+page.svelte` actual: cambiar `handleSubmit` para usar `goto('/calculos/resultado?data=...')` 
   - O usar navigation state de SvelteKit
3. **Componente de visualización**:
   - Usar datos existentes de `MemoriaOutput`
   - Renderizar cada paso con estructura técnica: fórmula → desarrollo → resultado
   - Incluir referencias a normas NOM (Art. 310-15, Tabla 250-122, etc.)
4. **Estilos**: Usar design tokens existentes de Tailwind

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/web/src/routes/+page.svelte` | Modified | Cambiar handleSubmit para redirigir |
| `frontend/web/src/routes/calculos/resultado/+page.svelte` | New | Nueva página de resultados |
| `frontend/web/src/lib/components/calculos/MemoriaTecnica.svelte` | New | Componente de memoria técnica |
| `frontend/web/src/lib/types/calculos.types.ts` | Modified | Añadir tipos si es necesario |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Datos muy grandes en URL | Medium | Usar POST con form o state de SvelteKit en lugar de query params |
| Migrar todos los campos del DTO | Low | El DTO ya tiene la info necesaria |
| Estilos不一致 | Low | Usar tokens existentes de Tailwind |

## Rollback Plan

1. Revertir cambios en `+page.svelte` para mostrar resultados inline como antes
2. Eliminar nueva ruta y componente creado
3. No se requiere cambio en backend

## Success Criteria

- [ ] Al hacer click en "Calcular Memoria" redirige a `/calculos/resultado`
- [ ] La página de resultados muestra todos los pasos del cálculo
- [ ] Cada paso incluye fórmula, desarrollo y resultado con unidades
- [ ] Incluye referencias a normas NOM
- [ ] Diseño responsive y usa design tokens existentes
- [ ] Tests pasan (QA checklist)
