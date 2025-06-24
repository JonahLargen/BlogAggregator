-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at timestamp not null,
    name VARCHAR(100) NOT NULL unique
);

-- +goose Down
DROP TABLE users;