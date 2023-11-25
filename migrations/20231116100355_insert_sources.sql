-- +goose Up
-- +goose StatementBegin
INSERT INTO sources (name, url)
VALUES ('habr-go', 'https://habr.com/ru/rss/hub/go/all/?fl=ru'),
       ('habr-microservices', 'https://habr.com/ru/rss/hubs/microservices/all/?fl=ru'),
       ('godev', 'https://go.dev/blog/feed.atom'),
       ('hashnode-golang', 'https://hashnode.com/n/golang/rss'),
       ('devto-golang', 'https://dev.to/feed/tag/golang');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM sources
WHERE name IN ('habr-go', 'habr-microservices','godev', 'hashnode-golang', 'devto-golang');
-- +goose StatementEnd
