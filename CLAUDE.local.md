# Auto-Tracking

GPS-трекинг сервис на Go. Принимает координаты с устройств, хранит треки поездок.

## Quick start/commands

```bash
cp .env.example .env   # настроить переменные окружения
make up                # TimescaleDB + MongoDB в Docker
make run               # запуск сервера локально (порт 8080)
```

## Architecture

Слоистая архитектура: handler → service → repository.

```
cmd/server/main.go          — точка входа, DI, graceful shutdown
internal/
  config/config.go          — конфигурация через env-переменные (caarlos0/env)
  domain/model/             — структуры: GPSPoint, Trip, User, Vehicle
  domain/geo/               — геовычисления (haversine)
  repository/timescale/     — GPS-точки (TimescaleDB, pgxpool)
  repository/mongo/         — поездки, пользователи, транспорт (MongoDB)
  service/                  — бизнес-логика (tracking, trip, stats)
  api/router.go             — маршруты (chi)
  api/handler/              — HTTP-хендлеры (device, auth, trip, stats)
  api/middleware/            — JWT и API-key авторизация
deployments/docker-compose.yml — TimescaleDB + MongoDB
```

## Tech

- **Go 1.26**, chi router
- **TimescaleDB** (pgx/v5, pgxpool) — GPS-точки с гипертаблицей по time
- **MongoDB** (mongo-driver/v2) — поездки, пользователи
- Конфигурация: `.env` файл, парсинг через `caarlos0/env`

## Конвенции

- SQL-запросы: именованные параметры `@param` + `pgx.NamedArgs`, "сырые" запросы к БД в формате 
`insert into table (
								  time
							  , trip_id
							  )
						values (
								  @time
							  , @trip_id
							  )`,
  `select
	       time
	     , trip_id
	from table`           
- Репозитории принимают `*pgxpool.Pool` / `*mongo.Database` напрямую (без интерфейсов)
- Сервисы принимают конкретные типы интерфейсов реализующие необходимые методы
- Ошибки оборачиваются через `fmt.Errorf("контекст: %w", err)`
- Конфиг — env-переменные с дефолтами, `.env` файл опционален
- Не писать однострочные `if err := something(); err != nil {}`

## API

Два набора эндпоинтов:
- **Device API** (`/api/device/`) — авторизация по `X-API-Key`
- **Web API** (`/api/v1/`) — авторизация по JWT Bearer token

## Make-команды

`make help` — список всех команд.

### Important notes
 see @docs/TECHNICAL_REQUIREMENTS.md for develop project and questions
