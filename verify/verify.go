package verify

import "regexp"

func IsPhone(phone string) bool {
	return regexp.MustCompile(Phone).MatchString(phone)
}

func IsEmail(email string) bool {
	return regexp.MustCompile(Email).MatchString(email)
}
