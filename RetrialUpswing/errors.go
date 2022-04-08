package RetrialUpswing

import "errors"

var (
	ErrBadRequest = errors.New("Bad Request")
	ErrMaxRetryExceeds = errors.New("Max Retry count exceeded")
	ErrRetryWindowAbsent  = errors.New("Retry window absent for the strategy")
	ErrInvalidRetryStrategy = errors.New("Invalid Retry Strategy")
	ErrInvalidRetrialSession = errors.New("Invalid Retry Session")
)
