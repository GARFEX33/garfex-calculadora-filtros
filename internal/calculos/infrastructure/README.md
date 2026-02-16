# Infrastructure Layer

Capa de infraestructura: adaptadores, repositories y handlers HTTP.

## Responsabilidades

- Implementar puertos definidos en application
- Conectar con sistemas externos (CSV, DB, HTTP)
- Manejar serialización HTTP
- **Define** adaptadores, **no** reglas de negocio

## Estructura

```
infrastructure/
├── adapter/
│   ├── driver/http/    # HTTP Handlers
│   └── driven/        # Repositories (CSV, DB)
└── router.go          # Configuración de rutas
```

## Adaptadores

| Tipo | Implementación |
|------|----------------|
| Driver | HTTP Handler (Gin) |
| Driven | CSV Repository |

## Dependencias

- ✅ `application/port/` (implementa interfaces)
- ✅ Frameworks (Gin, pgx)
- ❌ `domain/` (solo usa, no conoce)
