# Shared Kernel

Value objects compartidos entre todas las features.

## Responsabilidades

- Definir tipos de valor reutilizables
- **NO conocer** ninguna feature específica
- Solo tipos puros, sin lógica de negocio compleja

## Value Objects

| Value Object | Descripción |
|-------------|-------------|
| `Corriente` | Valor de corriente eléctrica (Amperes) |
| `Tension` | Valor de tensión eléctrica (Volts) |
| `Temperatura` | Temperatura en °C |
| `MaterialConductor` | Material del conductor (Cu, Al) |
| `Conductor` | Conductor con todas sus propiedades |
| `Charola` | Charola para cables |
| `ResistenciaReactancia` | R y X para cálculo de caída de tensión |
| `TablaEntrada` | Entrada genérica para tablas |

## Reglas

- Inmutables
- Sin dependencias externas
- Validadores en constructor

## Uso

```go
corriente := corriente.New(50.0)
tension := tension.New(220)
material := materialConductor.Cobre()
```
