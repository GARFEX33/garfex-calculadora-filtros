# Garfex — Calculadora de Memorias de Cálculo Eléctrico

API Go + Frontend SvelteKit para generar memorias de cálculo de instalaciones eléctricas según normativa NOM (México).

---

## Stack

| Capa | Tecnología |
|------|-----------|
| Backend | Go 1.24, Gin, pgx/v5 |
| Frontend | SvelteKit 2, Svelte 5, TypeScript, Tailwind 4 |
| Base de datos | PostgreSQL (Supabase) |
| PDF | Gotenberg 8 (Chromium) |
| Arquitectura | Hexagonal + Vertical Slicing |

---

## Deploy en servidor (Docker)

### Requisitos previos

El servidor debe tener corriendo:
- **PostgreSQL** en puerto `5434` (Supabase local)
- **Gotenberg** en puerto `3015`

### 1. Clonar el repo

```bash
cd ~
git clone https://github.com/GARFEX33/garfex-calculadora-filtros.git
cd garfex-calculadora-filtros
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
nano .env
```

Completar los valores reales:

```env
DB_HOST=host.docker.internal
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=tu_password_de_supabase
DB_NAME=postgres

GOTENBERG_URL=http://host.docker.internal:3015/forms/chromium/convert/html

# IP o dominio del servidor — la ve el navegador
PUBLIC_API_URL=http://192.168.1.X:8080
```

### 3. Levantar

```bash
docker compose up -d --build
```

### 4. Verificar

```bash
docker compose ps
# calculadora-filtros-api   → http://servidor:8080
# calculadora-filtros-web   → http://servidor:5173

curl http://localhost:8080/health
# {"status":"ok"}
```

### Comandos útiles

```bash
# Ver logs del API
docker compose logs -f calculadora-filtros-api

# Ver logs del frontend
docker compose logs -f calculadora-filtros-web

# Actualizar a la última versión
git pull && docker compose up -d --build

# Detener
docker compose down
```

---

## Desarrollo local

### Backend

```bash
# Copiar y completar variables
cp .env.example .env

# Correr
go run cmd/api/main.go

# Tests
go test ./...

# Build
go build ./cmd/api/...
```

### Frontend

```bash
cd frontend/web
npm install
npm run dev
# → http://localhost:5173
```

### Preview de PDF (hot-reload)

```bash
go run cmd/pdf_preview/main.go
# → http://localhost:3000?empresa=garfex
```

---

## API

Swagger UI disponible en `http://localhost:8080/swagger/index.html`

### Endpoint principal

```bash
POST /api/v1/calculos/memoria
```

```json
{
  "modo": "MANUAL_AMPERAJE",
  "amperaje_nominal": 50,
  "tension": 220,
  "tipo_canalizacion": "TUBERIA_PVC",
  "hilos_por_fase": 1,
  "longitud_circuito": 10,
  "itm": 100,
  "factor_potencia": 0.9,
  "estado": "Sonora",
  "sistema_electrico": "DELTA",
  "material": "Cu"
}
```

---

## Proyecto privado — GARFEX
