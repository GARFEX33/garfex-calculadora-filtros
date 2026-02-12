# Tablas NOM (CSV)

Datos de normativa NOM en formato CSV. Leidos por `infrastructure/repository/`.

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
5. Documentar en este CLAUDE.md
