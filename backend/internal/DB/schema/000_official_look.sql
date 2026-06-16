CREATE TABLE IF NOT EXISTS Users (
  user_id UUID PRIMARY KEY,
  email text,
  username text NOT NULL,
  password_hash VARCHAR(72) NOT NULL,
  date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  donation_percentage INT DEFAULT 10
);

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entry') THEN
        CREATE TYPE entry AS ENUM ('paycheck', 'donation');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS Ledgers (
  user_id UUID NOT NULL,
  transaction_id SERIAL PRIMARY KEY,
  ledger_entry entry,
  amount DECIMAL(18, 2) NOT NULL,
  description TEXT,
  charity_owed DECIMAL(18, 2),
  charity_fulfilled DECIMAL(18, 2),
  transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES Users (user_id)
);

CREATE TABLE IF NOT EXISTS denylist (
  jti uuid PRIMARY KEY,
  expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  token text PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP
);
