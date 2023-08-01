package util

import "strings"

const (
	TokenTypeBearer = "Bearer "
)

func BearerToken(raw string) (token string, err error) {
	if !strings.HasPrefix(raw, TokenTypeBearer) {
		err = ErrNotBearerToken
		return
	}

	token = raw[len(TokenTypeBearer):]
	if strings.Count(token, ".") != 2 {
		token = ""
		err = ErrNotBearerToken
		return
	}
	return
}
