package user

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewHTTPHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/getUser", GetUser)
	r.Post("/signUpUser", SignUpUser)
	r.Post("/validationNickname", ValidationNickname)
	r.Put("/updateUser", UpdateUser)
	r.Post("/getProfileUploadURL", GetProfileUploadURL)
	return r
}
