CREATE TABLE IF NOT EXISTS files (
    id bigserial PRIMARY KEY,
    filename VARCHAR(255),
    filepath TEXT,
    size_bytes INT,
    creator_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_files_creator_id ON files(creator_id);
