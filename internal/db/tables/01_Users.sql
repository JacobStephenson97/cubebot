CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(20) NOT NULL UNIQUE PRIMARY KEY,
    discord_username VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    avatar_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for quickly looking up users by their Discord ID
CREATE INDEX idx_users_id ON users(id); 