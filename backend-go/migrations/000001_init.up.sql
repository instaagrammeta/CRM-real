-- ====================================================================
-- 000001_init.up.sql
-- Manual SQL migration; the application also runs GORM AutoMigrate on
-- startup so this file is mainly a snapshot/backup for ops.
-- ====================================================================

CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    full_name       VARCHAR(255),
    age             INTEGER,
    personal_phones JSONB,
    work_phones     JSONB,
    login           VARCHAR(128) UNIQUE,
    password        VARCHAR(255),
    photo           VARCHAR(512),
    category        VARCHAR(128),
    role            VARCHAR(32) DEFAULT 'employee',
    telegram_chat_id BIGINT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_telegram_chat_id ON users(telegram_chat_id);

CREATE TABLE IF NOT EXISTS tasks (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255),
    description  TEXT,
    author_id    INTEGER,
    executor_id  INTEGER,
    photo        VARCHAR(512),
    status       VARCHAR(32) DEFAULT 'new',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER,
    title       VARCHAR(255),
    body        TEXT,
    type        VARCHAR(64),
    entity_type VARCHAR(64),
    entity_id   INTEGER,
    link        VARCHAR(512),
    is_read     BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);

CREATE TABLE IF NOT EXISTS telegram_subscribers (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER,
    chat_id    BIGINT UNIQUE,
    username   VARCHAR(128),
    first_name VARCHAR(128),
    last_name  VARCHAR(128),
    is_active  BOOLEAN DEFAULT TRUE,
    link_token VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Other tables are created by GORM AutoMigrate.
