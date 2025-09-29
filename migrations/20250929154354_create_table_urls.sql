-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls(
    id          SERIAL     PRIMARY KEY NOT NULL,
    url         TEXT                   NOT NULL UNIQUE,
    short_url   TEXT                   NOT NULL UNIQUE,
    created_at  TIMESTAMP              NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
