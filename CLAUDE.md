# CLAUDE.md

Ledger service for robust and high frequency operations.

Сервис Ledger является источником событий для высокочастотных операций. 
Он обеспечивает надежную обработку транзакций и поддерживает интеграцию с различными системами через gRPC и HTTP API.


## Сборка и тесты

- Сборка: make build
- Тесты: make test
- e2e-тесты: go test -tags e2e_functional,e2e_smoke -v ./e2e/...
- smoke e2e-тесты: go test -tags e2e_smoke -v ./e2e/...
- functional e2e-тесты: go test -tags e2e_functional -v ./e2e/...
- Линтер: make lint

## e2e-тесты

### smoke-тесты
- suites находятся в e2e/suites/smoke
- функции запуска находятся в e2e/tests/smoke_test.go
- для smoke-тестов используй тэг сборки e2e_smoke

### functional-тесты
- suites находятся в e2e/suites/functional
- функции запуска находятся в e2e/tests/functional_test.go
- Для functional-тестов используй тэг сборки e2e_functional

## Spec-Driven Development Workflow

1. **Расположение спецификаций** — `specs/<feature-slug>/`:
   - `spec.md` — требования, критерии приёмки, edge cases
   - `design.md` — техническое решение: схема событий, API-контракт,
     схема БД, инварианты event store'а
   - `tasks.md` — разбивка на задачи реализации
2. **Порядок работы:**
   - Задача без ссылки на спеку → сначала найди спеку в `specs/` или
     предложи её создать. Не начинай писать код без неё.
   - Спека есть, но неполная (нет критериев приёмки/edge cases) →
     дополни её и покажи диффом перед реализацией.
   - Реализация ссылается на конкретный пункт спеки в commit message:
     `Implements: specs/<feature-slug>/spec.md#section`.
3. **Статус спеки** в frontmatter: `draft` / `approved` / `in-progress` /
   `implemented` / `deprecated`. Не начинай реализацию по `draft`-спекам —
   уточни у пользователя, что она согласована.
4. **Расхождение спеки и кода** — спека источник истины. Сначала правишь
   спеку, потом e2e тесты, потом код, а не наоборот.

## Структура репозитория

```
/cmd/<service>       — точки входа (main.go)
/internal/<domain>   — бизнес-логика, приватные пакеты
/pkg/                — переиспользуемые публичные пакеты
/migrations/         — postgres-миграции (см. .claude/rules/postgres-migrations.md)
/e2e/                — e2e-тесты (см. .claude/rules/e2e-tests.md)
/specs/              — спецификации (spec-driven development)
```


## Что важно помнить

- выполняй задачи в отдельных ветках `feature/<feature-name>`, не в `main`/`master`
- Не начинай реализацию, если для неё нет approved-спеки в `specs/`.
- Публичные контракты (grpc/http, схема БД, формат события) описываются
  в `design.md` до генерации кода, а не после.
- E2E-тесты должны быть написаны до реализации функциональности, чтобы гарантировать соответствие кода спецификациям. 
- Не удаляй код, который был написан после взятия задачи в работу, без согласования с автором.


