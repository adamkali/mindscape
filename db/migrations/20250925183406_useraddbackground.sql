-- +goose Up
-- +goose StatementBegin
-- Add background column to user table as a filename
ALTER TABLE users ADD COLUMN background VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove background column from user table
ALTER TABLE users DROP COLUMN background;
-- +goose StatementEnd
