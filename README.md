# Notes API — REST API для управления заметками

Учебный проект на Go для создания REST API с базой данных PostgreSQL.

## 🎯 Цель проекта

Реализовать REST API-сервер, который позволяет:
- Регистрировать пользователей
- Создавать текстовые заметки
- Просматривать и редактировать только свои заметки
- Работать с пагинацией и сортировкой

---

## 🛠️ Технологии

### Backend:
- **Go 1.21+** — язык программирования
- **PostgreSQL 15** — база данных
- **Docker** — контейнеризация БД

### Библиотеки:
- `github.com/go-chi/chi/v5` — HTTP роутер
- `github.com/lib/pq` — PostgreSQL драйвер
- `github.com/joho/godotenv` — загрузка .env файлов
- `github.com/pressly/goose/v3` — миграции БД

### Инструменты:
- **goose** — управление миграциями
- **curl/Postman** — тестирование API

---

## 📐 Архитектура проекта
```
notes-api/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── models/                  # Структуры данных
│   │   ├── user.go              # Модель пользователя
│   │   ├── note.go              # Модель заметки
│   │   └── errors.go            # Кастомные ошибки
│   ├── storage/                 # Работа с БД
│   │   ├── storage.go           # Инициализация storage
│   │   ├── user_storage.go      # CRUD для users
│   │   └── note_storage.go      # CRUD для notes
│   └── handlers/                # HTTP обработчики
│       ├── user_handler.go      # Handler для /users
│       ├── note_handler.go      # Handler для /notes
│       └── response.go          # Вспомогательные функции ответов
├── migrations/                  # SQL миграции
│   ├── 001_create_users.sql
│   └── 002_create_notes.sql
├── docker-compose.yml           # Настройка PostgreSQL
├── .env                         # Переменные окружения (не в Git!)
├── .gitignore
├── go.mod
└── README.md
```

---

## 🔧 Установка и запуск

### Предварительные требования:
- Go 1.21+
- Docker Desktop
- Git

### 1. Клонирование проекта:
```bash
git clone https://github.com/Balyshev/notes-api.git
cd notes-api
```

### 2. Установка зависимостей:
```bash
go mod download
```

### 3. Установка goose (для миграций):
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 4. Создание .env файла:

Создай файл `.env` в корне проекта:
```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=notes_db
SERVER_PORT=8080
```

### 5. Запуск PostgreSQL:
```bash
docker-compose up -d
```

Проверка:
```bash
docker ps
```

### 6. Применение миграций:
```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=notes_db sslmode=disable" up
```

### 7. Запуск сервера:
```bash
go run cmd/api/main.go
```

Сервер запустится на `http://localhost:8080`

---

## 📋 API Endpoints

### Пользователи:

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| POST | `/users` | Регистрация пользователя | ✅ Реализовано |

### Заметки:

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| POST | `/users/{id}/notes` | Создать заметку | 🔄 В разработке |
| GET | `/users/{id}/notes` | Получить все заметки | 🔄 В разработке |
| GET | `/users/{id}/notes/{note_id}` | Получить одну заметку | 🔄 В разработке |
| PUT | `/users/{id}/notes/{note_id}` | Обновить заметку | 🔄 В разработке |
| DELETE | `/users/{id}/notes/{note_id}` | Удалить заметку | 🔄 В разработке |

---

## 🧪 Примеры использования API

### Создание пользователя:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alex"}'
```

**Ответ:**
```json
{
  "id": 1,
  "username": "alex",
  "created_at": "2025-11-23T12:00:00Z"
}
```

### Валидация (пустой username):
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": ""}'
```

**Ответ:**
```json
{
  "error": "username is required"
}
```

---

## 🧠 Ключевые концепции

### 1. **Миграции БД (Database Migrations)**

Миграции = версионированные изменения схемы базы данных.

**Зачем нужны:**
- Отслеживание изменений структуры БД
- Возможность отката (rollback)
- Синхронизация БД между разработчиками

