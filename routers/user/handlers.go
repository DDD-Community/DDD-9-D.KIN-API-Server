package user

import (
	"d.kin-app/internal/httpx"
	"d.kin-app/internal/middleware"
	"d.kin-app/models/user"
	"d.kin-app/routers/util"
	"errors"
	"net/http"
	"time"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	c, ok := middleware.GetUserClaims(r)
	if !ok {
		httpx.ErrorUnauthorized(w)
		return
	}
	u, err := user.GetUser(r.Context(), c.Subject)
	if errors.Is(err, user.ErrUserNotFound) {
		nowMilli := time.Now().UnixMilli()
		u = &user.User{
			UserId:      c.Subject,
			CreatedTime: nowMilli,
			UpdatedTime: nowMilli,
			ActiveDevices: map[string]user.Device{
				httpx.DeviceId(r): {
					LastAccessTime: nowMilli,
				},
			},
		}
		picture := c.GetPicture()
		if len(picture) > 0 {
			u.ImageURL = &picture
		}
		_ = user.CreateUser(r.Context(), u)
	} else if err != nil {
		httpx.ErrorInternalServerError(w)
		//TODO: Logging
		return
	}

	err = httpx.SetBodyJSON(w, http.StatusOK, userToDTO(u))
	if err != nil {
		//TODO: Logging
	}
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {
	c, ok := middleware.GetUserClaims(r)
	if !ok {
		httpx.ErrorUnauthorized(w)
		return
	}

	body, err := util.GetValidBodyJSON[struct {
		nicknameDTO
		yearOfBirthDTO
		genderDTO
	}](r)
	if err != nil {
		httpx.ErrorBadRequest(w, err)
		return
	}

	u, err := user.GetUser(r.Context(), c.Subject)
	switch {
	case err == nil:
		err = u.DoSignUp(body.Nickname, body.YearOfBirth, body.Gender)
	case errors.Is(err, user.ErrUserNotFound):
		nowMilli := time.Now().UnixMilli()
		u = &user.User{
			UserId:      c.Subject,
			Nickname:    body.Nickname,
			YearOfBirth: body.YearOfBirth,
			Gender:      body.Gender,
			CreatedTime: nowMilli,
			UpdatedTime: nowMilli,
			ActiveDevices: map[string]user.Device{
				httpx.DeviceId(r): {
					LastAccessTime: nowMilli,
				},
			},
		}
		picture := c.GetPicture()
		if len(picture) > 0 {
			u.ImageURL = &picture
		}
		err = user.CreateUser(r.Context(), u)
	default:
		httpx.ErrorInternalServerError(w)
		//TODO: Logging
		return
	}

	switch {
	case err == nil:
		err = httpx.SetBodyJSON(w, http.StatusOK, userToDTO(u))
	case errors.Is(err, user.ErrInvalidNickname):
		httpx.ErrorBadRequest(w, err)
	case errors.Is(err, user.ErrNicknameAlreadyExists):
		err = apiErrNicknameAlreadyExists.ResponseWrite(w)
	default:
		httpx.ErrorInternalServerError(w)
		//TODO: Logging
		return
	}

	if err != nil {
		//TODO: Logging
	}
}

func ValidationNickname(w http.ResponseWriter, r *http.Request) {
	body, err := util.GetValidBodyJSON[nicknameDTO](r)
	if err != nil {
		httpx.ErrorBadRequest(w, err)
		return
	}

	err = user.IsUsableNickname(r.Context(), body.Nickname)
	switch {
	case err == nil:
		httpx.WriteStatus(w, http.StatusNoContent)
	case errors.Is(err, user.ErrInvalidNickname):
		httpx.ErrorBadRequest(w, err)
	case errors.Is(err, user.ErrNicknameAlreadyExists):
		err = apiErrNicknameAlreadyExists.ResponseWrite(w)
	}

	if err != nil {
		//TODO: Logging
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	c, ok := middleware.GetUserClaims(r)
	if !ok {
		httpx.ErrorUnauthorized(w)
		return
	}

	body, err := util.GetValidBodyJSON[struct {
		ImageURL *string `json:"imageURL" validate:"omitempty,startswith=https://"`
		nicknameDTO
	}](r)
	if err != nil {
		httpx.ErrorBadRequest(w, err)
		return
	}

	u, _ := user.GetUser(r.Context(), c.Subject)
	if u == nil {
		err = apiErrNeedSignUpFirst.ResponseWrite(w)
		if err != nil {
			//TODO: Logging
		}
		return
	}

	err = u.ProfileUpdate(body.ImageURL, body.Nickname)
	switch {
	case err == nil:
		err = httpx.SetBodyJSON(w, http.StatusOK, userToDTO(u))
	case errors.Is(err, user.ErrInvalidNickname):
		httpx.ErrorBadRequest(w, err)
	case errors.Is(err, user.ErrNicknameAlreadyExists):
		err = apiErrNicknameAlreadyExists.ResponseWrite(w)
	default:
		httpx.ErrorInternalServerError(w)
		//TODO: Logging
		return
	}

	if err != nil {
		//TODO: Logging
	}
}

func GetProfileUploadURL(w http.ResponseWriter, r *http.Request) {
	c, ok := middleware.GetUserClaims(r)
	if !ok {
		httpx.ErrorUnauthorized(w)
		return
	}

	body, err := util.GetValidBodyJSON[struct {
		FileSize int64  `json:"fileSize" validate:"required,max=104857600"`
		MimeType string `json:"mimeType" validate:"required,startswith=image/"`
	}](r)
	if err != nil {
		httpx.ErrorBadRequest(w, err)
		return
	}

	u, _ := user.GetUser(r.Context(), c.Subject)
	if u == nil {
		err = apiErrNeedSignUpFirst.ResponseWrite(w)
		if err != nil {
			//TODO: Logging
		}
		return
	}

	res, err := u.ProfileImageUploadURL(user.ImageFile{
		Size:     body.FileSize,
		MimeType: body.MimeType,
	})
	if err != nil {
		httpx.ErrorInternalServerError(w)
		//TODO: Logging
		return
	}
	if res.S3UploadURL == nil || res.S3UploadMethod == nil {
		httpx.ErrorInternalServerError(w)
		return
	}

	var resp struct {
		UploadURL    string `json:"uploadURL"`
		UploadMethod string `json:"uploadMethod"`
		ImageURL     string `json:"imageURL"`
	}

	resp.UploadURL = *res.S3UploadURL
	resp.UploadMethod = *res.S3UploadMethod
	resp.ImageURL = res.ImageURL()

	err = httpx.SetBodyJSON(w, http.StatusOK, resp)
	if err != nil {
		//TODO: Logging
	}
}
