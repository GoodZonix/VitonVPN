# Настройка VPN-серверов (Xray-core, VLESS)

## Требования к VPS

- **ОС**: Ubuntu 22.04 LTS (рекомендуется)
- **Минимум**: 1 vCPU, 512 MB RAM, 10 GB SSD
- **Сеть**: не блокируемый провайдером IP (проверка на цензуру в регионе)

---

## 1. Установка Xray-core

```bash
bash -c "$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)" @ install
systemctl enable xray
systemctl start xray
```

Проверка: `xray version`

---

## 2. Генерация UUID

Один UUID на пользователя или один на сервер (в нашем случае Backend хранит UUID пользователя и подставляет в конфиг).

На сервере для теста:

```bash
xray uuid
# или
cat /proc/sys/kernel/random/uuid
```

В production UUID выдаёт Backend при генерации конфига.

---

## 3. VLESS + Reality (рекомендуется)

Reality маскирует трафик под настоящий TLS-сайт (например, www.google.com).

Генерация ключей Reality (на любой машине с Go):

```bash
xray x25519
# Вывод: Private key, Public key
```

Пример конфигурации Xray на сервере — см. `infra/xray/config-reality.json`.

Параметры:
- **dest**: реальный SNI, под который маскируемся (например, www.google.com:443)
- **serverNames**: список SNI для клиента
- **privateKey / publicKey**: пара x25519
- **shortIds**: короткие ID (можно сгенерировать: `openssl rand -hex 8`)

---

## 4. VLESS + TLS (альтернатива)

Классический TLS с вашим доменом и сертификатом (Let's Encrypt).

- Установка certbot, получение сертификата для домена.
- В конфиге Xray — путь к fullchain.pem и privkey.pem.
- Пример: `infra/xray/config-tls.json`.

---

## 5. Защита от блокировок

- **Reality**: предпочтительно — сложнее для DPI идентифицировать.
- **Порт 443**: использовать стандартный HTTPS-порт.
- **Nginx перед Xray (опционально)**: проксирование по path или заголовку, чтобы на одном IP был и обычный сайт, и VLESS.
- **CDN**: для TLS-варианта можно отдавать трафик через Cloudflare (осторожно с правилами), для Reality — обычно прямой IP.

---

## 6. Firewall

```bash
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp    # SSH
ufw allow 443/tcp   # VLESS
ufw allow 443/udp   # при использовании QUIC
ufw enable
```

---

## 7. Оптимизация скорости

- **Ядро**: `net.core.rmem_max`, `net.core.wmem_max` — увеличить (например, 2500000).
- **BBR**: включить TCP BBR (часто по умолчанию в современных ядрах): `sysctl net.ipv4.tcp_congestion_control=bbr`.
- **Многопоточность**: в Xray для inbounds можно не менять; основная нагрузка — TLS.

---

## 8. Регистрация сервера в Backend

После установки сервер добавляется в БД (через админку или API):

- Адрес (IP или домен)
- Порт (443)
- Регион (EU, US, ASIA, ...)
- Тип (reality / tls)
- Параметры Reality (publicKey, shortId, serverName) для генерации конфига на Backend.

Клиенты получают список серверов через GET /servers и готовые VLESS-ссылки через GET /config.
