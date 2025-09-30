-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS analytics
(
    id          SERIAL     PRIMARY KEY NOT NULL,
    short_url   VARCHAR(32) NOT NULL REFERENCES urls(short_url) ON DELETE CASCADE,
    user_agent  TEXT,
    device_type VARCHAR(32),
    os          VARCHAR(64),
    browser     VARCHAR(64),
    ip_address  INET,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS analytics;
-- +goose StatementEnd
