ALTER TABLE users ADD COLUMN password_hash TEXT;
ALTER TABLE users ALTER COLUMN keycloak_sub DROP NOT NULL;
ALTER TABLE users ALTER COLUMN keycloak_sub DROP UNIQUE;
CREATE UNIQUE INDEX IF NOT EXISTS users_keycloak_sub_idx ON users(keycloak_sub) WHERE keycloak_sub IS NOT NULL;