package error_messages

import "errors"

var (
	USERNAME_DUP  = errors.New("user name is duplicated.")
	EMAIL_DUP     = errors.New("email is duplicated.")
	NOT_FOUND_ROW = errors.New("with given data no row found.")
)
