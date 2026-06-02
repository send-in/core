CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,

    profile TEXT,
    picture TEXT,
    timezone TEXT,

    token TEXT,
    user_agent TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,

    public_id TEXT UNIQUE NOT NULL,

    first_name TEXT,
    last_name TEXT,

    bio TEXT,
    picture TEXT,

    company TEXT,
    country TEXT,
    timezone TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,

    label TEXT NOT NULL,
    value TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
    template_id UUID REFERENCES templates(id) ON DELETE SET NULL,

    name TEXT NOT NULL,
    picture TEXT,

    profile TEXT NOT NULL,
    company TEXT,
    timezone TEXT,

    message TEXT,

    schedule_time TIMESTAMP,

    is_sent BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_email
ON accounts(email);

CREATE INDEX idx_connections_account_id
ON connections(account_id);

CREATE INDEX idx_connections_public_id
ON connections(public_id);

CREATE INDEX idx_templates_account_id
ON templates(account_id);

CREATE INDEX idx_messages_account_id
ON messages(account_id);

CREATE INDEX idx_messages_template_id
ON messages(template_id);

CREATE INDEX idx_messages_schedule_time
ON messages(schedule_time);