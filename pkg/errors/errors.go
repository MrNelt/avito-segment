package errors

import "fmt"

var (
	ErrSegmentNotFound = fmt.Errorf("segment doesn't exist")
	ErrUserNotFound    = fmt.Errorf("user doesn't exist")
)
