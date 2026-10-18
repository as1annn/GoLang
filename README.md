# SIS — Task Processor (Checkpoint 1: Core Go)

```
go run .
go test ./...
```

Структура:
- `internal/task` — доменная модель, статусы, правила переходов, ошибки
- `internal/storage` — интерфейс `Repository` + `MemoryStore`
- `internal/service` — бизнес-логика
- `main.go` — демонстрация
