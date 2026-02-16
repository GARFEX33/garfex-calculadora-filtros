# {Feature} — {Layer} Layer

{Descripción en una línea del propósito de esta capa.}

## Trabajar en esta Capa

Esta capa es responsabilidad del **`{layer}-agent`**. El agente ejecuta su ciclo completo:

```
brainstorming-{layer} → writing-plans-{layer} → executing-plans-{layer}
```

**NO modificar directamente** — usar el sistema de orquestación.

## Estructura

```
internal/{feature}/{layer}/
├── {subdirectorio}/    # {qué contiene}
├── {subdirectorio}/    # {qué contiene}
└── {archivo}.go        # {qué contiene}
```

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/{feature}/{capa-inferior}` (si aplica)
- stdlib de Go

## Dependencias prohibidas

- `internal/{feature}/{capa-superior}/` — **nunca**
- Otras features — **nunca**
- Frameworks externos (si es domain)

## Cómo modificar esta capa

### Para nueva feature

```bash
# Orquestador despacha el agente correspondiente:
orchestrate-agents --agent {layer} --feature {feature}
```

### Para cambios menores

```bash
# Orquestador da instrucciones específicas:
# "{layer}-agent: agregar método X a entidad Y"
```

## Reglas de Oro

1. **{Regla 1}** — {explicación breve}
2. **{Regla 2}** — {explicación breve}
3. **{Regla 3}** — {explicación breve}
4. **{Regla 4}** — {explicación breve}
5. **{Regla 5}** — {explicación breve}

## Referencias

- Agente: `.opencode/agents/{layer}-agent.md`
- Skills: `brainstorming-{layer}`, `writing-plans-{layer}`, `executing-plans-{layer}`

## QA Checklist

- [ ] `go test ./internal/{feature}/{layer}/...` pasa
- [ ] Sin imports de capas prohibidas
- [ ] {Check específico de la capa}
- [ ] {Check específico de la capa}
