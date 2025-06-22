-- +goose Up
-- +goose StatementBegin
CREATE TABLE pages ( 
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    title TEXT NOT NULL,
    link VARCHAR(255) NOT NULL UNIQUE

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pages
-- +goose StatementEnd
