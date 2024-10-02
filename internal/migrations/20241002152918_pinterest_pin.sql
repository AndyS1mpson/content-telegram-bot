-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pin (
    id BIGSERIAL PRIMARY KEY,
    image_url TEXT NOT NULL,
    status integer not null default 1,
    tg_channel TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS tg_channel_idx ON pin (tg_channel);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pin;
-- +goose StatementEnd