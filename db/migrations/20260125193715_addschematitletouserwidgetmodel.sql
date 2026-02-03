-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_widgets  ADD COLUMN schema_title VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_widgets DROP COLUMN schema_title;
-- +goose StatementEnd
