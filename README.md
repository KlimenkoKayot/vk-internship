# PubSub Service (vk-internship)
Тестовое задание на Golang Intern VK

## Содержание

1. [Описание проекта](#описание-проекта)
2. [Архитектура и использованные техники](#архитектура-и-использованные-техники)
3. [Структура проекта](#структура-проекта)
4. [Запуск сервиса](#запуск-сервиса)
   - [Требования](#требования)
   - [Сборка и запуск](#сборка-и-запуск)
5. [Конфигурация](#конфигурация)
6. [API](#api)
7. [Тестирование](#тестирование)
8. [Логирование](#логирование)
9. [Лицензия](#лицензия)

## Описание проекта

Проект реализует сервис подписок на события по принципу Publisher-Subscriber с использованием gRPC. Состоит из двух основных компонентов:

1. Пакет `subpub` - реализация шины событий (Publisher-Subscriber)
2. Сервис `service` - gRPC сервис для подписки и публикации событий

## Архитектура и использованные техники

### Основные техники и паттерны:
- **Graceful Shutdown** - корректная обработка завершения работы сервиса
- **Dependency Injection (DI)** - внедрение зависимостей через интерфейсы
- **Интерфейсы** - абстракция реализации для легкого тестирования и замены компонентов
- **Publisher-Subscriber** - паттерн для обработки событий
- **Слоистая архитектура** (domain, infrastructure, interfaces)

### Логирование
Используется модуль `libs/logger` с поддержкой нескольких адаптеров (Logrus, Zap) и возможностью расширения.

## Структура проекта

```
.
├── libs/logger         # Библиотека логирования с поддержкой нескольких адаптеров
├── service/            # Основной сервис
│   ├── cmd/            # Точка входа
│   ├── config/         # Конфигурация
│   ├── internal/       # Внутренняя логика сервиса
│   └── pkg/grpc/       # gRPC спецификация и сервер
└── subpub/             # Реализация шины событий
```

## Запуск сервиса

### Требования
- Go 1.20+
- Установленный protoc и Go плагины для генерации gRPC кода

### Сборка и запуск

1. Установите зависимости:
```bash
cd service && go mod download
cd ../subpub && go mod download
```

2. Сгенерируйте gRPC код (из директории service):
```bash
cd ../service && protoc --go_out=. --go-grpc_out=. ./pkg/grpc/pb/pubsub.proto
```

3. Запустите сервис:
```bash
go run ./cmd/.
```

## Конфигурация

Конфигурация сервиса находится в `service/config/config.yaml`. Доступные параметры:
- Порт gRPC сервера

## API

Сервис предоставляет gRPC API с двумя методами:

### Подписка на события
```protobuf
rpc Subscribe(SubscribeRequest) returns (stream Event);
```

### Публикация события
```protobuf
rpc Publish(PublishRequest) returns (google.protobuf.Empty);
```

## Тестирование

Для тестирования используется стандартный пакет testing Go. Запуск тестов:
```bash
cd subpub && go test ./...
cd ../service && go test ./...
```

## Логирование

Сервис поддерживает несколько адаптеров логирования (Logrus, Zap). Адаптер выбирается в конфигурации.

## Лицензия

Проект распространяется под лицензией MIT. Подробности в файле LICENSE.