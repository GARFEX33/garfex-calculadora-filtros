# Proposal: Agregar unidades de voltaje (V/kV)

## Intent

Actualmente el sistema solo acepta valores de voltaje como enteros (127, 220, 240, 277, 440, 480, 600 V) sin soportar kilovoltios (kV). La potencia ya tiene unidades (W, KW, KVA, KVAR) y el usuario necesita poder enviar voltajes en kV (ej: 0.48 kV = 480 V) para consistencia con el estándar NOM y mejor usabilidad.

## Scope

### In Scope
- Agregar campo `tension_unidad` al input de la API (V o kV)
- Modificar el value object `Tension` para aceptar unidades V y kV
- Normalizar internamente a volts (como se hace con potencia → watts)
- Actualizar DTOs y handlers HTTP para aceptar el nuevo campo
- Mantener compatibilidad hacia atrás con inputs existentes (default: V)

### Out of Scope
- Agregar soporte para otros voltajes no NOM (validación sigue siendo estricta)
- Modificar la lógica de cálculo de caída de tensión (ya soporta fase-neutro/fase-fase)
- Agregar UI o documentación de API

## Approach

Seguir el patrón existente de `Potencia`:
1. Agregar tipo `UnidadTension` con constantes `V` y `kV`
2. Modificar `NewTension` para aceptar `(valor float64, unidad string)`
3. Normalizar a volts internamente: `kV = V * 1000`
4. Agregar `Unidad()` getter al value object
5. Actualizar `EquipoInput` DTO y handlers para recibir `tension_unidad`

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/shared/kernel/valueobject/tension.go` | Modified | Agregar soporte para unidades V/kV |
| `internal/calculos/application/dto/equipo_input.go` | Modified | Agregar campo tension_unidad |
| `internal/calculos/infrastructure/adapter/driver/http/` | Modified | Actualizar handlers que reciben tensión |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Breaking change en API | Low | Default "V" si no se especifica unidad |
| Validación NOM rota | Low | Mantener los mismos valores válidos (127, 220, etc.) |

## Rollback Plan

1. Revertir cambios en `tension.go` al constructor original `NewTension(valor int)`
2. Eliminar campo `tension_unidad` del DTO
3. Revertir handlers HTTP
4. Los tests existentes deben seguir pasando

## Dependencies

- Ninguna dependencia externa

## Success Criteria

- [ ] API acepta `tension: 480, tension_unidad: "V"` (comportamiento actual)
- [ ] API acepta `tension: 0.48, tension_unidad: "kV"` y lo convierte correctamente
- [ ] Value object Tension retorna la unidad correcta
- [ ] Tests existentes pasan sin modificaciones
- [ ] Documentación de API actualizada (comments)
