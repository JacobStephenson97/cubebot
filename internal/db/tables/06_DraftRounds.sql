CREATE TABLE IF NOT EXISTS draft_rounds (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INTEGER,
    round_number INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, in_progress, completed
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, round_number),
    FOREIGN KEY (session_id) REFERENCES draft_sessions(id) ON DELETE CASCADE,
    CONSTRAINT round_range CHECK (round_number BETWEEN 1 AND 3)
);

CREATE INDEX idx_draft_rounds_session_id ON draft_rounds(session_id); 