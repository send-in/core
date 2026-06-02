DROP INDEX IF EXISTS idx_messages_schedule_time;
DROP INDEX IF EXISTS idx_messages_template_id;
DROP INDEX IF EXISTS idx_messages_account_id;

DROP INDEX IF EXISTS idx_templates_account_id;

DROP INDEX IF EXISTS idx_connections_public_id;
DROP INDEX IF EXISTS idx_connections_account_id;

DROP INDEX IF EXISTS idx_accounts_email;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS connections;
DROP TABLE IF EXISTS accounts;