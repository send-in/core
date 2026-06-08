DROP INDEX IF EXISTS idx_accounts_plan;
DROP INDEX IF EXISTS idx_payments_account_id;

ALTER TABLE accounts
DROP COLUMN IF EXISTS lifetime_messages_used,
DROP COLUMN IF EXISTS lifetime_syncs_used,
DROP COLUMN IF EXISTS daily_schedules_used,
DROP COLUMN IF EXISTS daily_syncs_used,
DROP COLUMN IF EXISTS last_daily_reset_at,
DROP COLUMN IF EXISTS credits_renew_at,
DROP COLUMN IF EXISTS credits_remaining,
DROP COLUMN IF EXISTS plan_credits,
DROP COLUMN IF EXISTS plan;

DROP TABLE IF EXISTS payments;