package exception

import "errors"

var (
	ErrorIdNotFound    = errors.New("id not found")
	ErrorEmailNotFound = errors.New("email not found")
	ErrorEmailExist    = errors.New("email already exist")
	ErrorValidation    = errors.New("validation failed")
	ErrorFailedLogin    = errors.New("email or password wrong")
	ErrorInvalidToken    = errors.New("invalid token refresh")
)