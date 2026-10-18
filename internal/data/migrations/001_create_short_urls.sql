-- +goose Up
CREATE TABLE short_url_mappings (
    short_key TEXT PRIMARY KEY,
    target_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE short_url_mappings;
