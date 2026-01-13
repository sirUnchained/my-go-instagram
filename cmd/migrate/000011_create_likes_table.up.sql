CREATE TABLE IF NOT EXISTS likes (
    id              bigserial PRIMARY KEY,
    post            BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    creator         BIGINT REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_post_like UNIQUE (post, creator)
);

CREATE INDEX IF NOT EXISTS idx_likes_user ON likes(creator);
CREATE INDEX IF NOT EXISTS idx_likes_post ON likes(post);
