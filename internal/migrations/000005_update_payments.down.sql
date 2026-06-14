DROP INDEX IF EXISTS idx_payments_account_id;
DROP INDEX IF EXISTS idx_payments_order_id;
DROP INDEX IF EXISTS idx_payments_external_id;

ALTER TABLE payments
    ALTER COLUMN external_id SET NOT NULL;

ALTER TABLE payments
    DROP COLUMN IF EXISTS completed_at,
    DROP COLUMN IF EXISTS payload,
    DROP COLUMN IF EXISTS order_id,
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS plan;