package exception

import "errors"

var (
	ErrorIdNotFound     = errors.New("id not found")
	ErrorEmailNotFound  = errors.New("email not found")
	ErrorEmailExist     = errors.New("email already exist")
	ErrorEventExist     = errors.New("event already exist")
	ErrorEventNotFound  = errors.New("event not found")
	ErrorValidation     = errors.New("validation failed")
	ErrorFailedLogin    = errors.New("email or password wrong")
	ErrorInvalidToken   = errors.New("invalid token refresh")
	ErrorStockNotEnough = errors.New("stock not enough")
	ErrorQty            = errors.New("Quantity cannot zero")
)
