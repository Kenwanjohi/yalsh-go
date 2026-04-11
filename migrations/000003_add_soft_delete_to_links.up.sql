ALTER TABLE links ADD COLUMN deleted_at timestamptz;
ALTER TABLE links ADD COLUMN version integer NOT NULL DEFAULT 1;