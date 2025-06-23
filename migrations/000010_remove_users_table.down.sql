-- Recreate users table (if needed for rollback)
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20),
    profile_pic_url TEXT,
    bio TEXT,
    last_login_at TIMESTAMP
);
