# QA Checklist — Application Layer

Verificaciones obligatorias antes de commit.

## Tests

- [ ] `go test ./internal/calculos/application/...` pasa

## Estructura

- [ ] Use case < ~80 líneas
- [ ] Un use case = una responsabilidad (no combinar funcionalidades)

## Dependencias

- [ ] Sin lógica de negocio (fórmulas, validaciones complejas)
- [ ] Sin imports de infrastructure

## DTOs

- [ ] DTOs usan SOLO primitivos (string, int, float64)
- [ ] No exponen value objects ni entities de domain

## Conversiones

- [ ] Use case convierte DTO → value objects antes de llamar a domain
- [ ] Use case convierte resultado domain → DTO antes de retornar

## Errores

- [ ] Error wrapping con `%w`
