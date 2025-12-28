CREATE TABLE IF NOT EXISTS profiles(
    id bigserial PRIMARY KEY,
    fullname VARCHAR(255),
    bio VARCHAR(512),
    avatar text,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);