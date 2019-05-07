package request

import "regexp"

func VerifyPhone(phone string) bool {
	return regexp.MustCompile(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`).MatchString(phone)
}

func VerifyEmail(email string) bool {
	return regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`).MatchString(email)
}
