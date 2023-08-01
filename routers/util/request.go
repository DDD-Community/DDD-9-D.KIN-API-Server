package util

import (
	"d.kin-app/internal/httpx"
	"d.kin-app/internal/validation"
	"net/http"
)

func GetValidBodyJSON[T any](r *http.Request) (_ T, err error) {
	res, err := httpx.GetBodyJSON[T](r)
	if err != nil {
		return
	}

	err = validation.Valid(&res)
	if err != nil {
		return
	}

	return res, nil
}
