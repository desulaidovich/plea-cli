# plea-cli

CLI-инструмент для генерации Go-проектов из готового шаблона.

## Установка

```bash
go install github.com/desulaidovich/plea-cli@latest
```

## Сборка

```bash
git clone https://github.com/desulaidovich/plea-cli
cd plea-cli
go build -o plea-cli .
```

## Использование

```bash
plea-cli new --name <project-name> --repo <github-username-or-org>

# пример
plea-cli new --name my-service --repo desulaidovich
```

### Флаги

| Флаг        | Псевдоним | Описание                                     | Обязательный |
| ----------- | --------- | -------------------------------------------- | ------------ |
| `--name`    | `-n`      | Название проекта                             | да           |
| `--repo`    | `-r`      | GitHub username или организация              | нет          |
| `--output`  | `-o`      | Директория для генерации                     | нет          |
| `--verbose` | `-v`      | Подробный вывод (созданные файлы)            | нет          |

После генерации инструмент выведет следующие шаги:

```
next steps:
  cd <project-name>
  go mod init github.com/<repo>/<project-name>
  go mod tidy
```

## Структура генерируемого проекта

```
<name>/
├── cmd/
│   └── app/
│       └── main.go          # Точка входа
├── internal/
│   └── app/
│       └── app.go           # Логика приложения
├── pkg/
│   ├── log/
│   │   └── log.go           # Структурированное логирование (slog)
│   └── runner/
│       └── runner.go        # Graceful shutdown
└── Makefile
```

### Что входит в шаблон

**Functional options** — конфигурация `App` и `Runner` через `WithXxx`-опции с валидацией на старте.

**Logger** (`pkg/log`) — обёртка над `log/slog`:

- Форматы: `json` (по умолчанию) и `text`
- Уровни: `debug`, `info`, `warn`, `error`, `panic`
- Настраиваемый формат времени
- Метод `.With(map[string]any)` для добавления полей; структуры и карты раскладываются в `slog.Group`

**Runner** (`pkg/runner`) — запуск приложения с graceful shutdown:

- Перехватывает `SIGINT` / `SIGTERM`
- `Start` и `Stop` выполняются в `errgroup` — ошибка любого из них завершает оба
- Паника в `Start` или `Stop` перехватывается и возвращается как ошибка
- Таймаут остановки задаётся через `WithStopTimeout`
- Защита от повторного вызова `Run` через `atomic.Bool`

**Build-time variables** — `version` и `build` инжектируются через `-ldflags` при сборке (`make build`).

**Makefile**:

```makefile
make build   # сборка бинаря в bin/<project-name> с ldflags
make run     # go run ./cmd/app (цель по умолчанию)
```
