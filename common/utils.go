package common

import "regexp"

func ValidAddress(addr string) bool {
	match, _ := regexp.MatchString("^0x[a-fA-F0-9]{40}$", addr)
	return match
}

func ValidHash(hash string) bool {
	match, _ := regexp.MatchString("^0x([A-Fa-f0-9]{64})$", hash)
	return match
}
