BEGIN;

ALTER TABLE accounts
    ALTER COLUMN plan_credits SET DEFAULT 0,
    ALTER COLUMN credits_remaining SET DEFAULT 0;

ALTER TABLE connections
    ALTER COLUMN account_id TYPE uuid
    USING account_id::uuid;

ALTER TABLE messages
    ALTER COLUMN account_id TYPE uuid
    USING account_id::uuid;

ALTER TABLE templates
    ALTER COLUMN account_id TYPE uuid
    USING account_id::uuid;

ALTER TABLE connections
    ADD CONSTRAINT fk_connections_account
    FOREIGN KEY (account_id)
    REFERENCES accounts(id)
    ON DELETE CASCADE;

ALTER TABLE messages
    ADD CONSTRAINT fk_messages_account
    FOREIGN KEY (account_id)
    REFERENCES accounts(id)
    ON DELETE CASCADE;

ALTER TABLE templates
    ADD CONSTRAINT fk_templates_account
    FOREIGN KEY (account_id)
    REFERENCES accounts(id)
    ON DELETE CASCADE;

COMMIT;