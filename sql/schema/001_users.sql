-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW())
);

-- +goose Down
DROP TABLE users;