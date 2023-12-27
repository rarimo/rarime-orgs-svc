-- +migrate Up
ALTER TABLE organizations ADD COLUMN type VARCHAR(255);
ALTER TABLE organizations ADD COLUMN schema_url VARCHAR(255);

-- +migrate Down
ALTER TABLE organizations DROP COLUMN type;
ALTER TABLE organizations DROP COLUMN schema_url;
