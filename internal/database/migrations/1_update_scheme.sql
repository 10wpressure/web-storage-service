BEGIN;

-- Добавить ограничения длины для login и password_hash
ALTER TABLE users
    ADD CONSTRAINT login_length CHECK (length(login) <= 255),
    ADD CONSTRAINT password_hash_length CHECK (length(password_hash) <= 255);

-- Добавить новые колонки в таблицу сессий
ALTER TABLE sessions
    ADD COLUMN IF NOT EXISTS ip_address INET NOT NULL DEFAULT '0.0.0.0',
    ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 HOURS',
    ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT FALSE;

-- Установить ссылочную целостность для таблицы сессий и таблицы файлов
ALTER TABLE sessions
    ADD CONSTRAINT fk_sessions_users FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE assets
    ADD COLUMN IF NOT EXISTS deleted BOOLEAN NOT NULL DEFAULT FALSE,
    ADD CONSTRAINT fk_assets_users FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE;

-- Удалить временное значение по умолчанию из ip_address
ALTER TABLE sessions ALTER COLUMN ip_address DROP DEFAULT;

COMMIT;
