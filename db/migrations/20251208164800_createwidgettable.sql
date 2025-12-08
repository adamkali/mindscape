-- +goose Up
-- +goose StatementBegin
-- user_widgets table stores user widget instances with their configurations
CREATE TABLE user_widgets (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    schema_id UUID NOT NULL,

    -- Widget configuration stored as JSONB for flexibility
    -- This stores the user's custom values for the widget properties
    config JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Grid layout positioning
    position_x INT NOT NULL DEFAULT 0,
    position_y INT NOT NULL DEFAULT 0,
    width INT NOT NULL DEFAULT 1,
    height INT NOT NULL DEFAULT 1,

    -- Ordering for widgets at the same position
    z_index INT NOT NULL DEFAULT 0,

    -- Visibility toggle
    is_visible BOOLEAN NOT NULL DEFAULT true,

    -- Timestamps
    created_datetime TIMESTAMP NOT NULL DEFAULT now(),
    updated_datetime TIMESTAMP NOT NULL DEFAULT now(),

    -- Ensure user can't have duplicate widgets at the same position
    CONSTRAINT unique_user_widget_position UNIQUE(user_id, position_x, position_y),

    -- Validate positive dimensions
    CONSTRAINT positive_dimensions CHECK (width > 0 AND height > 0),
    CONSTRAINT positive_position CHECK (position_x >= 0 AND position_y >= 0)
);

-- Index for faster user widget lookups
CREATE INDEX idx_user_widgets_user_id ON user_widgets(user_id);

-- Index for schema lookups (to find all instances of a widget type)
CREATE INDEX idx_user_widgets_schema_id ON user_widgets(schema_id);

-- GIN index for JSONB config queries
CREATE INDEX idx_user_widgets_config ON user_widgets USING GIN(config);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_widgets_config;
DROP INDEX IF EXISTS idx_user_widgets_schema_id;
DROP INDEX IF EXISTS idx_user_widgets_user_id;
DROP TABLE user_widgets;
-- +goose StatementEnd
