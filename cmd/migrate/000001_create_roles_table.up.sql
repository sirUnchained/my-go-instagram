CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name varchar(64) UNIQUE NOT NULL,
    description varchar(512),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name, description) VALUES 
('USER', 'User can create posts, upload files and delete his/her own posts.'),
('ADMIN', 'Admin can create posts, upload files and delete others posts, see reports.'),
('BOSS', 'Boss or OMNIPOTENT can do anything.')
ON CONFLICT (name) DO NOTHING;