DROP INDEX IF EXISTS idx_connections_profile_urn;
DROP INDEX IF EXISTS idx_connections_recipient;

ALTER TABLE connections
DROP COLUMN profile_urn,
DROP COLUMN recipient;