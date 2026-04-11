DROP INDEX IF EXISTS links_unique_key_active;
ALTER TABLE links ADD CONSTRAINT links_key_key UNIQUE (key);