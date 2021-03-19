package helper

import "strings"

func ExtractData(cookie  string) []string {
	xs := strings.Split(cookie, " ")

	tokenValue := strings.Split(xs[1], "|")
	return  tokenValue
}
