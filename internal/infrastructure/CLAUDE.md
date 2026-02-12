# Infrastructure Layer

Implementa los ports definidos en `application/port/`.
Tecnologias: PostgreSQL (pgx/v5), CSV (encoding/csv).

## Estructura

- `repository/` — PostgresEquipoRepository, CSVTablaNOMRepository
- `client/` — PostgresClient (pgx pool)

## PostgresEquipoRepository

- Conecta a Supabase (PostgreSQL via pgx/v5)
- Tabla `equipos_filtros`: clave, tipo, voltaje, "qn/In", itm, bornes
- Mapea "qn/In" segun tipo: ACTIVO -> Amperaje, RECHAZO -> KVAR
- `context.Context` en todas las queries

## CSVTablaNOMRepository

Lee tablas NOM desde `data/tablas_nom/`. Los mapeos de canalizacion a tabla/columna viven aqui.

### Mapeo canalizacion -> tabla ampacidad

| TipoCanalizacion | Archivo CSV |
|---|---|
| TUBERIA_PVC / ALUMINIO / ACERO_PG / ACERO_PD | 310-15-b-16.csv |
| CHAROLA_CABLE_ESPACIADO | 310-15-b-17.csv |
| CHAROLA_CABLE_TRIANGULAR | 310-15-b-20.csv |

### Mapeo canalizacion -> columna R (Tabla 9)

| TipoCanalizacion | Columna resistencia |
|---|---|
| TUBERIA_PVC | res_{material}_pvc |
| TUBERIA_ALUMINIO | res_{material}_al |
| TUBERIA_ACERO_PG / ACERO_PD | res_{material}_acero |
| CHAROLA_CABLE_ESPACIADO | res_{material}_pvc |
| CHAROLA_CABLE_TRIANGULAR | res_{material}_pvc |

Charola no tiene conduit metalico, usa columna PVC (sin efecto de proximidad).

### Tablas de referencia impedancia

- `tabla-9-resistencia-reactancia.csv` — R (ohm/km) por calibre + material conduit
- `tabla-5-dimensiones-aislamiento.csv` — diametro exterior THW para DMG
- `tabla-8-conductor-desnudo.csv` — diametro desnudo + hilos para RMG

## Entorno de Desarrollo

### Variables de Entorno

```bash
DB_HOST=192.168.1.X   # dev: IP servidor Ubuntu | prod: localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=...
DB_NAME=postgres
ENVIRONMENT=development|production
```

### Topologia de Red

```
Servidor Ubuntu → Supabase (Docker) → PostgreSQL :5432 → Cloudflare Tunnel
Laptop Windows (dev) → conecta a 192.168.1.X:5432
```

## Convenciones

- Nunca importar `domain/service` — solo `domain/entity` y `domain/valueobject`
- `context.Context` como primer parametro siempre
- Errores wrapped con contexto: `fmt.Errorf("leer tabla %s: %w", nombre, err)`
- Conexion BD: inyeccion manual en `cmd/api/main.go`, sin globals
- Sin estado global mutable (`var db *pgxpool.Pool` a nivel de package)
