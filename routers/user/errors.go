package user

import (
	"d.kin-app/internal/httpx"
	"d.kin-app/models/user"
	"net/http"
)

var (
	apiErrNeedSignUpFirst       = &httpx.APIError{HTTPStatus: http.StatusConflict, Code: "ERR-100", Message: "need sign up first"}
	apiErrNicknameAlreadyExists = &httpx.APIError{HTTPStatus: http.StatusConflict, Code: "ERR-101", Message: user.ErrNicknameAlreadyExists.Error()}
)
