CREATE TABLE IF NOT EXISTS posts_tags (
    id          bigserial PRIMARY KEY,
    post        bigserial REFERENCES posts(id) ON DELETE CASCADE,
    tag         bigserial REFERENCES tags(id) ON DELETE CASCADE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_post_tag UNIQUE (post, tag)
);

CREATE INDEX IF NOT EXISTS idx_posts_tags_post ON posts_tags(post);
CREATE INDEX IF NOT EXISTS idx_posts_tags_tag ON posts_tags(tag);