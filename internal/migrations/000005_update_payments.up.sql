ALTER TABLE payments
    ADD COLUMN plan TEXT,
    ADD COLUMN currency TEXT,
    ADD COLUMN order_id TEXT,
    ADD COLUMN payload JSONB,
    ADD COLUMN completed_at TIMESTAMP;

ALTER TABLE payments
    ALTER COLUMN external_id DROP NOT NULL;

UPDATE payments
SET
    plan = 'pro',
    currency = 'USD'
WHERE plan IS NULL;

ALTER TABLE payments
    ALTER COLUMN plan SET NOT NULL;

ALTER TABLE payments
    ALTER COLUMN currency SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_external_id
    ON payments (external_id)
    WHERE external_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_order_id
    ON payments (order_id)
    WHERE order_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_payments_account_id
    ON payments (account_id);