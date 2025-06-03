CREATE TABLE poll_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    restaurant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    image_url TEXT,
    menu_url TEXT,
    UNIQUE (poll_id, restaurant_id)
);