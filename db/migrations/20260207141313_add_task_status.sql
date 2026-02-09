-- +goose Up
-- +goose StatementBegin
CREATE TABLE task_types (
	id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name      text UNIQUE NOT NULL,
	show_in_scheduled  boolean NOT NULL,
	show_in_completed  boolean NOT NULL,
	show_in_available  boolean NOT NULL,
	show_in_cancelled  boolean NOT NULL,
	created_at   timestamptz NOT NULL DEFAULT now()
);
INSERT INTO task_types (name, show_in_scheduled, show_in_completed, show_in_available, show_in_cancelled)
  VALUES
  ('AmbiguousTaskStatus', false, false, true, false),
  ('CancelledTaskStatus', false, false, true, true),
  ('DoneTaskStatus',       false, true,  false, false),
  ('HoldTaskStatus',       false, false, true, false),
  ('PendingTaskStatus',    true,  false, true, false),
  ('RecurringTaskStatus',  true,  false, true, false),
  ('UndoneTaskStatus',     true,  false, true, false),
  ('UrgentTaskStatus',     true,  false, true, false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_types;
-- +goose StatementEnd
