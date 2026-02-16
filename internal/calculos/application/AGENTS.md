# Calculos — Application Layer

Orquesta domain services. Define contratos (ports), no implementaciones.

## Trabajar en esta Capa

Esta capa es responsabilidad del **`application-agent`**. El agente ejecuta su ciclo completo:

```
brainstorming-application → writing-plans-application → executing-plans-application
```

**NO modificar directamente** — usar el sistema de orquestación.

## Estructura

```
internal/calculos/application/
├── port/           # Interfaces que infrastructure implementa
│   ├── driver/     # Ports de entrada (para ser llamados)
│   └── driven/     # Ports de salida (dependencias)
├── usecase/        # Orquestadores
│   └── helpers/    # Funciones auxiliares
└── dto/            # Entrada/salida
```

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/calculos/domain/entity`
- `internal/calculos/domain/service`
- stdlib de Go

## Dependencias prohibidas

- `internal/calculos/infrastructure/` — **nunca**
- Frameworks (Gin, pgx, etc.)

## Cómo modificar esta capa

### Para nueva feature

```bash
# Primero: domain-agent completa el dominio
# Luego:
orchestrate-agents --agent application --feature nueva-feature
```

### Para cambios en calculos existente

```bash
# Orquestador:
# "application-agent: agregar use case para exportar memoria a PDF"
```

## Flujo de Use Cases (orden obligatorio)

1. Corriente Nominal
2. Ajuste de Corriente
3. Seleccionar TipoCanalizacion
4. Resolver tabla ampacidad + temperatura
5. Conductor de Tierra
6. Dimensionar Canalización
7. Calcular Caída de Tensión

## Reglas de Application

### Use cases solo orquestan

```go
// BIEN: solo coordina
data, err := uc.repo.Find(...)
result, err := service.Calculate(...)
return dto.FromDomain(result), nil

// MAL: lógica de negocio aquí
if valor > 100 {  // esto va en domain
    return error
}
```

### Ports

- **Driver**: interfaces que expone application (para HTTP, gRPC, CLI)
- **Driven**: interfaces que application necesita (Repository, Clientes externos)

### DTOs

- Structs planos sin métodos
- Nunca exponer entidades de domain
- Mapping explícito `domain ↔ DTO`

## Referencias

- Agente: `application-agent`
- Skill: `.agents/skills/orchestrating-agents/SKILL.md`

## QA Checklist

- [ ] `go test ./internal/calculos/application/...` pasa
- [ ] Use case < ~80 líneas
- [ ] Sin lógica de negocio (fórmulas, validaciones complejas)
- [ ] Sin imports de infrastructure
- [ ] Error wrapping con `%w`
