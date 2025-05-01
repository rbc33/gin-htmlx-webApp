-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts ( 
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    title TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts
-- +goose StatementEnd
