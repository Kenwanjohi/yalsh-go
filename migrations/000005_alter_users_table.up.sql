ALTER TABLE users RENAME COLUMN password TO password_hash;
ALTER TABLE users ALTER COLUMN password_hash TYPE text USING password_hash::text;
ALTER TABLE users ADD COLUMN version integer NOT NULL DEFAULT 1;