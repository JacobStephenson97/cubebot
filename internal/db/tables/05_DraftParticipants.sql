CREATE TABLE IF NOT EXISTS draft_participants (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INTEGER,
    user_id INTEGER,
    team_number INTEGER, -- 1 or 2, assigned when draft starts
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, user_id),
    FOREIGN KEY (session_id) REFERENCES draft_sessions(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_draft_participants_session_id ON draft_participants(session_id);
CREATE INDEX idx_draft_participants_user_id ON draft_participants(user_id);
CREATE INDEX idx_draft_participants_team ON draft_participants(session_id, team_number); 