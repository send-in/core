ALTER TABLE connections
ADD COLUMN account_id UUID
REFERENCES accounts(id)
ON DELETE CASCADE;

UPDATE connections c
SET account_id = ac.account_id
FROM account_connections ac
WHERE ac.connection_id = c.id;

DROP TABLE account_connections;