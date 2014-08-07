package mydumpster

import (
	"strings"
)

var scapeChars = map[string]string{
	string([]byte{0}): `\0`, // This is null string a.k.a (\0, ^@...)
}

// Check all the replacements needed
func ScapeCharacters(str string) string {

	for k, v := range scapeChars {
		str = strings.Replace(str, k, v, -1)
	}
	return str

}

// Checks and error and the program dies (panic)
func CheckKill(e error) {

	if e != nil {
		panic(e)
	}
}
