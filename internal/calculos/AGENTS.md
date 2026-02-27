# Feature: calculos

Memoria de cálculo eléctrico según normativa NOM (México).

Esta feature implementa el cálculo completo de una instalación eléctrica:
corriente nominal → ajuste por temperatura/agrupamiento → selección de conductor →
conductor de tierra → dimensionamiento de canalización → caída de tensión (NOM).

> Fórmula de caída de tensión (impedancia efectiva NOM/IEEE-141) → ver fuente de verdad en [domain/AGENTS.md](domain/AGENTS.md#convenciones-de-cálculo--caída-de-tensión)

## Endpoints

### Endpoint Memoria de Cálculo

El endpoint `POST /api/v1/calculos/memoria` orquesta todos los pasos secuencialmente:

1. Corriente nominal (desde potencia o amperaje)
2. Ajuste (factores de temperatura, agrupamiento, uso)
3. Selección de conductor de alimentación
4. Selección de conductor de tierra
5. Dimensionamiento de canalización (tubería o charola)
6. Cálculo de caída de tensión

**Unidades de potencia soportadas:**
- `W` — Watts
- `KW` — Kilowatts (default)
- `KVA` — Kilovolt-amperes
- `KVAR` — Kilovars reactivos

**Unidades de tensión soportadas:**
- `V` — Volts (default, compatibilidad hacia atrás)
- `kV` — Kilovolts (se normaliza internamente a V)

El valor ingresado se normaliza a volts antes de validar contra la lista NOM (127, 220, 240, 277, 440, 480, 600 V). Ejemplo: `0.48 kV` → `480 V`.

## Estructura

```
internal/calculos/
├── domain/          ← entidades y servicios de cálculo puro
│   ├── entity/      ← Proyecto, TipoCanalizacion, SistemaElectrico, etc.
│   └── service/     ← Servicios de cálculo NOM
├── application/     ← ports, use cases, DTOs
│   ├── port/        ← TablaNOMRepository, EquipoRepository
│   ├── usecase/     ← OrquestadorMemoriaCalculo y micro use cases
│   │   └── helpers/ ← Funciones auxiliares
│   └── dto/         ← EquipoInput, MemoriaOutput
└── infrastructure/  ← adapters HTTP (driver), CSV y PostgreSQL (driven)
    └── adapter/
        ├── driver/http/     ← CalculoHandler, formatters
        ├── driven/csv/      ← CSVTablaNOMRepository
        └── driven/postgres/ ← CalcEquipoFiltroRepository
```

### Integración con equipos (catálogo de filtros)

El modo `LISTADO` en `POST /api/v1/calculos/memoria` permite seleccionar un equipo del catálogo de equipos (`equipos_filtros`).

El adapter `CalcEquipoFiltroRepository` consulta la tabla y mapea cada tipo de filtro:

- `TipoFiltro = A` → `FiltroActivo` (corriente directa)
- `TipoFiltro = KVA` → `Transformador` (calcula I desde KVA)
- `TipoFiltro = KVAR` → `FiltroRechazo` (calcula I desde KVAR)

## Cómo modificar esta feature

Ver guías por capa:
- [domain/AGENTS.md](domain/AGENTS.md) — entidades y servicios de cálculo puro
- [application/AGENTS.md](application/AGENTS.md) — ports, use cases, DTOs
- [infrastructure/AGENTS.md](infrastructure/AGENTS.md) — adapters HTTP, CSV y PostgreSQL

> Ver estructura y reglas en [docs/reference/structure.md](../../../docs/reference/structure.md)

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
