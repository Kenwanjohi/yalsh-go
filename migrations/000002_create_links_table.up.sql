CREATE TABLE IF NOT EXISTS links (
  id bigserial PRIMARY KEY,
  key varchar(20) UNIQUE NOT NULL,
  url text NOT NULL,
  clicks integer DEFAULT 0 NOT NULL,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at timestamptz,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_links_key ON links(key);