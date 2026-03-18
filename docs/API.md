# Backend REST API

Base URL: `https://api.yourapp.com` (или локально `http://localhost:8080`).

## Публичные эндпоинты

### POST /api/register

Регистрация пользователя.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "device_id": "unique-device-uuid",
  "device_name": "iPhone 15"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 604800,
  "user": { "id": "uuid", "email": "user@example.com" }
}
```

### POST /api/login

Вход.

**Request:** как у register (email, password, device_id, device_name).

**Response:** как у register.

---

## С авторизацией (Header: Authorization: Bearer &lt;token&gt;)

### GET /api/servers

Список серверов (доступен и без токена для выбора региона).

**Response (200):**
```json
{
  "servers": [
    {
      "id": "uuid",
      "name": "Netherlands 1",
      "region": "EU",
      "host": "nl1.vpn.example.com",
      "port": 443,
      "type": "reality",
      "is_active": true
    }
  ]
}
```

### GET /api/config

Получить VLESS-конфигурацию и список серверов. Проверяется подписка: если нет активной — vless_urls пустой.

**Response (200):**
```json
{
  "vless_urls": [
    "vless://user-uuid@nl1.vpn.example.com:443?type=tcp&security=reality&pbk=...&sid=...&sni=...&fp=chrome#Netherlands%201"
  ],
  "servers": [ ... ],
  "subscription": {
    "plan": "1m",
    "expires_at": "2025-04-14T12:00:00Z",
    "active": true
  }
}
```

### POST /api/subscribe

Оформить/продлить подписку (после успешной оплаты IAP/Google Play backend вызывает или клиент передаёт external_id).

**Request:**
```json
{
  "plan": "1m",
  "external_id": "iap_purchase_id_or_play_order_id"
}
```

**Response (200):**
```json
{
  "subscription": {
    "id": "uuid",
    "plan": "1m",
    "expires_at": "2025-04-14T12:00:00Z",
    "active": true
  }
}
```

### GET /api/subscribe

Статус подписки.

**Response (200):**
```json
{
  "active": true,
  "subscription": {
    "id": "uuid",
    "plan": "12m",
    "expires_at": "2026-03-14T12:00:00Z",
    "active": true
  }
}
```

---

## Коды ошибок

- **400** — неверное тело запроса или параметры.
- **401** — нет/неверный JWT.
- **403** — доступ запрещён (например, превышен лимит устройств).
- **409** — конфликт (email уже занят).
- **429** — rate limit (слишком много запросов).
