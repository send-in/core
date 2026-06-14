CREATE TABLE account_connections (
    account_id UUID NOT NULL
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    connection_id UUID NOT NULL
        REFERENCES connections(id)
        ON DELETE CASCADE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (
        account_id,
        connection_id
    )
);

CREATE INDEX idx_account_connections_account
ON account_connections(account_id);

CREATE INDEX idx_account_connections_connection
ON account_connections(connection_id);

INSERT INTO account_connections (
    account_id,
    connection_id
)
SELECT
    account_id::uuid,
    id
FROM connections
WHERE account_id IS NOT NULL;

ALTER TABLE connections
DROP COLUMN account_id;