ALTER TABLE users
  ADD COLUMN phone_number TEXT UNIQUE,
  ADD COLUMN profile_pic_url TEXT,
  ADD COLUMN bio TEXT,
  ADD COLUMN last_login_at TIMESTAMPTZ;
