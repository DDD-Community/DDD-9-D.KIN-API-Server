package httpx

import "errors"

var (
	ErrContentTypeMustApplicationJSON = errors.New("Content-Type must application/json")
)
