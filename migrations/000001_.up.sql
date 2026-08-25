CREATE SCHEMA IF NOT EXISTS auth_msginx;

CREATE TABLE IF NOT EXISTS auth_msginx.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(512) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT check_login_length CHECK (LENGTH(login) >= 3),
    CONSTRAINT check_role CHECK (role IN ('user', 'admin')),
    CONSTRAINT check_dates CHECK (created_at <= updated_at)
);

CREATE TABLE IF NOT EXISTS auth_msginx.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    refresh_token_hash BYTEA NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES auth_msginx.users(id) ON DELETE CASCADE,
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT check_dates CHECK (created_at < expires_at)
);

CREATE INDEX idx_sessions_user_id ON auth_msginx.sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON auth_msginx.sessions(expires_at);

CREATE OR REPLACE FUNCTION auth_msginx.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON auth_msginx.users
    FOR EACH ROW
    EXECUTE FUNCTION auth_msginx.update_updated_at_column();

COMMENT ON TABLE auth_msginx.users IS 'Таблица пользователей';
COMMENT ON COLUMN auth_msginx.users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN auth_msginx.users.login IS 'Логин пользователя';
COMMENT ON COLUMN auth_msginx.users.password_hash IS 'Хэш пароля';
COMMENT ON COLUMN auth_msginx.users.role IS 'Роль пользователя: user или admin';
COMMENT ON COLUMN auth_msginx.users.created_at IS 'Время регистрации пользователя';
COMMENT ON COLUMN auth_msginx.users.updated_at IS 'Время изменения данных';

COMMENT ON TABLE auth_msginx.sessions IS 'Таблица сессий пользователей';
COMMENT ON COLUMN auth_msginx.sessions.id IS 'Уникальный идентификатор сессии';
COMMENT ON COLUMN auth_msginx.sessions.refresh_token_hash IS 'Хэш refresh токен';
COMMENT ON COLUMN auth_msginx.sessions.user_id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN auth_msginx.sessions.is_revoked IS 'Флаг отзыва сессии';
COMMENT ON COLUMN auth_msginx.sessions.created_at IS 'Время создания сессии';
COMMENT ON COLUMN auth_msginx.sessions.expires_at IS 'Время истечения сессии';
