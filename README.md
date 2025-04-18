# Telegram User Channel Membership Service

Сервис для проверки членства пользователей в Telegram-каналах с использованием микросервисной архитектуры.

## Архитектура

Проект состоит из двух микросервисов:
- HTTP-сервис для обработки внешних запросов
- gRPC-сервис для взаимодействия с Telegram API и базой данных

## Требования

- Docker и Docker Compose
- Go 1.2 или выше (для локальной разработки)
- PostgreSQL (запускается в Docker)

## Установка и запуск

1. Клонируйте репозиторий

2. Сгенерируйте gRPC код из proto файлов:
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/telegram.proto
```

3. Удалите из интерфейса **TelegramServiceServer**, в файле telegram_grpc.pb.go, метод **mustEmbedUnimplementedTelegramServiceServer()** 

4. Запустите сервисы через Docker Compose:
```bash
docker-compose up --build
```

## API Endpoints

### Проверка членства пользователя в канале

```
POST http://localhost:8080/check-user
```

Тело запроса (JSON):
```json
{
    "bot_token": "your_bot_token",
    "channel_link": "@channel_name",
    "user_id": 123456789
}
```

Ответ:
```json
{
    "is_member": true,
    "error": ""
}
```

## База данных

Таблица `users` создается автоматически при первом запуске:

```sql
create table users(
    id serial primary key,
    telegram_id bigint not null unique,
    channel_id bigint not null,
    created_at timestamp default now()
);
```

## Разработка

Для локальной разработки:

1. Установите зависимости:
```bash
go mod download
```

2. Запустите PostgreSQL:
```bash
docker-compose up postgres
```

3. Запустите сервисы по отдельности:
```bash
# Terminal 1 - gRPC Service
go run ./grpc-service

# Terminal 2 - HTTP Service
go run ./http-service
```#   t e s t P r o j e c t  
 