# Tasks: Página de Resultados de Memoria de Cálculo

## Phase 1: Foundation (Estructura base y configuración)

- [ ] 1.1 Crear directorio `frontend/web/src/routes/calculos/resultado/`
- [ ] 1.2 Crear directorio `frontend/web/src/lib/components/calculos/secciones/` para sub-componentes
- [ ] 1.3 Verificar que existen los tipos `MemoriaOutput` y relacionados en `$lib/types/calculos.types.ts`

## Phase 2: Core Implementation - Página de Resultados

- [ ] 2.1 Crear `frontend/web/src/routes/calculos/resultado/+page.ts` con load function
  - Parsear query param `data` (base64 → JSON)
  - Validar que existen datos
  - Retornar `{ memoria }` al page component
- [ ] 2.2 Crear `frontend/web/src/routes/calculos/resultado/+page.svelte`
  - Usar datos del load
  - Incluir botón de "Nuevo Cálculo" para regresar
  - Importar y renderizar `MemoriaTecnica`

## Phase 3: Core Implementation - Componentes de Sección

- [ ] 3.1 Crear `frontend/web/src/lib/components/calculos/MemoriaTecnica.svelte` (componente principal)
  - Renderizar encabezado
  - Renderizar secciones en orden: Corriente → Conductor → Tierra → Caída → Canalización → Conclusión
- [ ] 3.2 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionEncabezado.svelte`
  - Mostrar: título, equipo, capacidad (kVA/A), voltaje, longitud
  - Usar datos de `memoria.equipo`, `memoria.tension`, `memoria.longitud_circuito`
- [ ] 3.3 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionCorriente.svelte`
  - Fórmula trifásica: `I = S / (√3 × V)`
  - Fórmula monofásica: `I = P / (V × cosθ)`
  - Mostrar desarrollo con valores numéricos
  - Referenciar sistema eléctrico: `memoria.sistema_electrico`
- [ ] 3.4 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionAlimentador.svelte`
  - Referenciar Artículo 310-15(b)(17)
  - Mostrar factor de diseño (125% o 135%)
  - Fórmula: `I_diseño = Factor × I_nominal`
  - Mostrar conductor seleccionado: calibre, material, aislamiento, ampacidad
- [ ] 3.5 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionTierra.svelte`
  - Referenciar Tabla 250-122
  - Mostrar ITM y calibre seleccionado
- [ ] 3.6 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionCaidaTension.svelte`
  - Fórmula trifásica: `%V = (√3 × I × Z × L / V) × 100`
  - Fórmula monofásica: `%V = (2 × I × Z × L / V) × 100`
  - Mostrar impedancia, desarrollo, resultado
  - Indicador de cumplimiento (verde/rojo)
- [ ] 3.7 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionCanalizacion.svelte`
  - Mostrar tipo (tubería/charola)
  - Mostrar tamaño comercial, fill factor
- [ ] 3.8 Crear `frontend/web/src/lib/components/calculos/secciones/SeccionConclusion.svelte`
  - Resumen de cumplimiento
  - Lista de observaciones (si existen)
  - Mensaje final según `memoria.cumple_normativa`

## Phase 4: Modificación del Flujo Existente

- [ ] 4.1 Modificar `frontend/web/src/routes/+page.svelte` - función `handleSubmit`
  - En caso de éxito: codificar `resultado` en base64
  - Llamar `goto('/calculos/resultado?data=' + encodedData)`
  - Eliminar código que muestra `ResultadosMemoria` inline
- [ ] 4.2 Opcional: mantener `ResultadosMemoria.svelte` para referencia o eliminar si no se usa

## Phase 5: Testing y Verificación

- [ ] 5.1 Verificar que la ruta `/calculos/resultado` carga correctamente
- [ ] 5.2 Probar flujo completo: formulario → calcular → resultados
- [ ] 5.3 Verificar que al hacer clic en "Nuevo Cálculo" regresa al formulario
- [ ] 5.4 Probar con datos de ejemplo (crear mock si es necesario)
- [ ] 5.5 Verificar diseño responsive (mobile, tablet, desktop)
- [ ] 5.6 Ejecutar `npm run qa` y verificar que pasa sin errores

## Phase 6: Cleanup

- [ ] 6.1 Eliminar comentarios de debug si existen
- [ ] 6.2 Verificar que no quedan console.log de desarrollo
- [ ] 6.3 Documentar cualquier decisión de diseño en comentarios (si es necesario)

---

## Order de Implementación Recomendado

1. **Primero**: Phase 1 (directorios)
2. **Segundo**: Phase 2 (página de resultados básica)
3. **Tercero**: Phase 3 (componentes de sección) - empezar por Encabezado y Conclusión
4. **Cuarto**: Phase 4 (modificar flujo actual)
5. **Quinto**: Phase 5 (testing)
6. **Sexto**: Phase 6 (cleanup)

**Justificación**: Los componentes de sección pueden probarse independientemente una vez que la página básica funcione.
