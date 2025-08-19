-- Create users table
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	phone_number VARCHAR(20) UNIQUE NOT NULL,
	scope VARCHAR(20) NOT NULL DEFAULT 'user',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_scope ON users(scope);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- Add constraints
DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_scope') THEN
		ALTER TABLE users ADD CONSTRAINT chk_users_scope CHECK (scope IN ('user', 'admin'));
	END IF;
END $$;

DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_phone_number') THEN
		ALTER TABLE users ADD CONSTRAINT chk_users_phone_number CHECK (phone_number ~ '^[+]?[1-9][0-9]{11,14}$');
	END IF;
END $$;