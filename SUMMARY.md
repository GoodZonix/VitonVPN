# VPN-сервис на VLESS (Xray-core) — Итог

Результат работы команды senior-разработчиков: архитектура, кодовая структура и примеры реализации полноценного VPN-стартапа уровня hitVPN.

---

## 1. Архитектура VPN

- **Протокол**: VLESS поверх Xray-core.
- **Шифрование/маскировка**: Reality (рекомендуется) или TLS.
- **Схема**: клиенты (iOS/Android/Web) → Backend API (конфиг + список серверов) → VPN edge (Xray на VPS) → интернет. Обход DPI за счёт Reality/TLS и порта 443.
- **Конфигурация**: клиент получает её через **GET /api/config** (JWT + активная подписка). Backend генерирует VLESS URL по пользовательскому UUID и списку серверов.
- **Распределение серверов**: по регионам (EU, US, Asia); выбор вручную или «Auto» (лучший по пингу). Масштабирование — добавление новых VPS и записей в БД.

Подробно: **ARCHITECTURE.md**.

---

## 2. VPN-серверы

- Документация по установке Xray, генерации UUID, настройке VLESS, Reality и TLS: **docs/SERVERS.md**.
- Примеры конфигов Xray: **infra/xray/config-reality.json**, **infra/xray/config-tls.json**.
- Описаны: firewall, оптимизация (BBR, буферы), защита от блокировок.

---

## 3. Backend

- **Стек**: Go, Chi, PostgreSQL, Redis, JWT.
- **Назначение**: регистрация/логин, лимит устройств, генерация VPN-конфигов, список серверов, подписки, проверка срока подписки.
- **Структура**: `backend/cmd/api`, `internal/config`, `internal/database`, `internal/models`, `internal/auth`, `internal/repository`, `internal/vless`, `internal/api/handlers`, `internal/api/middleware`.
- **БД**: схема в **backend/internal/database/schema.sql** (users, devices, servers, subscriptions, user_vpn_keys, traffic_logs).
- **API**: **docs/API.md** — POST /api/login, POST /api/register, GET /api/servers, GET /api/config, POST /api/subscribe, GET /api/subscribe.

---

## 4. Мобильное приложение (Flutter)

- **Расположение**: **mobile/**.
- **Функции**: большая кнопка Connect/Disconnect, выбор региона/сервера, авто-сервер, статус и пинг, автоматическое переподключение (заготовка), защита от утечки DNS (настройка + конфиг с backend).
- **Экраны**: главный (Home), логин/регистрация, список серверов, подписка, настройки. Провайдеры: Auth, VPN, Servers.
- **Тема**: светлая/тёмная (config/theme.dart). Для полноценного туннеля нужна интеграция с нативным Xray/V2Ray (Method Channel или готовый пакет).

---

## 5. UI/UX

- Описание стиля и экранов: **docs/UI_UX.md** (минимализм, темы, главный экран, выбор серверов, подписка, уникальность дизайна).

---

## 6. Платная подписка

- Планы: 1 месяц, 3 месяца, 12 месяцев. Интеграция Apple IAP и Google Play Billing — **docs/BILLING.md**.
- Backend: проверка подписки в GET /api/config и GET /api/subscribe, отключение VPN при истечении (логика на клиенте по ответу API), поддержка автопродления через webhook/серверные уведомления.

---

## 7. Панель администратора

- **Стек**: Next.js 14 (App Router), React, SWR.
- **Расположение**: **admin/**.
- **Страницы**: главная (навигация), серверы (таблица из GET /api/servers), пользователи, подписки, статистика — заглушки с указанием на необходимость admin API в backend.

---

## 8. DevOps

- **Docker**: **docker-compose.yml** — API, PostgreSQL, Redis; **infra/docker/Dockerfile.backend** для backend.
- **CI**: **.github/workflows/ci.yml** — сборка и тесты backend, сборка admin.
- **Деплой**: **docs/DEVOPS.md** — Docker, опционально Kubernetes, балансировка, автоматическое добавление серверов.

---

## 9. Безопасность

- **docs/SECURITY.md**: защита API (JWT, rate limit, CORS), защита конфигов, анти-абьюз, лимит устройств.

---

## 10. Результат — чек-лист

| Пункт | Где |
|-------|-----|
| Полная архитектура | ARCHITECTURE.md, README.md |
| Структура репозитория | Корень vpn-startup/, README |
| Пример конфигурации VLESS | infra/xray/config-reality.json, config-tls.json, vless.BuildVLESSURL в backend |
| Backend API | backend/, docs/API.md |
| Примеры кода | backend (Go), mobile (Flutter), admin (Next.js) |
| UI макеты/описание приложения | docs/UI_UX.md, mobile/lib/screens/ |
| План запуска VPN-стартапа | docs/LAUNCH_PLAN.md |
| Рекомендации по масштабированию | docs/LAUNCH_PLAN.md, docs/DEVOPS.md |

---

## Быстрый старт

```bash
# Backend (нужны PostgreSQL и Redis; или docker compose up)
cd backend && go mod tidy && go run ./cmd/api

# Admin
cd admin && pnpm install && pnpm dev

# Mobile
cd mobile && flutter pub get && flutter run
```

Переменные окружения backend: `DATABASE_URL`, `REDIS_URL`, `JWT_SECRET`, `HTTP_PORT`, `MAX_DEVICES_PER_USER`.
