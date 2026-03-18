# DevOps: деплой и инфраструктура

## Docker

- **Backend**: образ собирается из `backend/` с Dockerfile в `infra/docker/Dockerfile.backend`. Контекст сборки — `./backend`.
- **PostgreSQL и Redis**: образы из Docker Hub; схема БД подключается как init-скрипт в `docker-compose.yml`.

### Локальный запуск

```bash
# Из корня репозитория
docker compose up -d
# API: http://localhost:8080
# Подключение к БД: localhost:5432, user vpn, password vpn, db vpndb
# Redis: localhost:6379
```

Миграции: при первом запуске `schema.sql` выполняется автоматически (volume docker-entrypoint-initdb.d). Для последующих изменений схемы используйте миграции (например, golang-migrate) или ручное применение.

## Kubernetes (опционально)

- **Backend**: Deployment с 2–3 репликами, Service (ClusterIP или LoadBalancer), ConfigMap/Secret для `DATABASE_URL`, `REDIS_URL`, `JWT_SECRET`.
- **PostgreSQL/Redis**: можно вынести в managed-сервисы (RDS, ElastiCache, Cloud SQL, etc.) или развернуть в кластере (StatefulSet для PostgreSQL).
- **VPN-ноды (Xray)**: не в K8s; отдельные VPS с установленным Xray-core. Регистрация нод — через админку/БД.

Пример манифеста (фрагмент):

```yaml
# k8s/api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vpn-api
spec:
  replicas: 2
  selector:
    matchLabels: { app: vpn-api }
  template:
    metadata:
      labels: { app: vpn-api }
    spec:
      containers:
        - name: api
          image: your-registry/vpn-api:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: vpn-secrets
                  key: database-url
            - name: REDIS_URL
              valueFrom:
                secretKeyRef:
                  name: vpn-secrets
                  key: redis-url
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: vpn-secrets
                  key: jwt-secret
```

## CI/CD

- **GitHub Actions** (`.github/workflows/ci.yml`):
  - На push/PR в main/develop: сборка backend (`go build`), тесты (`go test`), сборка admin (`pnpm build`).
- Деплой: добавить job с деплоем в облако (e.g. build Docker-образа, push в registry, обновление K8s или ECS).

## Автоматическое добавление серверов

- Скрипт (например, `scripts/register-server.sh`) вызывает внутренний admin API или напрямую вставляет запись в таблицу `servers` (host, region, port, type, reality_pub_key, reality_short_id, reality_sni).
- На каждом VPS после установки Xray запускается скрипт с параметрами (или через Ansible/Terraform), который регистрирует ноду в БД.

## Балансировка нагрузки

- **API**: перед несколькими инстансами backend — load balancer (nginx, cloud LB, Ingress). Сессии не привязаны к инстансу (JWT stateless; Redis общий).
- **VPN-трафик**: клиент выбирает один сервер из списка (вручную или «Auto»). Балансировка — за счёт выбора пользователем или логики «лучший по пингу» на клиенте. При необходимости можно отдавать приоритет серверам с меньшей нагрузкой (метрики с нод в Redis, GET /servers с сортировкой).
