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

### Respuesta de Memoria de Cálculo

La respuesta de `POST /api/v1/calculos/memoria` está **agrupada por entidad** en lugar de campos planos secuenciales:

```json
{
  "equipo": { "clave": "...", "tipo": "...", "voltaje": 480, "amperaje": 100, "itm": 125 },
  "tipo_equipo": "CARGA",
  "factor_potencia": 0.85,
  "estado": "Ciudad de Mexico",
  "instalacion": {
    "tension": 480,
    "sistema_electrico": "ESTRELLA",
    "tipo_canalizacion": "TUBERIA_PVC",
    "material": "Cu",
    "longitud_circuito": 50,
    "hilos_por_fase": 1,
    "porcentaje_caida_maximo": 3.0
  },
  "corrientes": {
    "corriente_nominal": 45.08,
    "corriente_ajustada": 56.35,
    "corriente_por_hilo": 56.35,
    "factor_temperatura": 0.91,
    "factor_agrupamiento": 1.0,
    "factor_total_ajuste": 0.91,
    "temperatura_ambiente": 30,
    "temperatura_referencia": 75,
    "conductores_por_tubo": 3,
    "cantidad_conductores": 4,
    "tabla_ampacidad_usada": "Tabla 310-16"
  },
  "cable_fase": { "calibre": "8 AWG", "material": "Cu", "seccion_mm2": 8.37, ... },
  "cable_neutro": null,
  "cable_tierra": { "calibre": "12 AWG", ... },
  "canalizacion": {
    "resultado": { "tamano": "27", "area_total_mm2": 345.0, ... },
    "fill_factor": 0.35,
    "detalle_tuberia": { ... }
  },
  "proteccion": { "itm": 125 },
  "caida_tension": { "porcentaje": 2.1, "caida_volts": 10.08, "cumple": true, ... },
  "cumple_normativa": true,
  "observaciones": [...],
  "pasos": [...]
}
```

**Grupos de respuesta:**

| Grupo | Descripción |
| ----- | ----------- |
| `equipo` | Datos del filtro/equipo del catálogo |
| `instalacion` | Tensión, sistema eléctrico, canalización, material, longitud, hilos por fase, % caída |
| `corrientes` | Corriente nominal, ajustada, factores de corrección, temperatura, tabla NOM |
| `cable_fase` | Conductor de alimentación (antes `conductor_alimentacion`) |
| `cable_neutro` | Conductor neutro (nil para DELTA) |
| `cable_tierra` | Conductor de tierra |
| `canalizacion` | Resultado, fill factor, detalle charola/tubería |
| `proteccion` | ITM |
| `caida_tension` | Resultados de caída de tensión |

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
