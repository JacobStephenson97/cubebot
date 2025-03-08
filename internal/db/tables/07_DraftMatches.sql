CREATE TABLE IF NOT EXISTS draft_matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    round_id INTEGER,
    participant1_id INTEGER,
    participant2_id INTEGER,
    winner_id INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, in_progress, completed
    match_score VARCHAR(10), -- e.g. "2-1", "2-0"
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT different_participants CHECK (participant1_id != participant2_id),
    FOREIGN KEY (round_id) REFERENCES draft_rounds(id) ON DELETE CASCADE,
    FOREIGN KEY (participant1_id) REFERENCES draft_participants(id) ON DELETE CASCADE,
    FOREIGN KEY (participant2_id) REFERENCES draft_participants(id) ON DELETE CASCADE,
    FOREIGN KEY (winner_id) REFERENCES draft_participants(id) ON DELETE CASCADE
);

CREATE INDEX idx_draft_matches_round_id ON draft_matches(round_id);
CREATE INDEX idx_draft_matches_participant1_id ON draft_matches(participant1_id);
CREATE INDEX idx_draft_matches_participant2_id ON draft_matches(participant2_id); 