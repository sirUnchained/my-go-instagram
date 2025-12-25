CREATE TABLE IF NOT EXISTS files (
    id bigserial PRIMARY KEY,
    filename VARCHAR(255),
    filepath TEXT,
    size_bytes INT,
    creator bigserial REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_files_creator ON files(creator);
