package global_varables

import "errors"

var (
	USERNAME_DUP  = errors.New("user name is duplicated.")
	EMAIL_DUP     = errors.New("email is duplicated.")
	NOT_FOUND_ROW = errors.New("with given data no row found.")
)

const (
	USER_CTX   = "USER"
	USER_ROLE  = "USER"
	ADMIN_ROLE = "ADMIN"
)
