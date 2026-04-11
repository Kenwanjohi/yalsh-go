CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  username varchar(50) NOT NULL,
  email varchar(100) UNIQUE NOT NULL,
  password bytea NOT NULL,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);