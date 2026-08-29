ALTER TABLE IF EXISTS auth_msginx.sessions
    DROP CONSTRAINT IF EXISTS check_dates;

ALTER TABLE IF EXISTS auth_msginx.users 
    DROP CONSTRAINT IF EXISTS check_login_length,
    DROP CONSTRAINT IF EXISTS check_role,
    DROP CONSTRAINT IF EXISTS check_dates,


DROP INDEX IF EXISTS auth_msginx.idx_sessions_user_id;
DROP INDEX IF EXISTS auth_msginx.idx_sessions_expires_at;

DROP TRIGGER IF EXISTS auth_msginx.update_users_updated_at;
DROP FUNCTION IF EXISTS auth_msginx.update_updated_at_column;

DROP SCHEMA IF EXISTS auth_msginx CASCADE;