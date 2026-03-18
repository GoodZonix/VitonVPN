# Платная подписка и монетизация

## Планы

| План | Срок   | Рекомендуемая цена | Скидка   |
|------|--------|---------------------|----------|
| 1m   | 1 мес  | \$4.99              | —        |
| 3m   | 3 мес  | \$12.99             | ~13%     |
| 12m  | 12 мес | \$39.99             | ~33%     |

Backend хранит в `subscriptions` поля: `plan`, `started_at`, `expires_at`, `external_id` (идентификатор покупки в магазине).

## Apple In-App Purchases

1. В App Store Connect создать In-App Purchase продукты (consumable или auto-renewable subscription).
2. В приложении (Flutter) использовать пакет `in_app_purchase` для запроса покупки и получения транзакции.
3. После успешной покупки отправить на backend серверную верификацию (receipt) и вызвать POST /api/subscribe с `plan` и `external_id` (transactionId или productId).
4. Backend при необходимости верифицирует receipt через Apple API и создаёт/продлевает запись в `subscriptions`.

## Google Play Billing

1. В Google Play Console создать подписки (monthly, 3-month, annual).
2. В приложении использовать `in_app_purchase` (Flutter) или `billing_client` для Android.
3. После покупки — верификация через Google Play Developer API (backend) и вызов POST /api/subscribe с `external_id` (orderId/purchaseToken).

## Backend

- **Проверка подписки**: при GET /api/config и GET /api/subscribe проверять `subscriptions.expires_at > NOW()`.
- **Отключение VPN при истечении**: клиент при каждом подключении или по таймеру запрашивает статус; при `active: false` — разрывать туннель и показывать экран продления.
- **Автопродление**: для auto-renewable (Apple/Google) настроить webhook/серверные уведомления (Real-time developer notifications от Google, App Store Server Notifications от Apple). Backend при получении события продления обновляет `expires_at` в `subscriptions` (или добавляет новый период).

## Рекомендации

- Хранить в backend только факт подписки и срок; детали транзакций при необходимости хранить в логах или отдельной таблице для поддержки.
- Для тестирования использовать sandbox окружения Apple/Google.
