CREATE TABLE IF NOT EXISTS guilds (
    id INT AUTO_INCREMENT PRIMARY KEY,
    discord_guild_id VARCHAR(20) NOT NULL UNIQUE,
    guild_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_guilds_discord_guild_id ON guilds(discord_guild_id); 