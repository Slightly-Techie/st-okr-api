package helper

import (
	"strings"
)

func GenerateCompanyCode(name, id string) string {
	// Ensure both strings are long enough
	if len(name) < 3 || len(id) < 3 {
		return ""
	}

	// Get the first 3 letters of the first string
	firstPart := strings.ToUpper(name[:3])
	// Get the last 3 letters of the second string
	lastPart := strings.ToUpper(id[len(id)-3:])

	// Combine the two parts
	return firstPart + lastPart
}
