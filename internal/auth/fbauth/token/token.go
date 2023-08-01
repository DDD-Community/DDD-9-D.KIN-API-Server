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

func (c *Claims) UnmarshalJSON(bytes []byte) (err error) {
	var claims map[string]any
	err = json.Unmarshal(bytes, &claims)
	if err != nil {
		return
	}
	c.Set(claims)
	return
}

func (c *Claims) Set(claims map[string]any) {
	c.Token.Claims = claims
	c.Token.AuthTime, _ = c.Token.Claims["auth_time"].(int64)
	c.Token.Issuer, _ = c.Token.Claims["iss"].(string)
	c.Token.Audience, _ = c.Token.Claims["aud"].(string)
	c.Token.Expires, _ = c.Token.Claims["exp"].(int64)
	c.Token.IssuedAt, _ = c.Token.Claims["iat"].(int64)
	c.Token.Subject, _ = c.Token.Claims["sub"].(string)
	c.Token.UID, _ = c.Token.Claims["uid"].(string)
	fb, ok := c.Token.Claims["firebase"].(map[string]any)
	if !ok {
		return
	}

	c.Token.Firebase.SignInProvider, _ = fb["sign_in_provider"].(string)
	c.Token.Firebase.Tenant, _ = fb["tenant"].(string)
	c.Token.Firebase.Identities, _ = fb["identities"].(map[string]any)
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

func (c *Claims) GetName() (res string) {
	res, _ = c.Claims["name"].(string)
	return
}

func (c *Claims) GetPicture() (res string) {
	res, _ = c.Claims["picture"].(string)
	return
}

func (c *Claims) GetUserId() (res string) {
	res, _ = c.Claims["user_id"].(string)
	return
}

func (c *Claims) GetEmail() (res string) {
	res, _ = c.Claims["email"].(string)
	return
}

func (c *Claims) GetEmailVerified() (res bool) {
	res, _ = c.Claims["email_verified"].(bool)
	return
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
