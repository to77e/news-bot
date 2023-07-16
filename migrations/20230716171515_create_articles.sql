-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles
(
    id           SERIAL PRIMARY KEY,
    source_id    INT       NOT NULL,
    title        TEXT      NOT NULL,
    link         TEXT      NOT NULL,
    summary      TEXT      NOT NULL,
    published_at TIMESTAMP NOT NULL,
    created_at   TIMESTAMP NOT NULL,
    posted_at    TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd
