# Tablas NOM (CSV)

Datos de normativa NOM en formato CSV. Leidos por `infrastructure/repository/`.

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — convenciones al agregar mapeos en infrastructure

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Agregar nueva tabla NOM | `golang-patterns` (para el mapeo en infrastructure) |
| Modificar formato CSV existente | leer este AGENTS.md primero — validar contra fuente original |

## Tablas Disponibles (Fase 1)

### Ampacidad (seleccion de conductor de alimentacion)

| Archivo | Tabla NOM | Canalizacion | Rango |
|---------|-----------|-------------|-------|
| `310-15-b-16.csv` | 310-15(b)(16) | Tuberia conduit (4 tipos) | 14 AWG - 2000 MCM |
| `310-15-b-17.csv` | 310-15(b)(17) | Charola cable espaciado | 14 AWG - 2000 MCM |
| `310-15-b-20.csv` | 310-15(b)(20) | Charola triangular | 8 AWG - 1000 MCM |

### Tierra

| Archivo | Tabla NOM |
|---------|-----------|
| `250-122.csv` | 250-122 |

### Referencia impedancia (caida de tension)

| Archivo | Tabla NOM | Dato |
|---------|-----------|------|
| `tabla-9-resistencia-reactancia.csv` | Tabla 9 | R (ohm/km) por tipo conduit + reactancia |
| `tabla-5-dimensiones-aislamiento.csv` | Tabla 5 | Diametro exterior THW (mm) para DMG |
| `tabla-8-conductor-desnudo.csv` | Tabla 8 | Diametro desnudo (mm) + hilos para RMG |

### Canalizacion (dimensionamiento de tuberia)

| Archivo | Tabla NOM | Dato |
|---------|-----------|------|
| `tabla-conduit-dimensiones.csv` | Capítulo 9, Tabla 4 | Area interior (mm²) por tamaño de conduit |

## Formato CSV Ampacidad

```
seccion_mm2,calibre,cu_60c,cu_75c,cu_90c,al_60c,al_75c,al_90c
```

- Celdas vacias donde no aplica
- `310-15-b-20.csv` no tiene columnas 60C (charola triangular)
- Calibres en formato: `14 AWG`, `4/0 AWG`, `250 MCM`, etc.

## Seleccion de Temperatura

- <= 100A o calibres 14-1 AWG -> columna 60C
- > 100A o calibres > 1 AWG -> columna 75C
- 90C: solo con override explicito (todos los equipos certificados 90C)
- Charola triangular sin 60C -> fallback automatico a 75C

## Regla de Validacion

**NUNCA escribir datos de tablas NOM de memoria o imagenes.**
Siempre pedir datos en Excel/CSV al usuario y validar fila por fila contra la fuente original.

## Agregar nueva tabla (Fase 2+)

1. Obtener datos del usuario en Excel/CSV
2. Crear archivo CSV con formato estandar
3. Validar 100% contra fuente original (fila por fila)
4. Agregar mapeo en `infrastructure/repository/`
5. Documentar en este AGENTS.md

---

## CRITICAL RULES

### Integridad de Datos
- ALWAYS: Validar 100% contra fuente oficial NOM fila por fila antes de commitear
- ALWAYS: Celdas vacias donde la norma no aplica (no rellenar con 0)
- ALWAYS: Cabecera en primera fila, sin filas en blanco al inicio
- NEVER: Escribir o inferir datos NOM de memoria, imagenes o estimaciones
- NEVER: Modificar valores existentes sin validar contra la fuente original

### Formato
- ALWAYS: UTF-8, separador coma, punto decimal (no coma)
- ALWAYS: Calibres en formato estandar: `14 AWG`, `4/0 AWG`, `250 MCM`
- NEVER: Separadores de miles en valores numericos

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Archivo tabla ampacidad | `{articulo-nom}.csv` | `310-15-b-16.csv` |
| Archivo tabla referencia | `tabla-{numero}-{descripcion}.csv` | `tabla-9-resistencia-reactancia.csv` |
| Columna material+temp | `{material}_{temp}c` | `cu_60c`, `al_75c` |
| Columna resistencia | `res_{material}_{conduit}` | `res_cu_pvc`, `res_al_acero` |

---

## QA CHECKLIST

- [ ] Datos validados fila por fila contra fuente NOM original
- [ ] Cabecera con todas las columnas requeridas
- [ ] Sin valores inventados o interpolados
- [ ] Archivo agregado a la tabla "Tablas Disponibles" en este AGENTS.md
- [ ] Mapeo agregado en `infrastructure/repository/CSVTablaNOMRepository`
