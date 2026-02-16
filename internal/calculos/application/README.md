# Application Layer

Capa de aplicación: puertos, casos de uso y DTOs.

## Responsabilidades

- Definir puertos (interfaces) para infraestructura
- Orquestar casos de uso
- Transformar DTOs entre dominio e infraestructura
- **NO implementa** reglas de negocio ni adaptadores

## Estructura

```
application/
├── port/       # Interfaces para infraestructura
├── usecase/   # Orquestadores de lógica
└── dto/       # Data Transfer Objects
```

## Dependencias

- ✅ `domain/`
- ✅ `shared/kernel/valueobject/`
- ❌ `infrastructure/`

## Casos de Uso

- `OrquestadorMemoriaCalculo` - Orquestador principal
- `CalcularMemoria` - Memoria de cálculo completa
- `CalcularCorriente` - Cálculo de corriente
- `AjustarCorriente` - Ajuste por factores
- `DimensionarCanalizacion` - Dimensionamiento
- `SeleccionarConductor` - Selección de conductor
