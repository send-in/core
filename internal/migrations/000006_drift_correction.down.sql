BEGIN;

ALTER TABLE templates
    DROP CONSTRAINT IF EXISTS fk_templates_account;

ALTER TABLE messages
    DROP CONSTRAINT IF EXISTS fk_messages_account;

ALTER TABLE connections
    DROP CONSTRAINT IF EXISTS fk_connections_account;

ALTER TABLE templates
    ALTER COLUMN account_id TYPE text
    USING account_id::text;

ALTER TABLE messages
    ALTER COLUMN account_id TYPE text
    USING account_id::text;

ALTER TABLE connections
    ALTER COLUMN account_id TYPE text
    USING account_id::text;

ALTER TABLE accounts
    ALTER COLUMN plan_credits SET DEFAULT 5,
    ALTER COLUMN credits_remaining SET DEFAULT 5;

COMMIT;