-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;