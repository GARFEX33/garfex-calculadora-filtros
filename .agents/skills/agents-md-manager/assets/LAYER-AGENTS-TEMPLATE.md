# {Feature} — {Layer} Layer

{Descripción en una línea del propósito de esta capa.}

## Trabajar en esta Capa

Esta capa sigue las reglas de arquitectura hexagonal y vertical slicing.

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

直接amente, siguiendo las reglas de esta capa y las reglas globales en `docs/reference/structure.md`.

## Reglas de Oro

1. **{Regla 1}** — {explicación breve}
2. **{Regla 2}** — {explicación breve}
3. **{Regla 3}** — {explicación breve}
4. **{Regla 4}** — {explicación breve}
5. **{Regla 5}** — {explicación breve}

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)

## QA Checklist

- [ ] `go test ./internal/{feature}/{layer}/...` pasa
- [ ] Sin imports de capas prohibidas
- [ ] {Check específico de la capa}
- [ ] {Check específico de la capa}
