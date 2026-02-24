# Design: Número de Hilos de Tierra

## Technical Approach

El cambio implementa el cálculo automático del número de hilos de tierra según la normativa NOM, eliminando el valor hardcodeado `1` en el cálculo de dimensionamiento de canalización.

**Estrategia**: Agregar una función helper en el orquestador que calcule el número de tierras según las reglas de negocio y pasar ese valor al use case de tubería. El DTO `TuberiaInput` no necesita modificaciones ya que el campo no se expone al usuario.

## Architecture Decisions

### Decision: Dónde implementar la lógica de cálculo de hilos de tierra

**Choice**: En el `OrquestadorMemoriaCalculoUseCase`, antes de construir el `TuberiaInput`

**Alternatives considered**:
- En `CalcularTamanioTuberiaUseCase` — rechazado porque violaría el principio de "use cases solo orquestan"
- Agregar campo `NumTierras` al DTO `TuberiaInput` — rechazado porque no se expone al usuario

**Rationale**: 
- El orquestador es el responsable de construir los inputs para los use cases
- La lógica de decisión (charola vs tubería, número de tubos) ya está disponible en el orquestador
- Mantiene el use case simple y centrado en su responsabilidad (invocar servicio de dominio)

### Decision: Validación de número de tubos

**Choice**: Usar默认值 1 cuando `NumTuberias` no esté especificado o sea inválido

**Alternatives considered**:
- Retornar error inmediatamente — rechazado por compatibilidad hacia atrás
- Calcular como si fuera 1 tubo silently — chosen

**Rationale**: 
- Es consistente con el comportamiento existente
- El número de tubos es un campo obligatorio en `TuberiaInput` (`binding:"required,gt=0"`), pero el defaulting ocurre antes

## Data Flow

```
Usuario → OrquestadorMemoriaCalculoUseCase
                │
                ├─ Paso 1-3: Corriente, Ajuste, Conductores
                │
                ├─ Paso 4: Dimensionamiento Canalización
                │     │
                │     ├─ Si Charola:
                │     │     └─ NumTierras = 1 (hardcodeado por naturaleza)
                │     │
                │     └─ Si Tubería:
                │           ├─ Obtener NumTuberias del input
                │           ├─ CalcularNumHilosTierra(tipo, numTuberias)
                │           │     │
                │           │     ├─ Si tipo = Charola → return 1
                │           │     ├─ Si tipo = Tubería AND numTuberias ≤ 2 → return 1
                │           │     └─ Si tipo = Tubería AND numTuberias > 2 → return 2
                │           │
                │           └─ Crear TuberiaInput con NumTierras calculado
                │                 └─ Ejecutar CalcularTamanioTuberiaUseCase
                │
                └─ Paso 5: Caída de Tensión
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/calculos/application/usecase/orquestador_memoria_calculo.go` | Modify | Agregar función helper `calcularNumHilosTierra()` y llamarla antes de construir `TuberiaInput` en paso 4 |
| `internal/calculos/application/usecase/calcular_tamanio_tuberia.go` | Modify | Eliminar hardcode `1` en línea 76, usar `input.NumTierras` (requiere agregar campo al DTO) |
| `internal/calculos/application/dto/tuberia_input.go` | Modify | Agregar campo `NumTierras int` al struct (no se expone a usuario, solo uso interno) |

**Nota**: El DTO necesita el campo `NumTierras` para que el use case pueda recibir el valor calculado. Se agregará con tag `json:"-"` para no exponerse en JSON.

## Interfaces / Contracts

### Nuevo campo en TuberiaInput

```go
// internal/calculos/application/dto/tuberia_input.go
type TuberiaInput struct {
    NumFases         int    `json:"num_fases" binding:"required,gt=0"`
    CalibreFase      string `json:"calibre_fase" binding:"required"`
    NumNeutros       int    `json:"num_neutros" binding:"gte=0"`
    CalibreNeutro    string `json:"calibre_neutral"`
    CalibreTierra    string `json:"calibre_tierra" binding:"required"`
    TipoCanalizacion string `json:"tipo_canalizacion" binding:"required"`
    NumTuberias      int    `json:"num_tuberias" binding:"required,gt=0"`
    // Nuevo campo - no se expone al usuario
    NumTierras       int    `json:"-"`
}
```

### Función helper en orquestador

```go
// internal/calculos/application/usecase/orquestador_memoria_calculo.go

// calcularNumHilosTierra calcula el número de hilos de tierra según NOM.
// Reglas:
// - Charola: siempre 1 hilo
// - Tubería: si tubes <= 2 → 1 hilo; si tubes > 2 → 2 hilos
func calcularNumHilosTierra(tipoCanalizacion entity.TipoCanalizacion, numTuberias int) int {
    // Default a 1 tubo si es inválido
    if numTuberias <= 0 {
        numTuberias = 1
    }
    
    // Charola siempre usa 1 hilo de tierra
    if tipoCanalizacion.EsCharola() {
        return 1
    }
    
    // Tubería: más de 2 tubos = 2 hilos de tierra
    if numTuberias > 2 {
        return 2
    }
    
    return 1
}
```

### Modificación en use case de tubería

```go
// internal/calculos/application/usecase/calcular_tamanio_tuberia.go (línea 76)

// Antes:
1, // 1 tierra por tubo

// Después:
input.NumTierras, // Valor calculado por el orquestador
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Función `calcularNumHilosTierra` | Casos: charola (1), tubería 1-2 tubes (1), tubería >2 tubes (2), tubes=0 (default 1), tubes negativo (default 1) |
| Unit | Use case `CalcularTamanioTuberiaUseCase` con NumTierras variable | Verificar que se pase el valor correcto al dominio |
| Integration | Orquestador completo con tubería 1-2-3+ tubes | Verificar resultado de canalización con diferentes numTuberias |

**Casos de prueba específicos**:
1. Tubería con 1 tubo → 1 hilo de tierra
2. Tubería con 2 tubos → 1 hilo de tierra  
3. Tubería con 3 tubos → 2 hilos de tierra
4. Tubería con 4+ tubos → 2 hilos de tierra
5. Charola cable espaciado → 1 hilo de tierra
6. Charola cable triangular → 1 hilo de tierra

## Migration / Rollout

No se requiere migración de datos. El cambio es backwards-compatible porque:
- El cálculo anterior (siempre 1 hilo) se mantiene para los casos existentes (1-2 tubos)
- Solo afecta cuando hay más de 2 tubos, que anteriormente generaba un resultado incorrecto

**Feature flag**: No requerido — el comportamiento anterior era incorrecto y esto lo corrige.

## Open Questions

Ninguna. El diseño está completo y las specs son claras sobre los requisitos.
