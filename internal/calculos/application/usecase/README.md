# Use Cases

Casos de uso que orquestan la lógica de dominio.

## Use Cases

| Use Case | Descripción |
|----------|-------------|
| `OrquestadorMemoriaCalculo` | Orquesta todo el flujo de memoria de cálculo |
| `CalcularMemoria` | Calcula memoria completa |
| `CalcularCorriente` | Calcula corriente nominal |
| `AjustarCorriente` | Ajusta por temperatura/agrupamiento |
| `DimensionarCanalizacion` | Dimensiona canalización |
| `CalcularCaidaTension` | Calcula caída de tensión |
| `SeleccionarConductor` | Selecciona conductor |
| `SeleccionarTemperatura` | Selecciona temperatura |

## Estructura

```
usecase/
├── orquestador_memoria.go
├── calcular_memoria.go
├── calcular_corriente.go
├── ajustar_corriente.go
├── dimensionar_canalizacion.go
├── calcular_caida_tension.go
├── seleccionar_conductor.go
└── helpers/              # Funciones auxiliares
```

## Reglas

- Solo orquestación (sin lógica de negocio)
- < 80 líneas por caso de uso
- Error wrapping con `%w`
