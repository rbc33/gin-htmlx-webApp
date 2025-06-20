-- +goose Up
-- +goose StatementBegin
CREATE TABLE pages ( 
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    title TEXT NOT NULL,
    link TEXT NOT NULL

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pages
-- +goose StatementEnd
