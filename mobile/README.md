# VPN Startup — Mobile App (Flutter)

Мобильное приложение для подключения к VPN (VLESS / Xray-core).

## Функции

- Большая кнопка Connect / Disconnect
- Выбор региона и сервера (список с бэкенда)
- Автоматический лучший сервер (по пингу)
- Статус подключения и пинг
- Автоматическое переподключение (при обрыве)
- Защита от утечки DNS (настройки + конфиг с backend)
- Тёмная и светлая тема
- Регистрация / логин, проверка подписки, экран подписки (IAP)

## Структура

```
lib/
├── main.dart
├── app.dart
├── config/          # API URL, theme
├── core/             # API client, auth storage
├── models/           # Server, User, Subscription
├── providers/        # Auth, VPN state, Servers
├── screens/
│   ├── home_screen.dart
│   ├── servers_screen.dart
│   ├── subscription_screen.dart
│   ├── login_screen.dart
│   └── settings_screen.dart
└── widgets/          # Connect button, server tile
```

## Запуск

```bash
flutter pub get
flutter run
```

Для production: интеграция нативного VPN через Method Channel с Xray-core или использование готового пакета (например, обёртка над v2ray-core). Конфиг VLESS получается с GET /api/config и передаётся в движок VPN.

## Платформы

- iOS: требуется Network Extension (Packet Tunnel Provider) или использование системного VPN API с конфигом.
- Android: VpnService + локальный прокси (Xray в процессе приложения) или системный VPN с конфигом.

Реализация туннеля зависит от выбора библиотеки (например, libv2ray, xray-core для мобильных).
