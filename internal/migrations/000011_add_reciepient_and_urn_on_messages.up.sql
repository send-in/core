-- Up
ALTER TABLE messages
ADD COLUMN profile_urn TEXT,
ADD COLUMN recipient TEXT;