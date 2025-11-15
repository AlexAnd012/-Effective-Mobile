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