**Инструмент:** goose

**Пример миграции:**
```sql
-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
```

---

### 2. **Environment Variables (.env)**

**Проблема:** Пароли и настройки не должны быть в коде!

**Решение:** Файл `.env`
```env
DB_PASSWORD=secret123
```

**В коде:**
```go
password := os.Getenv("DB_PASSWORD")
```

---

### 3. **Repository Pattern (Storage Layer)**

**Идея:** Отделить бизнес-логику от работы с БД.
```
Handler → Storage → Database
```

**Преимущества:**
- Весь SQL в одном месте
- Легко тестировать
- Легко менять БД (PostgreSQL → MySQL)

**Пример:**
```go
// Handler не знает про SQL
user, err := h.storage.CreateUser(username)

// Storage содержит SQL
func (s *Storage) CreateUser(username string) (*User, error) {
    query := `INSERT INTO users (username) VALUES ($1) RETURNING id`
    // ...
}
```

---

### 4. **Middleware Pattern**

**Middleware** = функция, которая выполняется до/после handler'а.
```
HTTP Request → [Logger] → [Auth] → Handler → Response
```

**Пример:** Логирование каждого запроса
```go
r.Use(middleware.Logger)  // Автоматически логирует все запросы
```

**Вывод в консоль:**
```
2025/11/23 12:00:00 "POST /users HTTP/1.1" - 201 15ms
```

---

### 5. **Prepared Statements (защита от SQL Injection)**

❌ **Плохо (уязвимо):**
```go
query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
// Если username = "admin'; DROP TABLE users; --" → БД удалится!
```

✅ **Хорошо (безопасно):**
```go
query := `SELECT * FROM users WHERE username = $1`
db.QueryRow(query, username)  // Драйвер экранирует данные
```

---

## 📚 Что изучено в проекте

### Технологии:
- ✅ Структура Go проекта (cmd, internal, pkg)
- ✅ Работа с PostgreSQL через database/sql
- ✅ HTTP роутинг (chi router)
- ✅ JSON encoding/decoding
- ✅ Environment variables
- ✅ Docker Compose
- ✅ Миграции БД (goose)

### Паттерны:
- ✅ Repository Pattern
- ✅ Middleware Pattern
- ✅ DTO (Data Transfer Objects)
- ✅ Error Handling в Go
- ✅ RESTful API design

### SQL:
- ✅ CREATE TABLE
- ✅ INSERT с RETURNING
- ✅ SELECT с WHERE
- ✅ UPDATE
- ✅ DELETE
- ✅ Foreign Keys (REFERENCES)
- ✅ Индексы (CREATE INDEX)
- ✅ Prepared Statements ($1, $2)

---

## 🐛 Troubleshooting

### Docker контейнер не запускается:
```bash
docker-compose down
docker-compose up -d
docker ps
```

### Ошибка подключения к БД:

1. Проверь, запущен ли Docker: `docker ps`
2. Проверь .env файл
3. Измени `DB_HOST=localhost` на `DB_HOST=127.0.0.1`

### Ошибка "relation does not exist":

Примени миграции:
```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=notes_db sslmode=disable" up
```

---

## 🚀 Roadmap (что добавить дальше)

- [ ] Реализовать все endpoints для заметок
- [ ] Добавить аутентификацию (JWT)
- [ ] Написать unit-тесты
- [ ] Добавить Swagger документацию
- [ ] Dockerize приложение (не только БД)
- [ ] CI/CD pipeline
- [ ] Deploy на сервер

---

## 👤 Автор

**Balyshev**
- GitHub: [@Balyshev](https://github.com/Balyshev)

---

## 📝 Лицензия

MIT License — используй как хочешь!

---

## 🎓 Обучение

Этот проект создан в учебных целях для изучения:
- Go backend разработки
- REST API design
- Работы с PostgreSQL
- Архитектурных паттернов