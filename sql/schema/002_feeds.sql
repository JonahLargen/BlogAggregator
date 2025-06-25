-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(255) NOT NULL unique,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) 
    REFERENCES users(id) 
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;