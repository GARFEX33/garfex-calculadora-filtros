# Comandos y Endpoints

## Stack

Go 1.22+, Gin, PostgreSQL (pgx/v5), testify, golangci-lint

## Comandos de Desarrollo

```bash
go test ./...           # Tests
go test -race ./...     # Tests con race detector
go build ./...          # Compilacion
go vet ./...            # Analisis estatico
golangci-lint run       # Linting completo
```

## Iniciar Servidor

**IMPORTANTE:** Asegurarse de que el puerto 8080 esté libre antes de iniciar:

```bash
# Opción 1: Compilar y ejecutar (recomendado)
go build -o server.exe ./cmd/api/main.go
./server.exe

# Opción 2: Ejecutar directamente (sin compilar)
go run cmd/api/main.go

# Opción 3: Puerto personalizado (si 8080 está ocupado)
set PORT=8090
go run cmd/api/main.go
```

**Verificar que el servidor está corriendo:**

```bash
curl http://localhost:8080/health
# Respuesta esperada: {"status":"ok"}
```

## Endpoints Disponibles

### Amperaje Nominal
```bash
curl -X POST http://localhost:8080/api/v1/calculos/amperaje \
  -H "Content-Type: application/json" \
  -d '{"potencia_watts":5000,"tension":220,"sistema_electrico":"MONOFASICO","factor_potencia":0.9}'
```

### Corriente Ajustada
```bash
curl -X POST http://localhost:8080/api/v1/calculos/corriente-ajustada \
  -H "Content-Type: application/json" \
  -d '{"amperaje_nominal":50,"estado":"Sonora","temperatura_ambiente":30,"cantidad_conductores":3,"factor_uso":1.0}'
```

### Conductor Alimentación
```bash
curl -X POST http://localhost:8080/api/v1/calculos/conductor-alimentacion \
  -H "Content-Type: application/json" \
  -d '{"corriente_ajustada":60,"tipo_canalizacion":"TUBERIA_PVC","material":"Cu"}'
```

### Conductor Tierra
```bash
curl -X POST http://localhost:8080/api/v1/calculos/conductor-tierra \
  -H "Content-Type: application/json" \
  -d '{"corriente_ajustada":60,"tipo_canalizacion":"TUBERIA_PVC","material":"Cu"}'
```

### Tuberia
```bash
curl -X POST http://localhost:8080/api/v1/calculos/tuberia \
  -H "Content-Type: application/json" \
  -d '{"conductor_seccion_mm2":13.3,"hilos_por_fase":1}'
```

### Caida de Tension
```bash
curl -X POST http://localhost:8080/api/v1/calculos/caida-tension \
  -H "Content-Type: application/json" \
  -d '{"calibre":"2 AWG","material":"Cu","corriente_ajustada":50,"longitud_circuito":100,"tension":220,"limite_caida":3,"tipo_canalizacion":"TUBERIA_PVC","sistema_electrico":"DELTA","tipo_voltaje":"FF","hilos_por_fase":1}'
```

### Charola
```bash
# Espaciado
curl -X POST http://localhost:8080/api/v1/calculos/charola/espaciado \
  -H "Content-Type: application/json" \
  -d '{"ancho_mm":600,"cantidad_conductores":5}'

# Triangular
curl -X POST http://localhost:8080/api/v1/calculos/charola/triangular \
  -H "Content-Type: application/json" \
  -d '{"ancho_mm":600,"cantidad_conductores":5}'
```
