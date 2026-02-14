# Garfex Calculadora de Memorias de Cálculo

Backend API en Go para memorias de cálculo de instalaciones eléctricas según normativa NOM (México).

## Estado del Proyecto

✅ **Fase 1 - Completa**  
✅ **Fase 2 - Completa**

Las dos fases principales del proyecto están implementadas y operativas.

---

## Fases Completadas

### ✅ Fase 1: Domain Layer + API Core
- **4 tipos de equipos:** Filtro Activo, Filtro de Rechazo, Transformador, Carga
- **6 servicios de cálculo:**
  - Corriente nominal por tipo de equipo
  - Ajuste de corriente (factores de temperatura y agrupamiento)
  - Selección de conductor de alimentación
  - Selección de conductor de tierra (Cu/Al)
  - Dimensionamiento de canalización (tubería/charola)
  - Cálculo de caída de tensión (IEEE-141/NOM)
- **7 tablas NOM CSV:** Ampacidad, tierra, temperatura, impedancia
- **6 tipos de canalización:** 4 tuberías + 2 charolas
- **API REST** con Gin

**Planes completados:** 11

### ✅ Fase 2: Factores Automáticos
- **Factores de ajuste automáticos** basados en tablas NOM
  - Factor de temperatura según estado de México
  - Factor de agrupamiento según sistema eléctrico
- **Cálculo de canalización mejorado:**
  - Tubería: fill factors 53%/31%/40% según cantidad de conductores
  - Charola espaciada: fórmula NOM con espacios
  - Charola triangular: fórmula NOM para arreglos
- **Soporte multi-tubo:** Canalización en paralelo (`hilos_por_fase`)

---

## Stack Tecnológico

- **Backend:** Go 1.22+
- **Framework:** Gin
- **Base de datos:** PostgreSQL (pgx/v5)
- **Tablas NOM:** CSV
- **Testing:** testify
- **Linting:** golangci-lint

---

## Arquitectura

Arquitectura hexagonal/clean architecture:

```
internal/
├── domain/          # Lógica de negocio pura
│   ├── entity/      # Entidades y tipos
│   ├── valueobject/ # Value objects (inmutables)
│   └── service/     # Servicios de cálculo
├── application/     # Orquestación
│   ├── dto/         # Data Transfer Objects
│   ├── port/        # Interfaces (contratos)
│   └── usecase/     # Casos de uso
├── infrastructure/  # Implementaciones
│   └── repository/  # Repositorios CSV/BD
└── presentation/    # API REST
    └── handler/     # HTTP handlers
```

---

## Uso Rápido

### Iniciar servidor

```bash
# Compilar y ejecutar
go build -o server.exe ./cmd/api/main.go
./server.exe

# O ejecutar directamente
go run cmd/api/main.go
```

### Endpoint principal

```bash
curl -X POST http://localhost:8080/api/v1/calculos/memoria \
  -H "Content-Type: application/json" \
  -d '{
    "modo": "MANUAL_AMPERAJE",
    "amperaje_nominal": 50,
    "tension": 220,
    "tipo_canalizacion": "TUBERIA_PVC",
    "hilos_por_fase": 1,
    "longitud_circuito": 10,
    "itm": 100,
    "factor_potencia": 0.9,
    "estado": "Sonora",
    "sistema_electrico": "DELTA",
    "material": "Cu"
  }'
```

### Verificar salud

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Testing

```bash
# Tests unitarios
go test ./...

# Tests con race detector
go test -race ./...

# Build
go build ./...

# Linting
golangci-lint run
```

---

## Documentación

### Planes de Diseño (completados)

Todos los planes están en `docs/plans/completed/`:

1. Arquitectura inicial
2. Domain layer
3. Nuevos equipos (Transformador + Carga)
4. Tablas NOM canalización
5. Ports CSV infrastructure
6. Caída de tensión IEEE-141
7. Material conductor tierra
8. Canalización multi-tubo
9. Fase 2: Memoria de cálculo completa

### Guías por Capa

Cada capa tiene su propio `AGENTS.md` con guías específicas:

- `internal/domain/AGENTS.md`
- `internal/application/AGENTS.md`
- `internal/infrastructure/AGENTS.md`
- `internal/presentation/AGENTS.md`

---

## Releases

- **v1.0.0** - Fase 1 + Fase 2 Completas (2026-02-13)
  - API REST operativa
  - Cálculo completo de memorias NOM
  - 4 equipos, 6 servicios, 11 tablas CSV
  - Factores automáticos y multi-tubo

---

## Licencia

Proyecto privado - GARFEX
