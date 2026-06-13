```markdown
# Ozinse Backend

API для стриминговой платформы Ozinse на Go с использованием Gin.

## Стек

- **Go 1.23+**
- **Gin** — HTTP фреймворк
- **PostgreSQL** — база данных
- **JWT** — авторизация
- **Docker** — контейнеризация
- **Swagger** — документация API

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/AdWayne/ozinse.git
cd ozinse
```

### 2. Настроить .env

Создай `.env` файл в корне:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123
DB_NAME=ozinse
JWT_SECRET=super_secret_key
REFRESH_SECRET=super_secret_refresh_key
BASE_URL=http://localhost:8080
```

### 3. Запустить через Docker

```bash
docker-compose up --build
```

Сервисы:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8081
- **PostgreSQL**: localhost:5432

### 4. Запустить без Docker

```bash
# Установить зависимости
go mod download

# Запустить PostgreSQL и выполнить миграции вручную
# migrations/001_init.sql — создание таблиц
# migrations/002_insert.sql — начальные данные

# Запустить сервер
go run cmd/api/main.go
```

## Учётные записи по умолчанию

| Роль | Email | Пароль |
|------|-------|--------|
| Админ | admin@ozinse.com | admin123 |
| Пользователь | (регистрация) | — |

## API Endpoints

### Авторизация
| Метод | Путь | Описание |
|-------|------|----------|
| POST | /api/v1/auth/register | Регистрация |
| POST | /api/v1/auth/login | Вход |
| POST | /api/v1/auth/refresh | Обновление токенов |
| POST | /api/v1/auth/reset-password | Восстановление пароля |

### Профиль (🔒)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/v1/profile/me | Получить профиль |
| PUT | /api/v1/profile/me | Обновить профиль |
| PUT | /api/v1/profile/password | Сменить пароль |

### Публичный контент
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/v1/projects | Список проектов (с фильтрами) |
| GET | /api/v1/projects/featured | Подборки для главной |
| GET | /api/v1/projects/{id} | Детали проекта |
| GET | /api/v1/projects/{id}/seasons | Сезоны и эпизоды |

### Избранное (🔒)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/v1/favorites | Список избранного |
| POST | /api/v1/favorites/{project_id} | Добавить в избранное |
| DELETE | /api/v1/favorites/{project_id} | Удалить из избранного |

### Справочники
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/v1/categories | Категории |
| GET | /api/v1/genres | Жанры |
| GET | /api/v1/age-ratings | Возрастные рейтинги |

### Админка (🔒🛡️)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/v1/admin/projects | Все проекты |
| POST | /api/v1/admin/projects | Создать проект |
| PUT | /api/v1/admin/projects/{id} | Обновить проект |
| DELETE | /api/v1/admin/projects/{id} | Удалить проект |
| POST | /api/v1/admin/projects/{id}/seasons | Добавить сезон |
| POST | /api/v1/admin/seasons/{id}/episodes | Добавить эпизод |
| PUT | /api/v1/admin/projects/featured-order | Порядок в подборках |
| POST/PUT/DELETE | /api/v1/admin/categories | Управление категориями |
| POST/PUT/DELETE | /api/v1/admin/genres | Управление жанрами |
| POST/PUT/DELETE | /api/v1/admin/age-ratings | Управление рейтингами |
| GET | /api/v1/admin/users | Список пользователей |
| POST | /api/v1/admin/users/{id}/assign-role | Назначить роль |
| GET/POST/PUT/DELETE | /api/v1/admin/roles | Управление ролями |
| POST | /api/v1/admin/upload | Загрузка файлов |

## Структура проекта

```
ozinse/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа
├── internal/
│   ├── config/                  # Конфигурация
│   ├── database/                # Подключение к БД и миграции
│   ├── model/                   # Структуры данных
│   ├── repository/              # SQL запросы
│   ├── service/                 # Бизнес-логика
│   └── handler/                 # HTTP обработчики
├── pkg/
│   └── jwt/                     # JWT сервис
├── migrations/                  # SQL миграции
├── static/                      # Статические файлы
├── uploads/                     # Загруженные файлы
├── swagger.yaml                 # Swagger документация
├── Dockerfile
├── docker-compose.yml
└── .env
```

## Фильтры для GET /projects

| Параметр | Тип | Пример |
|----------|-----|--------|
| search | string | ?search=ғарыш |
| category_id | int | ?category_id=3 |
| genre_id | int | ?genre_id=5 |
| age_rating_id | int | ?age_rating_id=3 |
| project_type | string | ?project_type=MOVIE |
| page | int | ?page=1 |
| limit | int | ?limit=20 |

## Ошибки

Формат ответа при ошибке:

```json
{
    "error_code": "NOT_FOUND",
    "message": "Проект не найден",
    "details": null
}
```

HTTP коды: 200, 201, 400, 401, 403, 404, 500.

## Документация API

Swagger UI доступен по адресу http://localhost:8081 (при запуске через Docker).

Или открой `swagger.yaml` в [Swagger Editor](https://editor.swagger.io).
