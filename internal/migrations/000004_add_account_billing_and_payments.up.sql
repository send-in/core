CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID NOT NULL
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    status TEXT NOT NULL DEFAULT 'pending',

    plan_credits INTEGER NOT NULL,
    amount BIGINT NOT NULL,

    provider TEXT NOT NULL,
    external_id TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE accounts
ADD COLUMN plan TEXT NOT NULL DEFAULT 'free',

ADD COLUMN plan_credits INTEGER NOT NULL DEFAULT 5,
ADD COLUMN credits_remaining INTEGER NOT NULL DEFAULT 5,

ADD COLUMN credits_renew_at TIMESTAMP NULL,

ADD COLUMN last_daily_reset_at TIMESTAMP NULL,

ADD COLUMN daily_syncs_used INTEGER NOT NULL DEFAULT 0,
ADD COLUMN daily_schedules_used INTEGER NOT NULL DEFAULT 0,

ADD COLUMN lifetime_syncs_used INTEGER NOT NULL DEFAULT 0,
ADD COLUMN lifetime_messages_used INTEGER NOT NULL DEFAULT 0;

CREATE INDEX idx_payments_account_id
ON payments(account_id);

CREATE INDEX idx_accounts_plan
ON accounts(plan);