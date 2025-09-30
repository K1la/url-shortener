# URL Shortener

Сервис сокращения URL с аналитикой переходов. Создавайте короткие ссылки, отслеживайте статистику и анализируйте переходы пользователей.

## 🚀 Возможности

- **Создание коротких ссылок**: Автоматическая генерация или кастомные алиасы
- **Редирект**: Переходы по коротким ссылкам с сохранением аналитики
- **Аналитика**: Подробная статистика переходов по дням/месяцам
- **Современный UI**: Темная тема с адаптивным дизайном
- **Кэширование**: Redis для быстрого доступа к данным
- **База данных**: PostgreSQL для надежного хранения

## 📋 Требования

- Docker & Docker Compose
- Go 1.24+ (для локальной разработки)
- Node.js (опционально, для фронтенда)

## 🏗️ Архитектура проекта

```
url-shortener/
├── cmd/
│   └── shortener/
│       └── main.go                 # Точка входа приложения
├── internal/
│   ├── api/
│   │   ├── handler/                # HTTP обработчики
│   │   │   ├── create.go          # Создание коротких ссылок
│   │   │   ├── get.go             # Получение и редирект
│   │   │   ├── handler.go         # Базовый обработчик
│   │   │   └── interface.go       # Интерфейсы
│   │   ├── response/              # HTTP ответы
│   │   │   └── response.go
│   │   ├── router/                # Маршрутизация
│   │   │   └── router.go
│   │   └── server/                # HTTP сервер
│   │       └── server.go
│   ├── cache/                     # Кэширование
│   │   └── redis.go
│   ├── config/                    # Конфигурация
│   │   ├── config.go
│   │   └── types.go
│   ├── model/                     # Модели данных
│   │   └── model.go
│   ├── repository/                # Слой данных
│   │   ├── create.go             # Создание записей
│   │   ├── get.go                # Получение данных
│   │   └── repo.go               # Базовый репозиторий
│   └── service/                   # Бизнес-логика
│       ├── create.go             # Создание ссылок
│       ├── get.go                # Получение ссылок
│       ├── interface.go          # Интерфейсы сервисов
│       ├── service.go            # Базовый сервис
│       └── utils.go              # Утилиты
├── migrations/                    # SQL миграции
│   ├── 20250929154354_create_table_urls.sql
│   └── 20250929155733_create_table_redirect_analytics.sql
├── web/                          # Фронтенд
│   ├── index.html                # Главная страница
│   ├── analytics.html            # Страница аналитики
│   ├── app.js                    # JS для главной страницы
│   ├── analytics.js              # JS для аналитики
│   └── styles.css                # Стили
├── docker-compose.yml            # Docker Compose конфигурация
├── Dockerfile                    # Docker образ
├── go.mod                        # Go модули
└── go.sum                        # Go зависимости
```

## 🛠️ Технологии

### Backend
- **Go 1.24** - основной язык
- **Gin** - HTTP фреймворк
- **PostgreSQL** - основная база данных
- **Redis** - кэширование
- **Goose** - миграции БД

### Frontend
- **Vanilla JavaScript** - без фреймворков
- **CSS3** - современные стили
- **Responsive Design** - адаптивная верстка

## 🚀 Быстрый старт

### 1. Клонирование репозитория
```bash
git clone https://github.com/K1la/url-shortener
cd url-shortener
```

### 2. Запуск с Docker Compose
```bash
# Создание .env файла (опционально)
cat > .env << EOF
DB_HOST=db
DB_PORT=5432
DB_USER=your_name
DB_PASSWORD=your_password
DB_NAME=postgres-url-shortener
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=/migrations
REDIS_PASSWORD=redispass
EOF

# Запуск всех сервисов
docker compose up --build
```

### 3. Доступ к приложению
- **Веб-интерфейс**: http://localhost:8080
- **API**: http://localhost:8080/api/
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 📚 API Документация

### Создание короткой ссылки
```http
POST /api/shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url",
  "user_short_url": "my-custom-alias"  // опционально
}
```

**Ответ:**
```json
{
  "result": {
    "id": "uuid",
    "url": "https://example.com/very/long/url",
    "short_url": "abc123",
    "created_at": "2025-09-30T10:00:00Z"
  }
}
```

### Переход по короткой ссылке
```http
GET /s/{short_url}
```
**Ответ:** HTTP 302 редирект на оригинальный URL

### Аналитика переходов
```http
GET /api/analytics/{short_url}
```

**Ответ:**
```json
{
  "result": {
    "short_url": "abc123",
    "total_clicks": 42,
    "daily": {
      "2025-09-30": 15,
      "2025-09-29": 27
    },
    "monthly": {
      "2025-09": 42
    },
    "user_agent": {
      "Mozilla/5.0...": 25,
      "Chrome/91.0...": 17
    }
  }
}
```

## 🗄️ База данных

### Таблица `urls`
```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    short_url VARCHAR(32) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Таблица `analytics`
```sql
CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(32) REFERENCES urls(short_url) ON DELETE CASCADE,
    user_agent TEXT,
    device_type VARCHAR(32),
    os VARCHAR(64),
    browser VARCHAR(64),
    ip_address INET,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## 🎨 Фронтенд

### Главная страница (`/`)
- Создание коротких ссылок
- Кастомные алиасы
- Копирование ссылок
- Переход к аналитике

### Страница аналитики (`/analytics.html`)
- Ввод алиаса для анализа
- Статистика переходов
- Графики по дням/месяцам
- Таблица с детальной информацией


### Docker Compose сервисы
- **shortener** - основное приложение
- **db** - PostgreSQL база данных
- **redis** - Redis кэш
- **migrator** - выполнение миграций

## 🚀 Развертывание

### Продакшн
1. Настройте переменные окружения
2. Используйте внешние БД (PostgreSQL, Redis)
3. Настройте SSL/TLS
4. Используйте reverse proxy (nginx)

### Мониторинг
- Логи приложения в stdout
- Метрики PostgreSQL
- Мониторинг Redis
- Health checks для Docker

## 🧪 Тестирование

### Локальная разработка
```bash
# Запуск только БД
docker compose up db redis

# Запуск приложения локально
go run cmd/shortener/main.go
```

### Тестирование API
```bash
# Создание ссылки
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Получение аналитики
curl http://localhost:8080/api/analytics/abc123
```

