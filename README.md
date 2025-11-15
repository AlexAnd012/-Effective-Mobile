# Тестовое задание Junior Golang Developer Effective Mobile
## Задача:
спроектировать и реализовать REST-сервис для агрегации данных об онлайн-подписках пользователей.
## Требования:
1 Выставить HTTP-ручки для CRUDL-операций над записями о подписках. Каждая запись
содержит:  
> 1 Название сервиса, предоставляющего подписку  
> 2 Стоимость месячной подписки в рублях  
> 3 ID пользователя в формате UUID  
> 4 Дата начала подписки (месяц и год)  
> 5 Опционально дата окончания подписки

2 Выставить HTTP-ручку для подсчета суммарной стоимости всех подписок за выбранный  
период с фильтрацией по id пользователя и названию подписки   
3 СУБД – PostgreSQL. Должны быть миграции для инициализации базы данных  
4 Покрыть код логами  
5 Вынести конфигурационные данные в .env/.yaml-файл  
6 Предоставить swagger-документацию к реализованному API  
7 Запуск сервиса с помощью docker compose  
## Примечания:  
1 Проверка существования пользователя не требуется. Управление пользователями вне  
зоны ответственности вашего сервиса  
2 Стоимость любой подписки – целое число рублей, копейки не учитываются  
Пример тела запроса на создание записи о подписке:  
json  
{  
“service_name”: “Yandex Plus”,  
“price”: 400,  
“user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”,  
“start_date”: “07-2025”  
}  
  

# REST-сервис для хранения и расчёта стоимостей пользовательских подписок.  
Поддерживает CRUDL, расчёт суммарной стоимости за период, миграции PostgreSQL, логирование, middleware, конфиги из .env, Swagger и запуск через Docker Compose.  
# Стек
Go (стандартная библиотека: net/http, context)
chi роутинг (github.com/go-chi/chi/v5)
pgx/pgxpool PostgreSQL драйвер/пул соединений (github.com/jackc/pgx/v5/pgxpool)
slog структурные логи (log/slog)
swaggo/http-swagger Swagger UI 
docker compose запуск postgres + приложение
migrate/migrate применение миграций

# Реализовано 
## CRUDL для подписок:  
POST /api/v1/subscriptions   
GET /api/v1/subscriptions/{id}  
GET /api/v1/subscriptions (фильтры + пагинация)  
PUT /api/v1/subscriptions/{id} (полная замена)  
DELETE /api/v1/subscriptions/{id}  
## Расчёт суммы за период:  
GET /api/v1/cost/total?from=MM-YYYY&to=MM-YYYY[&user_id=&service_name=]
## Здоровье:
GET /healthz жив ли процесс  
GET /readyz готов ли сервис (ping БД с таймаутом)  
## Логи  
(access + recovery), request-id, конфиги из .env  
## Миграции PostgreSQL (migrations/0001_init.up.sql)  
## Swagger UI (/swagger/index.html)
# Архитектура и расположение
.
├── cmd/  
│   └── main.go                     # точка входа 
├── internal/   
│   ├── app/  
│   │   └── server.go               # обёртка над http.Server (start/shutdown)  
│   ├── config/  
│   │   └── config.go               # чтение .env, валидация, ошибки на пустые  
│   ├── domain/  
│   │   ├── errors.go               # ошибки валидации
│   │   └── subscription.go         # доменная модель + валидация дат/цен  
│   ├── dto/  
│   │   ├── subscription_dto.go     # Create/Update/List/Response  
│   │   └── cost_dto.go             # TotalCostQuery/Response  
│   ├── http_server/  
│   │   ├── httx/   
│   │   │   ├── handlers/  
│   │   │   │   ├── handlers_health.go  # /healthz, /readyz   
│   │   │   │   ├── handlers_subscription.go # CRUDL  
│   │   │   │   └── handlers_cost.go    # /cost/total  
│   │   │   └── responses.go          # JSON/Error helpers  
│   │   ├── middleware/  
│   │   │   ├── accesslog.go        # access-log  
│   │   │   ├── recovery.go         # panic → 500 + лог стека  
│   │   │   └── requestid.go        # request-id  
│   │   └── router/  
│   │       └── router.go           # конструктор chi-маршрутизатора  
│   ├── logging/  
│   │   └── logger.go               # фабрика slog.Logger (уровни)  
│   ├── repo/  
│   │   ├── postgres/  
│   │   │   └── postgres.go         # init pgxpool + Ping с таймаутом  
│   │   └── subscription_repo.go    # интерфейс и реализация на PostgreSQL (CRUD+CalcTotal)  
│   └── service/  
│       └── subscription.go         # бизнес-логика, валидации, маппинг DTO  
├── migrations/  
│   └── 0001_init.up.sql            # схема таблицы subscriptions + индексы  
├── docs/                           # сгенерированные swag-файлы (когда подключено)  
├── .env                            # конфигурация приложения  
├── .env.example                    # пример конфигурации приложения  
├── Dockerfile  
├── docker-compose.yml  
├── go.mod  
└── Makefile  
## Слои  
domain «чистая» модель и инварианты.  
dto вход/выход API (что читает/пишет HTTP).  
service бизнес-логика: валидации, парсинг дат MM-YYYY, превращение DTO в Domain, вызовы репозитория.  
repo SQL к PostgreSQL. В тот числе CalcTotal без помесячной генерации - через формулы пересечения периодов.  
http_server веб-уровень: хендлеры, middleware, маршрутизация, JSON-ответы.
