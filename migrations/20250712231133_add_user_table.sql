-- +goose Up
-- +goose StatementBegin
CREATE TABLE users ( 
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    passwd VARCHAR(255) NOT NULL

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE 
-- +goose StatementEnd