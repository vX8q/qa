# QA 

Описание

Сервис позволяет:
- создавать вопросы (`POST /questions/`)
- просматривать вопросы (`GET /questions/`)
- получать вопрос с ответами (`GET /questions/{id}`)
- создавать ответы на вопросы (`POST /questions/{id}/answers`)
- удалять вопросы и ответы (`DELETE /questions/{id}`, `DELETE /answers/{id}`)

Миграции базы данных выполняются через GORM или Goose (по необходимости).

---

Технологии

- Go 1.24
- PostgreSQL 16 (Alpine)
- GORM ORM
- Docker + Docker Compose
- SQLite (тестирование)
- `httptest` и `testify` для unit-тестов

---

Требования

- Docker
- Docker Compose
- Go 

---

Установка и запуск

1. Клонировать репозиторий:

``bash
1. git clone https://github.com/vX8q/qa
2. поднять сервис через Docker Compose docker-compose up -d --build
3. проверка
docker-compose ps
docker-compose logs app --tail=50
docker-compose logs db --tail=50

---

Вопросы

Создать вопрос:

curl -X POST http://localhost:8080/questions/ \
-H "Content-Type: application/json" \
-d '{"text":"Тест вопрос"}'

Список вопросов:

curl http://localhost:8080/questions/

Получить вопрос с ответами:

curl http://localhost:8080/questions/1

---

Ответы

Создать ответ:

curl -X POST http://localhost:8080/questions/1/answers \
-H "Content-Type: application/json" \
-d '{"user_id":"550e8400-e29b-41d4-a716-446655440000","text":"Ответ"}'


Получить ответ:

curl http://localhost:8080/answers/1

---

Тесты
docker-compose exec app go test ./...
