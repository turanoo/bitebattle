CREATE TABLE poll_options (
    id TEXT PRIMARY KEY,
    poll_id TEXT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    restaurant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    image_url TEXT,
    menu_url TEXT,
    UNIQUE (poll_id, restaurant_id)
);