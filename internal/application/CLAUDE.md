# Application Layer

Orquesta domain services. Define contratos (ports), no implementaciones.

## Estructura

- `port/` — Interfaces que infrastructure implementa
- `usecase/` — CalcularMemoriaUseCase (orquesta los 7 pasos)
- `dto/` — EquipoInput, MemoriaOutput (entrada/salida de la API)

## Ports (interfaces)

- **EquipoRepository** — buscar equipos en BD (PostgreSQL)
- **TablaNOMRepository** — leer tablas CSV de ampacidad, tierra, impedancia

Las interfaces se definen aqui, se implementan en `infrastructure/`.
Pequenas y enfocadas (pocos metodos por interface).

## Flujo del UseCase (orden obligatorio)

1. Corriente Nominal (segun TipoEquipo)
2. Ajuste de Corriente (factores)
3. Seleccionar TipoCanalizacion — determina tabla NOM
4. Resolver tabla ampacidad + columna temperatura — llamar SeleccionarConductorAlimentacion
5. Conductor de Tierra (ITM -> tabla 250-122)
6. Dimensionar Canalizacion (40% fill)
7. Resolver datos Tablas 9/5/8 — llamar CalcularCaidaTension

## Seleccion de Temperatura (logica aqui, no en domain)

- <= 100A -> 60C (o 75C si charola triangular sin columna 60C)
- > 100A -> 75C
- 90C solo con `temperatura_override: 90` explicito del usuario

## DTOs

- **EquipoInput:** modo (LISTADO/MANUAL_AMPERAJE/MANUAL_POTENCIA), datos del equipo, parametros de instalacion, TipoCanalizacion, TemperaturaOverride
- **MemoriaOutput:** resultado completo de todos los pasos para el reporte

## Convenciones

- `context.Context` como primer parametro en operaciones I/O
- Errores de flujo: `ErrEquipoNoEncontrado`, `ErrModoInvalido`
- DTOs son structs planos sin logica de negocio
- Nunca importar infrastructure — solo domain
