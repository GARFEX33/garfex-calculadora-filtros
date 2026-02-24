# Delta for Equipo Selection Auto-Fill

## ADDED Requirements

### Requirement: Auto-fill Installation Fields from Selected Equipment

When a user selects an equipment from the catalog in LISTADO mode, the system MUST automatically populate the installation fields with values from the selected equipment.

The following fields MUST be auto-filled:
- **Tensión**: MUST be set to the equipment's `voltaje` value
- **Sistema Eléctrico**: MUST be set by mapping equipment's `conexion` field:
  - DELTA → DELTA
  - ESTRELLA → ESTRELLA  
  - MONOFASICO → MONOFASICO
  - BIFASICO → BIFASICO
- **Tipo de Voltaje**: MUST be set by mapping equipment's `tipo_voltaje` field:
  - FF → FASE_FASE
  - FN → FASE_NEUTRO

#### Scenario: Select equipment with all connection fields populated

- GIVEN the user is in LISTADO mode
- AND the equipment has `conexion: "DELTA"` and `tipo_voltaje: "FF"` and `voltaje: 440`
- WHEN the user selects the equipment from the catalog
- THEN the "Tensión" field MUST be pre-filled with "440"
- AND the "Sistema Eléctrico" dropdown MUST show "Delta" as selected
- AND the "Tipo de Voltaje" radio MUST show "Fase-Fase" as selected

#### Scenario: Select equipment with null connection fields

- GIVEN the user is in LISTADO mode
- AND the equipment has `conexion: null` and `tipo_voltaje: null`
- WHEN the user selects the equipment from the catalog
- THEN the installation fields MUST remain empty and editable
- AND no auto-fill occurs

#### Scenario: Switch from LISTADO to MANUAL mode

- GIVEN the user has selected an equipment in LISTADO mode
- AND the installation fields are auto-filled and disabled
- WHEN the user switches to MANUAL_AMPERAJE or MANUAL_POTENCIA mode
- THEN all installation fields MUST become editable again
- AND the auto-filled values from the equipment MUST be cleared

### Requirement: Disable Installation Fields When Equipment Selected

When an equipment is selected in LISTADO mode, the system MUST disable the following installation fields to prevent manual editing:
- Tensión (input)
- Sistema Eléctrico (dropdown)
- Tipo de Voltaje (radio buttons)

The fields MUST remain disabled as long as an equipment is selected.

#### Scenario: Equipment selected - fields disabled

- GIVEN the user has selected an equipment in LISTADO mode
- WHEN the user attempts to change the "Sistema Eléctrico" dropdown
- THEN the dropdown MUST be disabled
- AND the user MUST NOT be able to select a different value

#### Scenario: Equipment deselected - fields enabled

- GIVEN the user had selected an equipment in LISTADO mode
- AND the installation fields are disabled
- WHEN the user clears the equipment selection (or no equipment is selected)
- THEN all installation fields MUST become editable again

### Requirement: Display Connection Details in Equipment Card

The equipment selection card MUST display the connection and voltage type information when available.

#### Scenario: Display full equipment details

- GIVEN the user has selected an equipment with `conexion: "ESTRELLA"` and `tipo_voltaje: "FN"`
- WHEN the equipment card is displayed
- THEN the card MUST show the connection type "ESTRELLA"
- AND MUST show the voltage type "FN" (or mapped "Fase-Neutro")

## MODIFIED Requirements

### Requirement: Equipment Type Definition

The `EquipoFiltro` type in the frontend MUST include the following fields that were previously missing:

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `conexion` | `string \| null` | Backend `conexion` | Connection type: DELTA, ESTRELLA, MONOFASICO, BIFASICO |
| `tipo_voltaje` | `string \| null` | Backend `tipo_voltaje` | Voltage reference: FF (Fase-Fase), FN (Fase-Neutro) |

(Previously: These fields were not present in the frontend type definition)

### Requirement: CalcularMemoriaRequest Includes Equipment Connection Data

The frontend MUST send the equipment's connection data to the backend when in LISTADO mode.

(Previously: Only basic equipment data was sent; connection fields were not mapped to installation fields)
