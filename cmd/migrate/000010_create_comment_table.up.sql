CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    post bigint REFERENCES posts(id) ON DELETE CASCADE not null,
    parent bigint REFERENCES comments(id) ON DELETE CASCADE,
    creator bigint REFERENCES users(id) ON DELETE CASCADE not null,
    content VARCHAR(2048),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent);
CREATE INDEX IF NOT EXISTS idx_comments_creator ON comments(creator);