-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content (
    id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    type TEXT NOT NULL,
    status integer not null default 1,
    channel TEXT NOT NULL,
    created_at timestamptz not null default now()
);
CREATE INDEX IF NOT EXISTS tg_channel_idx ON content (channel);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content;
-- +goose StatementEnd