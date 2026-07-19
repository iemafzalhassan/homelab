ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
ALTER TABLE users ALTER COLUMN keycloak_sub SET NOT NULL;
DROP INDEX IF EXISTS users_keycloak_sub_idx;
ALTER TABLE users ADD CONSTRAINT users_keycloak_sub_key UNIQUE (keycloak_sub);