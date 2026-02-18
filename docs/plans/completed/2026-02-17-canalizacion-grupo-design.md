# Diseño: Endpoint Canalización de Grupo de Alimentadores

**Fecha:** 2026-02-17  
**Estado:** Aprobado  
**Autor:** Orquestador + Usuario

---

## Problema

En instalaciones eléctricas, un circuito puede tener múltiples hilos por fase (1, 2, 3, etc.). Cuando hay varios hilos por fase:

- **Cables pequeños** — Podrían caber todos en un solo tubo
- **Cables grandes** — Es mejor dividir en múltiples tubos (ej: 2 hilos/fase = 2 tubos)

El sistema actual (`DimensionarCanalizacionUseCase`) calcula para un solo alimentador. Se necesita un endpoint dedicado para calcular la canalización de un grupo de alimentadores con número de tubos configurable.

---

## Solución

Crear un endpoint nuevo `POST /api/v1/calculos/canalizacion-grupo` que:

1. Recibe lista de conductores (fases) + sección de tierra + número de tubos
2. Aplica regla NOM: **1 conductor de tierra por tubo**
3. Reutiliza el servicio de dominio `CalcularCanalizacion` existente
4. Retorna tamaño de tubería, áreas y factor de relleno

---

## Arquitectura

### Componentes

```
internal/calculos/
├── domain/
│   ├── entity/canalizacion.go          ← MODIFICAR: agregar FactorRelleno
│   └── service/calculo_canalizacion.go ← MODIFICAR: asignar FactorRelleno
├── application/
│   ├── dto/canalizacion_grupo.go       ← NUEVO
│   └── usecase/calcular_canalizacion_grupo.go ← NUEVO
└── infrastructure/adapter/driver/http/
    └── calculo_handler.go              ← MODIFICAR: agregar handler
cmd/api/main.go                         ← MODIFICAR: wiring
```

### Flujo de datos

```
HTTP Request (JSON primitivos)
       │
       ▼
┌─────────────────────────────┐
│   CanalizacionGrupoInput    │
│   - conductores []          │
│   - seccion_tierra_mm2      │
│   - tipo_canalizacion       │
│   - numero_de_tubos         │
└─────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│  CalcularCanalizacionGrupoUseCase   │
│  1. Validar DTO                     │
│  2. Convertir → ConductorParaCanal. │
│  3. Agregar tierras (1 por tubo)    │
│  4. Obtener tabla NOM (port)        │
│  5. Llamar servicio dominio         │
│  6. Convertir resultado → DTO       │
└─────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────┐
│  CanalizacionGrupoOutput    │
│   - tamano                  │
│   - area_total_mm2          │
│   - area_por_tubo_mm2       │
│   - numero_de_tubos         │
│   - factor_relleno          │
└─────────────────────────────┘
       │
       ▼
HTTP Response (JSON)
```

---

## DTOs

### Input

```go
type ConductorGrupoInput struct {
    Cantidad   int     `json:"cantidad" binding:"required,min=1"`
    SeccionMM2 float64 `json:"seccion_mm2" binding:"required,gt=0"`
}

type CanalizacionGrupoInput struct {
    Conductores      []ConductorGrupoInput `json:"conductores" binding:"required,min=1,dive"`
    SeccionTierraMM2 float64               `json:"seccion_tierra_mm2" binding:"required,gt=0"`
    TipoCanalizacion string                `json:"tipo_canalizacion" binding:"required"`
    NumeroDeTubos    int                   `json:"numero_de_tubos" binding:"required,min=1"`
}
```

### Output

```go
type CanalizacionGrupoOutput struct {
    Tamano         string  `json:"tamano"`
    AreaTotalMM2   float64 `json:"area_total_mm2"`
    AreaPorTuboMM2 float64 `json:"area_por_tubo_mm2"`
    NumeroDeTubos  int     `json:"numero_de_tubos"`
    FactorRelleno  float64 `json:"factor_relleno"`
}
```

---

## Regla de Negocio Clave

> **1 conductor de tierra por tubo** — El sistema agrega automáticamente `numero_de_tubos` conductores de tierra al cálculo.

Ejemplo con 2 tubos:
- Input: 6 fases + 1 tierra (sección)
- Cálculo interno: 6 fases + **2 tierras** (1 por tubo)

---

## Modificación al Dominio

La entidad `Canalizacion` debe incluir `FactorRelleno`:

```go
type Canalizacion struct {
    Tipo           TipoCanalizacion
    Tamano         string
    AnchoRequerido float64
    NumeroDeTubos  int
    FactorRelleno  float64  // ← AGREGAR
}
```

El servicio `CalcularCanalizacion` asigna el factor al construir el resultado.

---

## Endpoint

### Request

```http
POST /api/v1/calculos/canalizacion-grupo
Content-Type: application/json

{
  "conductores": [
    {"cantidad": 6, "seccion_mm2": 53.49}
  ],
  "seccion_tierra_mm2": 21.15,
  "tipo_canalizacion": "TUBERIA_PVC",
  "numero_de_tubos": 2
}
```

### Response

```json
{
  "tamano": "1 1/4",
  "area_total_mm2": 363.24,
  "area_por_tubo_mm2": 181.62,
  "numero_de_tubos": 2,
  "factor_relleno": 0.31
}
```

---

## Enfoque Seleccionado

**Use Case Nuevo + Reutilizar Servicio Existente**

- Reutiliza 100% del dominio existente (`CalcularCanalizacion`)
- Respeta arquitectura hexagonal
- Mínimo código nuevo
- Único cambio en domain: agregar campo `FactorRelleno`

---

## Criterios de Éxito

- [ ] Endpoint responde correctamente con distintas combinaciones de conductores/tubos
- [ ] Factor de relleno se calcula según NOM (0.53, 0.31, o 0.40)
- [ ] Tierra se agrega automáticamente (1 por tubo)
- [ ] Tests unitarios para use case
- [ ] Tests de integración para endpoint
- [ ] `go test ./...` pasa
- [ ] `go build ./...` compila

---

## Referencias

- Servicio existente: `internal/calculos/domain/service/calculo_canalizacion.go`
- Entidad: `internal/calculos/domain/entity/canalizacion.go`
- Use case similar: `internal/calculos/application/usecase/dimensionar_canalizacion.go`
