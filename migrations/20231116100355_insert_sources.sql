-- +goose Up
-- +goose StatementBegin
INSERT INTO sources (name, url)
    VALUES ('habr', 'https://habr.com/ru/rss/hub/go/all/?fl=ru'),
           ('godev', 'https://go.dev/blog/feed.atom'),
           ('hashnode', 'https://hashnode.com/n/golang/rss'),
           ('devto', 'https://dev.to/feed/tag/golang');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sources
    WHERE name IN ('habr', 'godev', 'hashnode', 'devto');
-- +goose StatementEnd
