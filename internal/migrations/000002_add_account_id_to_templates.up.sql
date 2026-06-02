ALTER TABLE templates
ADD COLUMN account_id UUID;

ALTER TABLE templates
ADD CONSTRAINT fk_templates_account
FOREIGN KEY (account_id)
REFERENCES accounts(id)
ON DELETE CASCADE;

CREATE INDEX idx_templates_account_id
ON templates(account_id);