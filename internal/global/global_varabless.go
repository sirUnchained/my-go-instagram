package global_varables

import "errors"

var (
	USERNAME_DUP  = errors.New("user name is duplicated")
	EMAIL_DUP     = errors.New("email is duplicated")
	DUP_ITEM      = errors.New("somthing gets duplicated")
	NOT_FOUND_ROW = errors.New("with given data no row found")
)

const (
	USER_CTX        = "USER"
	TARGET_USER_CTX = "TARGET_USER"
	USER_ROLE       = "USER"
	ADMIN_ROLE      = "ADMIN"
)

// WARNING if you want to change, this please make sure to change
// migrate file `000013_create_reports_table.up.sql` and in
// payloads `CreateReportPayload` too.
const (
	REPORT_SPAM           = "spam_report"
	REPORT_PORN_CONTENT   = "porn_content"
	REPORT_RACIST_CONTENT = "racist_content"
	REPORT_OTHER_CONTENT  = "other"
)
