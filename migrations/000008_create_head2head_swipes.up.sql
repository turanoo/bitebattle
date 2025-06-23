CREATE TABLE head2head_swipes (
    id TEXT PRIMARY KEY,
    match_id TEXT NOT NULL REFERENCES head2head_matches(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    restaurant_id TEXT NOT NULL,
    restaurant_name TEXT NOT NULL,
    liked BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
