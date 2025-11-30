-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT TIMEZONE('utc', NOW()),
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;