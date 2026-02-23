# Feature: calculos

Memoria de cálculo eléctrico según normativa NOM (México).

Esta feature implementa el cálculo completo de una instalación eléctrica:
corriente nominal → ajuste por temperatura/agrupamiento → selección de conductor →
conductor de tierra → dimensionamiento de canalización → caída de tensión (NOM).

### Fórmula de caída de tensión (NOM / IEEE-141)

La caída de tensión usa **impedancia efectiva**, no la magnitud vectorial:

```
Zef = R·cosθ + X·senθ      (senθ = √(1 - cos²θ))
e   = factor × (I/N) × L × Zef
%e  = (e / V_referencia) × 100
```

Donde:
- **I** = Corriente nominal (A)
- **N** = Número de conductores en paralelo por fase
- **L** = Longitud del circuito (km)
- **R, X** = Resistencia y reactancia por conductor (Ω/km) de Tabla 9 NOM
- **cosθ** = Factor de potencia (0 < FP ≤ 1)
- **Zef** = Impedancia efectiva por conductor

#### Factores por sistema eléctrico

| Sistema         | factor | Voltaje referencia |
| --------------- | ------ | ------------------ |
| MONOFASICO 1F2H | 2      | Vfn                |
| BIFASICO 2F3H  | 2      | Vfn                |
| DELTA 3F3H     | √3     | Vff                |
| ESTRELLA 3F4H  | √3     | Vfn                |

#### Conversión de voltaje

- Si el sistema requiere Vfn pero el usuario ingresa Vff: `Vref = Vff / √3`
- Si el sistema requiere Vff pero el usuario ingresa Vfn: `Vref = Vfn × √3`

> **NO usar** `Z = √(R² + X²)` — eso es la magnitud, no la impedancia efectiva.

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
└── infrastructure/  ← adapters HTTP (driver) y CSV (driven)
    └── adapter/
        ├── driver/http/     ← CalculoHandler, formatters
        └── driven/csv/      ← CSVTablaNOMRepository
```

## Cómo modificar esta feature

Trabajar directamente en las capas internas:

- `domain/` — entidades y servicios de cálculo puro
- `application/` — ports, use cases, DTOs
- `infrastructure/` — adapters HTTP y CSV

> Ver estructura y reglas en [docs/reference/structure.md](../../../docs/reference/structure.md)

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
