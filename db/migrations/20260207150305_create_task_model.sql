-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
	id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	task_type_id  uuid NOT NULL REFERENCES task_types(id) ON DELETE RESTRICT,
	name          text NOT NULL,
	description   text,
	created_at    timestamptz NOT NULL DEFAULT now(),
	due_at        timestamptz,
	updated_at    timestamptz,
	completed_at  timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
