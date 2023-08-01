package validation

import "unicode/utf8"

//TODO: 정해지면 수정
//var invalidNicknamePattern = regexp.MustCompile("")

func IsValidNickname(nickname string) bool {
	cnt := utf8.RuneCountInString(nickname)
	return cnt >= 1 && cnt <= 8 //!invalidNicknamePattern.MatchString(nickname)
}
