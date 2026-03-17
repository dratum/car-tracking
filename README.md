# Car-Tracking

GPS-трекинг сервис для автомобиля. Принимает координаты с ESP32 + NEO-6M, хранит треки поездок, отображает маршруты на карте и считает статистику пробега.

## Tech Stack

### Backend
- **Go 1.25** + Chi router
- **TimescaleDB** (pgx/v5) — GPS-точки (hypertable)
- **MongoDB** — поездки, пользователи

### Frontend
- **SvelteKit 2** (Svelte 5), TypeScript
- **Tailwind CSS v4**
- **Leaflet** — карты (OpenStreetMap)
- **pnpm** — пакетный менеджер

### Infrastructure
- **Docker Compose** — TimescaleDB + MongoDB + App
- **GitHub Actions** — CI (build, test, lint)

## Quick Start

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- pnpm
- Node.js 22+

### Local Development

```bash
# 1. Clone & configure
git clone https://github.com/dratum/car-tracking.git
cd car-tracking
cp .env.example .env

# 2. Start databases
make up

# 3. Run backend (terminal 1)
make run

# 4. Run frontend (terminal 2)
make web-dev
```

Backend: http://localhost:8080
Frontend: http://localhost:5173

### Docker (production)

```bash
make docker-build
make up
```

All-in-one на http://localhost:8080 (Go раздаёт API + SPA).

## Make Commands

| Command | Description |
|---------|-------------|
| `make up` | Start containers (DB + App) |
| `make down` | Stop containers |
| `make logs` | Container logs (follow) |
| `make run` | Run backend locally |
| `make build` | Build Go binary |
| `make docker-build` | Build Docker image |
| `make web-install` | Install frontend deps |
| `make web-dev` | Frontend dev server |
| `make web-build` | Build frontend for production |

## API

### Device API (`/api/device/`)
Auth: `X-API-Key` header

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/device/trip/start` | Start trip |
| POST | `/api/device/location` | Send GPS point |
| POST | `/api/device/trip/end` | End trip |

### Web API (`/api/v1/`)
Auth: `Authorization: Bearer <JWT>`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login |
| GET | `/api/v1/trips` | List trips |
| GET | `/api/v1/trips/:id` | Trip details |
| GET | `/api/v1/trips/:id/points` | GPS points |
| GET | `/api/v1/stats?period=week` | Mileage stats |

## Project Structure

```
cmd/server/              — entrypoint
internal/
  api/handler/           — HTTP handlers
  api/middleware/         — JWT & API-key auth
  config/                — env config
  domain/model/          — data models
  domain/geo/            — haversine distance
  repository/timescale/  — GPS points (TimescaleDB)
  repository/mongo/      — trips, users (MongoDB)
  service/               — business logic
web/                     — SvelteKit SPA
deployments/             — docker-compose
```

## License

MIT
