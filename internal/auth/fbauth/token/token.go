package token

import (
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	audience = "chevit-37fba"
	issuer   = "https://securetoken.google.com/chevit-37fba"
)

type Token struct {
	*jwt.Token
}

func (t *Token) FirebaseClaims() *Claims {
	return t.Token.Claims.(*Claims)
}

var (
	_ jwt.Claims       = (*Claims)(nil)
	_ json.Unmarshaler = (*Claims)(nil)
)

type Claims struct {
	auth.Token
}

// UnmarshalJSON
// TODO: refactor 두번 호출 하는 부분 개선 필요
func (c *Claims) UnmarshalJSON(bytes []byte) (err error) {
	err = json.Unmarshal(bytes, &c.Token)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &c.Token.Claims)
}

func (c *Claims) Valid() error {
	vErr := new(jwt.ValidationError)
	if c.Audience != audience {
		vErr.Inner = jwt.ErrTokenInvalidAudience
		vErr.Errors |= jwt.ValidationErrorAudience
	}

	if c.Issuer != issuer {
		vErr.Inner = jwt.ErrTokenInvalidIssuer
		vErr.Errors |= jwt.ValidationErrorIssuer
	}

	now := jwt.TimeFunc().Unix()

	if dt := c.Expires - now; dt <= 0 {
		vErr.Inner = fmt.Errorf("%s by %s", jwt.ErrTokenExpired, time.Duration(dt)*time.Second)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if dt := now - c.IssuedAt; dt < 0 {
		vErr.Inner = jwt.ErrTokenUsedBeforeIssued
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if dt := now - c.AuthTime; dt < 0 {
		vErr.Inner = jwt.ErrTokenNotValidYet
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

// CASE-1
/**
{
  "name": "이름",
  "picture": "https://lh3.googleusercontent.com/a/AAcHTtfvUqmT64pJ9PvAKpvAPvls7glr0Ocjp_qszBGrGrkQyA=s96-c",
  "iss": "https://securetoken.google.com/chevit-37fba",
  "aud": "chevit-37fba",
  "auth_time": 1689001774,
  "user_id": "B9GFoNzb6Th15QbOUajzMLBOFB73",
  "sub": "B9GFoNzb6Th15QbOUajzMLBOFB73",
  "iat": 1689250331,
  "exp": 1689253931,
  "email": "email@gmail.com",
  "email_verified": true,
  "firebase": {
    "identities": {
      "google.com": [
        "1234567890"
      ],
      "email": [
        "email@gmail.com"
      ]
    },
    "sign_in_provider": "google.com"
  }
}
*/

// CASE-2
/**
  {
    "iss": "https://securetoken.google.com/chevit-37fba",
    "aud": "chevit-37fba",
    "auth_time": 1688566888,
    "user_id": "Earpl406inMwYhFQGkzqWMS8S3E2",
    "sub": "Earpl406inMwYhFQGkzqWMS8S3E2",
    "iat": 1688566888,
    "exp": 1688570488,
    "email": "test@test.com",
    "email_verified": false,
    "firebase": {
      "identities": {
        "email": [
          "test@test.com"
        ]
      },
      "sign_in_provider": "password"
    }
  }
*/
