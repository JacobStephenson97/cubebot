CREATE TABLE IF NOT EXISTS draft_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    format_id INTEGER,
    guild_id VARCHAR(20),
    created_by_user_id VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'queue', -- queue, drafting, matches_in_progress, completed, cancelled
    external_draft_url VARCHAR(255), -- Link to the external draft
    channel_id VARCHAR(20), -- Channel ID of the draft session
    message_id VARCHAR(20), -- Message ID of the draft session
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (format_id) REFERENCES draft_formats(id) ON DELETE CASCADE,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_draft_sessions_guild_id ON draft_sessions(guild_id);
CREATE INDEX idx_draft_sessions_status ON draft_sessions(status); 