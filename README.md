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

### Из флагов

```bash
plea-cli new --name <project-name> --module <module-path>
```

#### Флаги

| Флаг          | Псевдоним | Описание                                               | По умолчанию | Обязательный |
| ------------- | --------- | ------------------------------------------------------ | ------------ | ------------ |
| `--name`      | `-n`      | Название проекта                                       | —            | да           |
| `--module`    | `-m`      | Имя модуля в `go.mod`                                  | —            | да           |
| `--output`    | `-o`      | Директория для генерации                               | `.`          | нет          |
| `--log-level` | `-ll`     | Уровень логирования (`debug`, `info`, `warn`, `error`) | `debug`      | нет          |
| `--verbose`   | `-v`      | Подробный вывод (директории, файлы, команды)           | `false`      | нет          |

```bash
plea-cli new \
  --name my-service \
  --module github.com/username/my-service \
  --output ./projects \
  --log-level info \
  --verbose
```

### Из манифеста

Создайте `plea.yaml` (или `plea.yml` / `plea.json`) в корне проекта:

```yaml
name: my-app
module: github.com/username/my-app
output: ./projects
log_level: info
verbose: true
```

```json
{
  "name": "my-app",
  "module": "github.com/username/my-app",
  "output": "./projects",
  "log_level": "info",
  "verbose": true
}
```

Запустите:

```bash
plea-cli manifest                      # ищет ./plea.yaml по умолчанию
plea-cli manifest --path ./plea.json   # явный путь к файлу
```

#### Флаги

| Флаг     | Псевдоним | Описание               | По умолчанию  |
| -------- | --------- | ---------------------- | ------------- |
| `--path` | `-p`      | Путь к файлу манифеста | `./plea.yaml` |

## Структура генерируемого проекта

```
<output>/<name>/
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
├── go.mod
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

## Требования

- Go 1.24+
