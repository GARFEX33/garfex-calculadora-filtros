# Design: desarrollo-formula-corriente-pdf

## Technical Approach

This change implements the step-by-step visualization of the nominal current formula in the PDF calculation memory. The strategy consists of:

1. **Extend DTO** `MemoriaOutput` with new structs for the development
2. **Create helper function** in the orchestrator that generates development data replicating frontend logic
3. **Add template functions** to iterate over maps and format steps
4. **Modify template** `seccion_corriente.html` to use the new structured data

## Architecture Decisions

### Decision: Development logic location

**Choice**: Helper function within the orchestrator's `usecase` package  
**Alternatives considered**: 
- Create new use case `GenerarDesarrolloCorrienteUseCase`
- Move logic to domain service

**Rationale**: The current development is a derived calculation that doesn't make sense as an independent use case. The orchestrator already has access to all necessary data (`TipoEquipo`, `CorrienteNominal`, `Tension`, `FactorPotencia`, `SistemaElectrico`). Additionally, keeping it in the orchestrator avoids creating a new dependency in the use case graph.

### Decision: Data structure for ValoresReferencia

**Choice**: `map[string]string` for reference values  
**Alternatives considered**:
- Fixed struct with optional fields
- Array of key-value pairs

**Rationale**: 
- Each equipment type has different values (KVA vs Potencia vs Amperaje)
- `map[string]string` provides flexibility and the template can iterate over existing keys
- It's the same pattern used by the frontend in `getInfoCalculo()`

### Decision: Field name in MemoriaOutput

**Choice**: `DesarrolloCorriente`  
**Alternatives considered**: 
- `DesarrolloFormula`
- `DetalleCorriente`
- `PasosCorriente`

**Rationale**: `DesarrolloCorriente` is consistent with frontend nomenclature and the section name in the PDF ("Cálculo de Corriente Nominal").

### Decision: Replace vs. enhance current template

**Choice**: Replace existing conditional logic with structured data  
**Alternatives considered**: 
- Add new section alongside current logic
- Keep both and use feature flag

**Rationale**: 
- The current template uses nested `{{if}}` for each equipment type
- The new structured data makes that conditional logic unnecessary
- The result is cleaner and more maintainable

## Data Flow

```
OrquestadorMemoriaCalculoUseCase.Execute()
    │
    ├── Step 1: CalcularCorrienteUseCase → CorrienteNominal
    │
    ├── Step 2: AjustarCorrienteUseCase → CorrienteAjustada
    │
    ├── ... (other steps)
    │
    └── AT THE END (after Step 1):
        invoke generarDesarrolloCorriente(memoria) → DatosDesarrolloCorriente
        → assign to output.DesarrolloCorriente
            │
            ▼
    PDF Template (seccion_corriente.html)
        │
        ├── {{.Memoria.DesarrolloCorriente.TipoCalculo}}
        ├── {{.Memoria.DesarrolloCorriente.FormulaUsada}}
        ├── range .Memoria.DesarrolloCorriente.PasosDesarrollo
        └── range .Memoria.DesarrolloCorriente.ValoresReferencia
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/calculos/application/dto/memoria_output.go` | Modify | Add structs `DatosDesarrolloCorriente`, `PasoDesarrollo` and field `DesarrolloCorriente` to `MemoriaOutput` |
| `internal/calculos/application/usecase/orquestador_memoria_calculo.go` | Modify | Add function `generarDesarrolloCorriente()` and invoke it after Step 1 |
| `internal/pdf/infrastructure/adapter/driven/template/html_renderer.go` | Modify | Add template functions `rangeMap`, `index`, `add` to iterate over maps |
| `internal/pdf/templates/partials/seccion_corriente.html` | Modify | Replace conditional logic with structured data from `DesarrolloCorriente` |

## Interfaces / Contracts

### New DTO Structs

```go
// DatosDesarrolloCorriente contains the step-by-step development of the nominal current calculation.
type DatosDesarrolloCorriente struct {
    TipoCalculo       string            `json:"tipo_calculo"`
    FormulaUsada      string            `json:"formula_usada"`
    PasosDesarrollo   []PasoDesarrollo  `json:"pasos_desarrollo"`
    ValoresReferencia map[string]string `json:"valores_referencia"`
}

// PasoDesarrollo represents an individual step in the calculation development.
type PasoDesarrollo struct {
    Numero      int    `json:"numero"`
    Descripcion string `json:"descripcion"`
    Resultado   string `json:"resultado"`
}
```

### Field added to MemoriaOutput

```go
type MemoriaOutput struct {
    // ... existing fields ...
    
    // DesarrolloCorriente contains the step-by-step development of the current calculation.
    // Calculated in the orchestrator after determining the nominal current.
    DesarrolloCorriente *DatosDesarrolloCorriente `json:"desarrollo_corriente,omitempty"`
}
```

### New Template Functions

```go
// rangeMap iterates over a map[string]string, returning slice of {key, value}
"rangeMap": func(m map[string]string) []struct{Key, Value string} {
    result := make([]struct{Key, Value string}, 0, len(m))
    for k, v := range m {
        result = append(result, struct{Key, Value string}{k, v})
    }
    return result
},
// add sums two integers (for indices)
"add": func(a, b int) int {
    return a + b
},
```

## Testing Strategy

| Layer | What to test | Approach |
|-------|-------------|----------|
| Unit (DTO) | JSON serialization of new structs | Unmarshaling test with known values |
| Unit (orchestrator) | Generated development data for each equipment type | Compare against expected values (KVA, Potencia, etc.) |
| Unit (template) | Template rendering with complete data | Verify generated HTML contains expected elements |
| Integration | Endpoint `/api/v1/calculos/memoria` returns new field | Integration test with HTTP client |

### Test cases for generarDesarrolloCorriente

| Equipment Type | System | Verify |
|----------------|--------|--------|
| FILTRO_ACTIVO | any | type="Amperaje directo", valores["Amperaje"] |
| TRANSFORMADOR | any | type="Desde KVA (Transformador)", valores["KVA"], valores["Voltaje"] |
| FILTRO_RECHAZO | any | type="Desde KVAR (Filtro de Rechazo)", valores["KVAR"], valores["Voltaje"] |
| CARGA | ESTRELLA/DELTA | type="Desde Potencia (Sistema Trifásico)", valores["Potencia"], valores["Factor de Potencia"], valores["Sistema"] |
| CARGA | MONOFASICO/BIFASICO | type="Desde Potencia (Sistema Monofásico)", valores["Potencia"], valores["Factor de Potencia"], valores["Sistema"] |

## Migration / Rollout

No migration required. The change is backward-compatible:
- Field `DesarrolloCorriente` is `omitempty`
- If not populated, template can fallback to previous behavior
- Recommendation: implement new template and remove old logic in the same PR

## Open Questions

- [ ] Should `DesarrolloCorriente` also be calculated for the case where only nominal current is requested (endpoint `/api/v1/calculos/corriente`)? For now: NO, only for complete memory.
- [ ] Should consistency validation with frontend be added? The specification indicates values must match (max 0.01 difference). For now, replicate exact frontend formula.
