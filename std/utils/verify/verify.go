package verify

import (
	"net/http"
	"regexp"
)

var (
	regexMobile = regexp.MustCompile("^1[0-9]{10}$")
	regexName   = regexp.MustCompile("^[a-zA-Z0-9]{5,8}$")
	regexPwd    = regexp.MustCompile("^[@$_a-zA-Z0-9]{5,8}$")
	regexUrl    = regexp.MustCompile("^[/\\-?=_.a-zA-Z0-9]{1,255}$")
	regexText   = regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9]{1,12}$")
)

func MobileValid(in string) bool {
	return regexMobile.MatchString(in)
}

func UsernameValid(in string) bool {
	return regexName.MatchString(in)
}

func PasswordValid(in string) bool {
	return regexPwd.MatchString(in)
}

func UrlValid(in string) bool {
	return regexUrl.MatchString(in)
}

func TextValid(in string) bool {
	return regexText.MatchString(in)
}

func MethodValid(in string) bool {
	return in == http.MethodGet ||
		in == http.MethodPut ||
		in == http.MethodPost ||
		in == http.MethodDelete
}
