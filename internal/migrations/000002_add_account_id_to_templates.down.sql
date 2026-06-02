DROP INDEX IF EXISTS idx_templates_account_id;

ALTER TABLE templates
DROP CONSTRAINT IF EXISTS fk_templates_account;

ALTER TABLE templates
DROP COLUMN IF EXISTS account_id;