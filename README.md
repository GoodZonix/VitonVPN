# VPN Startup — VLESS (Xray-core)

Полноценный VPN-сервис уровня стартапа на протоколе VLESS (Xray-core) с поддержкой Reality/TLS, backend на Go, мобильными приложениями (Flutter) и админ-панелью.

## Структура репозитория

```
vpn-startup/
├── ARCHITECTURE.md          # Общая архитектура системы
├── docs/                    # Документация
│   ├── SERVERS.md           # Настройка VPS и Xray
│   ├── SECURITY.md          # Безопасность и анти-абьюз
│   └── LAUNCH_PLAN.md       # План запуска и масштабирование
├── infra/                   # Инфраструктура
│   ├── xray/                # Конфиги Xray (VLESS, Reality, TLS)
│   ├── docker/              # Dockerfile для backend, xray (если нужно)
│   └── k8s/                 # Манифесты Kubernetes (опционально)
├── backend/                 # Backend API (Go)
├── admin/                   # Панель администратора (Next.js)
├── mobile/                  # Мобильное приложение (Flutter)
└── scripts/                 # CI/CD, деплой, генерация конфигов
```

## Быстрый старт

- **Серверы**: см. `docs/SERVERS.md` и `infra/xray/`
- **Backend**: `cd backend && go run ./cmd/api`
- **Admin**: `cd admin && pnpm dev`
- **Mobile**: `cd mobile && flutter run`

## Стек

| Компонент   | Технологии                          |
|------------|--------------------------------------|
| VPN ядро   | Xray-core, VLESS, Reality / TLS      |
| Backend    | Go, PostgreSQL, Redis, REST API     |
| Admin      | Next.js, React, Admin UI             |
| Mobile     | Flutter (iOS + Android)              |
| DevOps     | Docker, GitHub Actions, опционально K8s |

## Лицензия

Proprietary — стартап.
