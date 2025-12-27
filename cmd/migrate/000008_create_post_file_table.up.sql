CREATE TABLE IF NOT EXISTS posts_files (
    id bigserial PRIMARY KEY,
    post bigserial REFERENCES posts(id) ON DELETE CASCADE,
    file bigserial REFERENCES files(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_post_file UNIQUE (post, file)
);

CREATE INDEX IF NOT EXISTS idx_posts_files_post ON posts_files(post);
CREATE INDEX IF NOT EXISTS idx_posts_files_file ON posts_files(file);