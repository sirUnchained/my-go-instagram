CREATE TYPE report_reason AS ENUM (
    'porn content',
    'racist content',
    "spam report",
    'other'
);

CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    creator INTEGER REFERENCES users(id),
    post INTEGER NOT NULL REFERENCES posts(id),
    comment INTEGER REFERENCES comments(id),
    reason report_reason NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_report_target CHECK (
        post IS NOT NULL OR comment IS NOT NULL
    )
);

CREATE INDEX idx_reports_creator ON reports(creator);
CREATE INDEX idx_reports_post ON reports(post);
CREATE INDEX idx_reports_comment ON reports(comment);
CREATE INDEX idx_reports_created_at ON reports(created_at);