# Domain Layer

Capa de negocio pura. Sin dependencias externas (sin Gin, pgx, CSV).
Recibe datos ya interpretados — no conoce archivos ni BD.
Para patrones Go idiomaticos, usa skill `golang-patterns`.

## Estructura

- `entity/` — Equipos, ITM, MemoriaCalculo, TipoCanalizacion, TipoEquipo
- `valueobject/` — Corriente, Tension, Conductor (inmutables)
- `service/` — 6 servicios de calculo electrico

## Tipos de Equipos (CalculadorCorriente)

| Tipo | Parametro | Formula In |
|------|-----------|-----------|
| FiltroActivo | AmperajeNominal | In = directo |
| FiltroRechazo | KVAR | In = KVAR / (KV x sqrt3) |
| Transformador | KVA | In = KVA / (KV x sqrt3) |
| Carga | KW + FP + Fases | In = KW / (KV x factor x FP) |

Todos implementan la interfaz `CalculadorCorriente`. Carga y Transformador tambien implementan `CalculadorPotencia`.

## Value Objects

- **Corriente:** valida > 0, unidad "A", metodos `Multiplicar`/`Dividir`
- **Tension:** solo valores NOM: 127, 220, 240, 277, 440, 480, 600
- **Conductor:** `NewConductor(ConductorParams{})`, campos requeridos: Calibre, Material, SeccionMM2. TipoAislamiento vacio para conductores desnudos (tierra)

## TipoCanalizacion (6 valores)

```
TUBERIA_PVC, TUBERIA_ALUMINIO, TUBERIA_ACERO_PG, TUBERIA_ACERO_PD,
CHAROLA_CABLE_ESPACIADO, CHAROLA_CABLE_TRIANGULAR
```

Los 4 tipos de tuberia comparten tabla de ampacidad. Cada tipo mapea a columna R diferente en Tabla 9.

## TipoEquipo (4 valores)

```
FILTRO_ACTIVO, FILTRO_RECHAZO, TRANSFORMADOR, CARGA
```

## Servicios

1. **CorrienteNominal** — calcula In segun tipo de equipo
2. **AjusteCorriente** — aplica factores (temperatura, agrupamiento)
3. **SeleccionarConductorAlimentacion** — recibe `[]EntradaTablaConductor` pre-resueltos
4. **SeleccionarConductorTierra** — ITM + tabla tierra pre-resuelta
5. **CalcularCanalizacion** — NOM 40% fill para tuberia
6. **CalcularCaidaTension** — metodo impedancia: Z=sqrt(R^2+X^2), VD=sqrt3 x I x Z x L_km

## Caida de Tension (metodo impedancia)

- **R:** recibida de Tabla 9 (pre-resuelta por infrastructure)
- **X:** calculada geometricamente con DMG/RMG
- **RMG** = (diametro_desnudo/2) x factorHilos[numHilos]
- **DMG** = diametro_exterior_thw x factorDMG[tipoCanalizacion]
- Factores hilos: 1=0.7788, 7=0.726, 19=0.758, 37=0.768, 61=0.772
- Factores DMG: tuberia/triangular=1.0, espaciado=2.0
- Diseno completo: `docs/plans/2026-02-11-caida-tension-impedancia-design.md`

## MemoriaCalculo

Struct resultado que agrupa todos los pasos: CorrienteNominal, CorrienteAjustada, FactoresAjuste, Potencias (KVA/KW/KVAR), ConductorAlimentacion, ConductorTierra, TipoCanalizacion, Canalizacion, TemperaturaUsada, CaidaTension (ResultadoCaidaTension), CumpleNormativa.

## Convenciones

- Constructores: `NewXxx(XxxParams{})` para structs con muchos parametros
- Errores: `ErrXxx = errors.New(...)`, wrap con `fmt.Errorf("%w: ...", ErrXxx)`
- Tests: table-driven con `t.Run`, testify, mismo directorio (`_test.go`)
- Sin panic, sin context en structs, receptores consistentes (valor o puntero, no mezclar)

## Patron para agregar nuevo equipo (Fase 2+)

1. Crear entity en `entity/xxx.go` implementando `CalculadorCorriente`
2. Agregar constante en `TipoEquipo`
3. Agregar test table-driven en `entity/xxx_test.go`
4. Los servicios existentes funcionan sin cambio (reciben interfaces)
