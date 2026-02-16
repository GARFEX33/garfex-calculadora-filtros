# DTOs

Data Transfer Objects para entrada y salida de la API.

## DTOs de Entrada

| DTO | Descripción |
|-----|-------------|
| `EquipoInput` | Datos de entrada para cálculo de equipo |

## DTOs de Salida

| DTO | Descripción |
|-----|-------------|
| `MemoriaOutput` | Memoria de cálculo completa |

## Reglas

- Structs planos sin métodos complejos
- Tags JSON para serialización
- **NO exponer** entidades de dominio directamente

## Conversión

La conversión dominio ↔ DTO se hace en los use cases:

```go
output := memoriaOutput.FromDomain(memoria)
```
