-- Down
ALTER TABLE messages
DROP COLUMN IF EXISTS profile_urn,
DROP COLUMN IF EXISTS recipient;