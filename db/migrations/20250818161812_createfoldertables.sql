-- +goose Up
-- +goose StatementBegin
-- folders table
CREATE TABLE folders (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    created_datetime TIMESTAMP NOT NULL DEFAULT now(),
    updated_datetime TIMESTAMP NOT NULL DEFAULT now()
);
-- notes table
CREATE TABLE notes (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    content VARCHAR NOT NULL,
    created_datetime TIMESTAMP NOT NULL DEFAULT now(),
    updated_datetime TIMESTAMP NOT NULL DEFAULT now()
);
-- bookmarks table
Create table bookmarks (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    link VARCHAR NOT NULL,
    icon VARCHAR,
    description VARCHAR,
    created_datetime TIMESTAMP NOT NULL DEFAULT now(),
    updated_datetime TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookmarks;
DROP TABLE notes;
DROP TABLE folders;
-- +goose StatementEnd
