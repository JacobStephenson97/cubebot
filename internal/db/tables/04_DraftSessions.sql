CREATE TABLE IF NOT EXISTS draft_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    format_id INTEGER,
    guild_id INTEGER,
    created_by_user_id INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'queue', -- queue, drafting, matches_in_progress, completed, cancelled
    team_size INTEGER NOT NULL DEFAULT 3, -- 3 for 3v3, 4 for 4v4, 5 for 5v5
    total_players INTEGER NOT NULL DEFAULT 6, -- 6 for 3v3, 8 for 4v4, 10 for 5v5
    external_draft_url VARCHAR(255), -- Link to the external draft
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (format_id) REFERENCES draft_formats(id),
    FOREIGN KEY (guild_id) REFERENCES guilds(id),
    FOREIGN KEY (created_by_user_id) REFERENCES users(id)
);

CREATE INDEX idx_draft_sessions_guild_id ON draft_sessions(guild_id);
CREATE INDEX idx_draft_sessions_status ON draft_sessions(status); 