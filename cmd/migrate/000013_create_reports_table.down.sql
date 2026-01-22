DROP INDEX IF EXISTS idx_reports_creator;
DROP INDEX IF EXISTS idx_reports_post;
DROP INDEX IF EXISTS idx_reports_parent_comment;
DROP INDEX IF EXISTS idx_reports_created_at;

DROP TABLE IF EXISTS reports;

DROP TYPE IF EXISTS report_reason;