package awsx

import (
	"errors"
	"github.com/aws/smithy-go"
)

func CompareErrorCodes(err error, errorCodes ...string) bool {
	var ae smithy.APIError
	if !errors.As(err, &ae) {
		return false
	}

	code := ae.ErrorCode()
	for _, c := range errorCodes {
		if code == c {
			return true
		}
	}

	return false
}
