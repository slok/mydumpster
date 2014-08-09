package mydumpster

import (
	"strings"
)

var replacementChars = map[string]string{
	// Null char(\0, ^@...)
	string([]byte{0x00}): `\0`,

	// New line, LF (\n)
	string([]byte{0x0A}): `\n`,

	// Carriage return, CR (\r)
	string([]byte{0x0D}): `\r`,

	// Scape, ESC, ^[, (\)
	string([]byte{0x1A}): `\`,

	// Apostrophe, ', (\')
	string([]byte{0x27}): `\'`,
}

// Check all the replacements needed
func ReplaceCharacters(str string) string {

	for k, v := range replacementChars {
		str = strings.Replace(str, k, v, -1)
	}
	return str

}

func SearchStr(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// Checks and error and the program dies (panic)
func CheckKill(e error) {

	if e != nil {
		panic(e)
	}
}
