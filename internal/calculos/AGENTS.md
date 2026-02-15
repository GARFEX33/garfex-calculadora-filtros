# Feature: calculos

Memoria de cálculo eléctrico según normativa NOM (México).

Esta feature implementa el cálculo completo de una instalación eléctrica:
corriente nominal → ajuste por temperatura/agrupamiento → selección de conductor →
conductor de tierra → dimensionamiento de canalización → caída de tensión (IEEE-141).

## Estructura interna

```
calculos/
  domain/          ← entidades y servicios de cálculo puro
  application/     ← ports, use cases, DTOs
  infrastructure/  ← adapters HTTP (driver) y CSV (driven)
```

## Reglas de aislamiento

- Esta feature NO importa `equipos/` ni ninguna otra feature
- Solo importa `shared/kernel/` para value objects compartidos
- `cmd/api/main.go` es el único que instancia y conecta esta feature

## Agentes responsables

| Capa | Agente | AGENTS.md |
| ---- | ------ | --------- |
| domain/ | `domain-agent` | `internal/calculos/domain/AGENTS.md` |
| application/ | `application-agent` | `internal/calculos/application/AGENTS.md` |
| infrastructure/ | `infrastructure-agent` | `internal/calculos/infrastructure/AGENTS.md` |
