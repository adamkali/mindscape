-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
	id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id            uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	refresh_token_hash text NOT NULL UNIQUE,
	user_agent         varchar(255),
	expires_at         timestamptz NOT NULL,
	created_at         timestamptz NOT NULL DEFAULT now(),
	last_used_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
