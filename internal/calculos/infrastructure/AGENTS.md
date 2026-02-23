# Calculos — Infrastructure Layer

Implementa los ports definidos en `application/port/`. Tecnologías: CSV (encoding/csv), HTTP (Gin).

## Estructura

```
internal/calculos/infrastructure/
├── adapter/
│   ├── driven/
│   │   ├── csv/              # CSVTablaNOMRepository (tablas NOM)
│   │   └── postgres/         # CalcEquipoFiltroRepository (equipos_filtros)
│   └── driver/
│       └── http/
│           ├── formatters/   # NombreTablaAmpacidad, GenerarObservaciones
│           ├── middleware/   # CORS, RequestLogger
│           └── handler.go    # CalculoHandler
└── router.go                 # Configuración de rutas Gin
```

> **Nota:** Las subcarpetas `adapter/driver/http/` y `adapter/driven/csv/` heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/calculos/domain/entity`
- `internal/calculos/application/port` (interfaces a implementar)
- `internal/calculos/application/usecase` (para llamar desde handlers)
- Gin, encoding/csv

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../docs/reference/structure.md)

## Adapters

### Driven (implementan ports)

- **CSVTablaNOMRepository** — lee tablas NOM desde CSV
- **CSVSeleccionarTemperatura** — temperaturas por estado
- **CalcEquipoFiltroRepository** — consulta `equipos_filtros` en PostgreSQL y mapea a entidades de cálculo. Implementa `calculos/application/port.EquipoRepository`.

### CalcEquipoFiltroRepository — Mapeo de TipoFiltro a entidad

El adapter consulta la tabla `equipos_filtros` y convierte cada registro a la entidad de dominio correcta:

| `TipoFiltro` BD | Campo `Amperaje` | Entidad calculos | Fórmula corriente |
|---|---|---|---|
| `A` | Corriente directa (A) | `FiltroActivo` | `I = Amperaje` (directo) |
| `KVA` | Potencia aparente (KVA) | `Transformador` | `I = KVA / (kV × √3)` |
| `KVAR` | Potencia reactiva (KVAR) | `FiltroRechazo` | `I = KVAR / (kV × √3)` |

> **Nota**: El pool de PostgreSQL se comparte con el repositorio de equipos (`equipos/application/port.EquipoFiltroRepository`).

### Driver (HTTP)

- **CalculoHandler** — endpoints REST
  - `POST /api/v1/calculos/amperaje` — calcular amperaje nominal
  - `POST /api/v1/calculos/corriente-ajustada` — calcular corriente ajustada con factores NOM
  - `POST /api/v1/calculos/conductor-alimentacion` — seleccionar conductor de alimentación
  - `POST /api/v1/calculos/conductor-tierra` — seleccionar conductor de tierra
  - `POST /api/v1/calculos/tuberia` — dimensionar tubería
  - `POST /api/v1/calculos/charola/espaciado` — calcular espaciado en charola
  - `POST /api/v1/calculos/charola/triangular` — calcular configuración triangular
  - `POST /api/v1/calculos/caida-tension` — calcular caída de tensión

- **MemoriaHandler** — orquestador completo
  - `POST /api/v1/calculos/memoria` — memoria de cálculo completa (orquesta todos los pasos)

#### Body: POST /api/v1/calculos/caida-tension

```json
{
  "calibre":           "2 AWG",
  "material":          "Cu",
  "tipo_canalizacion": "TUBERIA_PVC",
  "corriente_nominal": 70.0,
  "longitud_circuito": 30.0,
  "tension":           127,
  "sistema_electrico": "MONOFASICO",
  "tipo_voltaje":      "FASE_NEUTRO",
  "hilos_por_fase":    1,
  "limite_caida":      3.0,
  "factor_potencia":   0.9
}
```

> `factor_potencia` (cosθ) es **required**. Rango: `(0, 1]`. Sin él → HTTP 400.
> El campo `impedancia` en la respuesta es `Zef = R·cosθ + X·senθ` por conductor (Ω/km).
> La caída se calcula como: `e = factor × (I/N) × L × Zef`

#### Body: POST /api/v1/calculos/memoria

```json
{
  "modo": "MANUAL_POTENCIA",
  "tipo_equipo": "CARGA",
  "potencia_nominal": 15,
  "potencia_unidad": "KW",
  "tension": 480,
  "tension_unidad": "V",
  "factor_potencia": 0.85,
  "itm": 30,
  "sistema_electrico": "ESTRELLA",
  "estado": "Ciudad de Mexico",
  "tipo_canalizacion": "TUBERIA_PVC",
  "longitud_circuito": 50,
  "tipo_voltaje": "FASE_FASE"
}
```

> `potencia_unidad` acepta: W, KW, KVA, KVAR (default: KW)
> `tension_unidad` acepta: V, kV (default: V). Con `kV`, el valor se normaliza internamente: `0.48 kV` → `480 V`
> `tipo_voltaje` acepta: FASE_NEUTRO, FASE_FASE

**Ejemplo con modo LISTADO (equipo de la BD):**
```json
{
  "modo": "LISTADO",
  "clave": "GAR-100A-480V",
  "itm": 125,
  "tension": 480,
  "sistema_electrico": "ESTRELLA",
  "estado": "Ciudad de Mexico",
  "tipo_canalizacion": "TUBERIA_PVC",
  "longitud_circuito": 50,
  "tipo_voltaje": "FASE_FASE"
}
```
> En modo `LISTADO`, los campos `amperaje_nominal`, `potencia_nominal` y `tipo_equipo` se obtienen automáticamente de la BD según la `clave`.

**Ejemplo con kV:**
```json
{
  "tension": 0.48,
  "tension_unidad": "kV"
}
```

## Mapeo de Errores HTTP — Tensión

| Error                    | HTTP | Causa                              |
| ------------------------ | ---- | ---------------------------------- |
| `ErrVoltajeInvalido`     | 400  | Valor no está en lista NOM         |
| `ErrUnidadTensionInvalida` | 400 | Unidad no es "V" ni "kV"          |

### Formatters

- **NombreTablaAmpacidad** — nombres descriptivos de tablas
- **GenerarObservaciones** — observaciones del cálculo

## Mapeo de Errores HTTP

| Error domain/application    | HTTP status |
| --------------------------- | ----------- |
| ErrModoInvalido             | 400         |
| ErrTipoCanalizacionInvalido | 400         |
| ErrSistemaElectricoInvalido | 400         |
| ErrTipoEquipoInvalido       | 400         |
| Validación                  | 400         |
| ErrConductorNoEncontrado    | 422         |
| ErrCanalizacionNoDisponible | 422         |
| CALCULO_NO_POSIBLE          | 422         |
| Error interno               | 500         |

## Reglas de Oro — Capa Infrastructure

*Estas reglas son específicas para la capa Infrastructure de cálculos. Ver [docs/reference/structure.md](../../../docs/reference/structure.md) para reglas globales.*

1. **Implementar exactamente el port** — no agregar métodos
2. **Sin lógica de negocio** — solo traducción de datos
3. **Handlers solo coordinan** — bind → use case → response
4. **Inyección de dependencias** — constructor, no globals
5. **Context.Context** — primer parámetro en I/O

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)

## QA Checklist

- [ ] `go test ./internal/calculos/infrastructure/...` pasa
- [ ] Repositorios implementan ports exactamente
- [ ] Sin estado global
- [ ] Sin lógica de negocio
