-- Users
CREATE TABLE IF NOT EXISTS users (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email            VARCHAR(255) UNIQUE NOT NULL,
    password_hash    VARCHAR(255) NOT NULL,
    wallet_balance   NUMERIC(12,2) NOT NULL DEFAULT 0,
    last_billing_at  TIMESTAMPTZ DEFAULT NOW(),
    trial_expires_at TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- User devices (for device limit)
CREATE TABLE IF NOT EXISTS devices (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id   VARCHAR(255) NOT NULL,
    name        VARCHAR(100),
    last_seen   TIMESTAMPTZ DEFAULT NOW(),
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, device_id)
);

CREATE INDEX idx_devices_user_id ON devices(user_id);

-- VPN servers (edge nodes)
CREATE TABLE IF NOT EXISTS servers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    region          VARCHAR(50) NOT NULL,
    host            VARCHAR(255) NOT NULL,
    port            INT NOT NULL DEFAULT 443,
    type            VARCHAR(20) NOT NULL DEFAULT 'reality',
    reality_pub_key VARCHAR(255),
    reality_short_id VARCHAR(50),
    reality_sni     VARCHAR(255),
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_servers_region ON servers(region);
CREATE INDEX idx_servers_active ON servers(is_active);

-- Subscriptions (plans: 1m, 3m, 12m)
CREATE TABLE IF NOT EXISTS subscriptions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan            VARCHAR(10) NOT NULL,
    started_at      TIMESTAMPTZ NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    auto_renew      BOOLEAN DEFAULT false,
    external_id     VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_expires ON subscriptions(expires_at);

-- User VPN config: one UUID per user (used in VLESS URL)
CREATE TABLE IF NOT EXISTS user_vpn_keys (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    vless_uuid  UUID NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Traffic stats (optional, for admin)
CREATE TABLE IF NOT EXISTS traffic_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id),
    server_id   UUID REFERENCES servers(id),
    bytes_up    BIGINT DEFAULT 0,
    bytes_down  BIGINT DEFAULT 0,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_traffic_logs_user_created ON traffic_logs(user_id, created_at);

-- Wallet topups (e.g. YooMoney)
CREATE TABLE IF NOT EXISTS wallet_topups (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount       NUMERIC(12,2) NOT NULL,
    provider     VARCHAR(50) NOT NULL,
    provider_tx  VARCHAR(255) NOT NULL,
    status       VARCHAR(20) NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_wallet_topups_user_id ON wallet_topups(user_id);
CREATE INDEX idx_wallet_topups_provider_tx ON wallet_topups(provider_tx);

-- Wallet usage records (optional, for analytics/support)
CREATE TABLE IF NOT EXISTS wallet_usage (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id),
    server_id   UUID REFERENCES servers(id),
    amount      NUMERIC(12,4) NOT NULL,
    started_at  TIMESTAMPTZ NOT NULL,
    ended_at    TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_wallet_usage_user_created ON wallet_usage(user_id, created_at);

-- Telegram account linking
CREATE TABLE IF NOT EXISTS telegram_accounts (
    tg_user_id BIGINT PRIMARY KEY,
    user_id    UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telegram_accounts_user_id ON telegram_accounts(user_id);
