CREATE TABLE IF NOT EXISTS follows (
    id              BIGSERIAL PRIMARY KEY,
    follower        BIGINT REFERENCES users(id) ON DELETE CASCADE,
    following       BIGINT REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_follower_following UNIQUE(follower, following)
);

CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower);
CREATE INDEX IF NOT EXISTS idx_follows_folloing ON follows(following);