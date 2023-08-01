package user

import "d.kin-app/models/user"

func userToDTO(u *user.User) (res responseDTO) {
	res = responseDTO{
		UserId:      u.UserId,
		ImageURL:    u.ImageURL,
		Nickname:    nil,
		YearOfBirth: nil,
		Gender:      nil,
		NeedSignUp:  false,
	}
	if len(u.Nickname) > 0 {
		res.Nickname = &u.Nickname
	}

	if u.YearOfBirth > 0 {
		res.YearOfBirth = &u.YearOfBirth
	}

	if len(u.Gender) > 0 {
		res.Gender = &u.Gender
	}

	res.NeedSignUp = res.Nickname != nil &&
		res.YearOfBirth != nil &&
		res.Gender != nil
	return
}

type responseDTO struct {
	UserId      string       `json:"userId"`
	ImageURL    *string      `json:"imageURL"`
	Nickname    *string      `json:"nickname"`
	YearOfBirth *int16       `json:"yearOfBirth"`
	Gender      *user.Gender `json:"gender"`
	NeedSignUp  bool         `json:"needSignUp"`
}

type nicknameDTO struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=8"`
}

type yearOfBirthDTO struct {
	YearOfBirth int16 `json:"yearOfBirth" validate:"required,year_of_birth"`
}

type genderDTO struct {
	Gender user.Gender `json:"gender" validate:"required,gender"`
}
