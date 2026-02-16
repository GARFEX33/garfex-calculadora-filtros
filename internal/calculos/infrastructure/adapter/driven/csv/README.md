# CSV Repository

Implementación de repositorios usando archivos CSV.

## Repositories

| Repository | Descripción |
|------------|-------------|
| `CSVTablaNOMRepository` | Acceso a tablas NOM desde CSV |
| `SeleccionarTemperaturaRepository` | Selección de temperatura por estado |

## Archivos CSV

```
data/tablas_nom/
├── 310-15-b-2-a.csv     # Tabla de ampacidad
├── 310-15-b-3-a.csv     # Factores de temperatura
├── 250-122.csv          # Conductor de tierra
├── tabla-9-resistencia-reactancia.csv
├── tabla-conduit-dimensiones.csv
└── ...
```

## Uso

```go
repo, err := csv.NewCSVTablaNOMRepository("data/tablas_nom")
datos, err := repo.ObtenerTablaAmpacidad(ctx, canalizacion, material, temp)
```
