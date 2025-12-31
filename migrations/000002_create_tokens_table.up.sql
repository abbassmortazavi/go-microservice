CREATE TABLE tokens
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER     NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    hash_token text        NOT NULL,
    expired_at TIMESTAMP   NOT NULL,
    is_revoked BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Indexes for better performance
CREATE INDEX idx_tokens_user_id ON tokens (user_id);
CREATE INDEX idx_tokens_expired_at ON tokens (expired_at);
CREATE INDEX idx_tokens_hash_token ON tokens (hash_token);