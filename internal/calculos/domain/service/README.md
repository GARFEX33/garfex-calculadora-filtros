# Domain Services

Servicios de cálculo puro según normativa NOM y IEEE-141.

## Servicios

| Servicio | Norma | Descripción |
|----------|-------|-------------|
| `CalcularCorrienteNominal` | NOM | Calcula corriente nominal I = P / (V × FP) |
| `AjusteCorriente` | NOM-001 | Ajusta corriente por temperatura y agrupamiento |
| `CalcularConductor` | NOM-001 | Selecciona calibre de conductor |
| `CalculoCanalizacion` | NOM | Dimensiona tubo/conduit |
| `CalculoTierra` | NOM-001 | Calcula conductor de puesta a tierra |
| `CalculoCaidaTension` | IEEE-141 | Calcula caída de tensión con impedancia |
| `SeleccionarTemperatura` | NOM | Selecciona temperatura según estado |
| `CalcularFactorTemperatura` | NOM | Factor de corrección por temperatura |
| `CalcularFactorAgrupamiento` | NOM | Factor de corrección por agrupamiento |
| `CalcularCharolaEspaciado` | NOM | Espaciamiento en charolas |
| `CalcularCharolaTriangular` | NOM | Arreglo triangular de conductores |
| `CalcularFactorUso` | NOM | Factor de utilización |

## Reglas

- **Puro**: sin imports externos a stdlib
- **Sin contexto**: no recibe `context.Context`
- **Errores de negocio**: retorna errores del dominio

## Ejemplo

```go
corriente, err := service.CalcularCorrienteNominal(5000, 220, 0.9)
```
