-- +goose Up
-- +goose StatementBegin
CREATE TABLE images ( 
    uuid varchar(36) DEFAULT(UUID()) PRIMARY KEY,
    name TEXT NOT NULL,
    alt TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE images;
-- +goose StatementEnd
