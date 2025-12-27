CREATE TABLE IF NOT EXISTS profiles(
    id bigserial PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255),
    bio text,
    avatar text,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);