# Domain Entities

Entidades del dominio de cálculo eléctrico.

## Entidades

| Entidad | Descripción |
|---------|-------------|
| `MemoriaCalculo` | Agregado raíz: memoria de cálculo completa |
| `TipoCanalizacion` | Tipo de canalización (PVC, EMT, IMC, RMC, etc.) |
| `SistemaElectrico` | Sistema eléctrico (Monofásico, Trifásico Delta, Trifásico Estrella) |
| `ITM` | Interruptor termomagnético |
| `Equipo` | Equipo eléctrico |
| `Carga` | Carga del sistema |
| `FiltroActivo` | Filtro activo de armónicos |
| `FiltroRechazo` | Filtro de rechazo |
| `Transformador` | Transformador |
| `TipoEquipo` | Tipo de equipo |

## Valores

- **NO** tienen lógica de negocio compleja
- **SÍ** tienen validaciones básicas de construcción
- Usan value objects de `shared/kernel/`

## Ejemplo

```go
tipo := entity.NewTipoCanalizacion("TUBERIA_PVC")
```
