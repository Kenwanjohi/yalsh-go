ALTER TABLE links DROP CONSTRAINT IF EXISTS links_key_key;
CREATE UNIQUE INDEX links_unique_key_active ON links (key) WHERE deleted_at IS NULL;