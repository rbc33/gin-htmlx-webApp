-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_permalinks (
  id INTEGER UNIQUE PRIMARY KEY AUTO_INCREMENT,
  permalink VARCHAR(255) UNIQUE,
  post_id INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post_permalinks;
-- +goose StatementEnd
