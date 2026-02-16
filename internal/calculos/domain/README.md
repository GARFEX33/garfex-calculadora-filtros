# Domain Layer

Capa de dominio: entidades, value objects y lógica de negocio pura.

## Responsabilidades

- Definir entidades del negocio (TipoCanalizacion, SistemaElectrico, ITM, etc.)
- Implementar servicios de cálculo según normativa NOM y IEEE-141
- **NO depende** de application, infrastructure ni frameworks externos

## Estructura

```
domain/
├── entity/     # Entidades del dominio
└── service/    # Servicios de cálculo puro
```

## Reglas

- Sin imports externos (solo stdlib)
- Sin JSON tags, context, logging
- Errores de negocio definidos en el dominio

## Servicios Disponibles

- CalculoCorrienteNominal
- AjusteCorriente
- CalculoConductor
- CalculoCanalizacion
- CalculoTierra
- CalculoCaidaTension (IEEE-141)
- SeleccionarTemperatura
- CalcularFactorTemperatura
- CalcularFactorAgrupamiento
- CalcularCharolaEspaciado
- CalcularCharolaTriangular
- CalcularFactorUso
