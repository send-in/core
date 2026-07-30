ALTER TABLE connections
ADD COLUMN profile_urn TEXT,
ADD COLUMN recipient TEXT;

CREATE INDEX idx_connections_profile_urn
ON connections (profile_urn);

CREATE INDEX idx_connections_recipient
ON connections (recipient);