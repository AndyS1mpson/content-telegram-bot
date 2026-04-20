-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pin (
    id BIGINT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'pin',
    status INTEGER NOT NULL DEFAULT 1,
    channel TEXT NOT NULL,
    query TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, channel)
);
CREATE INDEX IF NOT EXISTS pin_status_channel_query_idx ON pin (channel, query, status);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pin;
-- +goose StatementEnd
