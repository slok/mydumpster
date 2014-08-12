package mydumpster

import (
	"fmt"
	"testing"
)

type TestPair struct {
	censorship  Censorship
	value       string
	censoredStr string
	censoredNil bool
}

var tests = []TestPair{
	// Basic tests
	TestPair{Censorship{"", "", "", false, false, ""}, "Go rules", "Go rules", false},
	TestPair{Censorship{"", "", "_", false, false, ""}, "Go rules", "_Go rules", false},
	TestPair{Censorship{"", "_", "", false, false, ""}, "Go rules", "Go rules_", false},
	TestPair{Censorship{"", "", "", true, false, ""}, "Go rules", "", false},
	TestPair{Censorship{"", "", "", false, true, ""}, "Go rules", "", true},
	TestPair{Censorship{"", "", "", false, false, "Go is awesome"}, "Go rules", "Go is awesome", false},

	// Complex tests
	TestPair{Censorship{"", "_", "***", false, false, ""}, "Go rules", "***Go rules_", false},
	TestPair{Censorship{"", "__", "__", true, false, ""}, "Go rules", "", false},
	TestPair{Censorship{"", "__", "__", true, true, ""}, "Go rules", "", true},
	TestPair{Censorship{"", "__", "__", true, true, "Go is awesome"}, "Go rules", "Go is awesome", false},
}

// Global for all the tests ------------
// Execute at the beggining of all the tests
func setup(t *testing.T) {

}

// Execute at the end of all the tests
func tearDown(t *testing.T) {

}

// Unit tests --------------------------
func TestCensore(t *testing.T) {
	setup(t)
	defer tearDown(t)

	for _, pair := range tests {
		v, n := pair.censorship.censore(pair.value)
		if v != pair.censoredStr {
			t.Error(
				"For", fmt.Sprintf("'%v'", pair.censorship),
				"expected", fmt.Sprintf("'%s'", pair.censoredStr),
				"got", fmt.Sprintf("'%s'", v),
			)
		}

		if n != pair.censoredNil {
			t.Error(
				"For", fmt.Sprintf("'%v'", pair.censorship),
				"expected", fmt.Sprintf("'%t'", pair.censoredNil),
				"got", fmt.Sprintf("'%t'", n),
			)
		}
	}
}
