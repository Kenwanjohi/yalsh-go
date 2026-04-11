ALTER TABLE users DROP COLUMN IF EXISTS version;
ALTER TABLE users ALTER COLUMN password_hash TYPE bytea USING password_hash::bytea;
ALTER TABLE users RENAME COLUMN password_hash TO password;