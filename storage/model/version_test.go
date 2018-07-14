package model

import "testing"

var versionToIntTests = []struct {
	version  string
	expected int
}{
	{"7.3.15", 70315},
	{"6.0.0", 60000},
}

func TestVersionToInt(t *testing.T) {
	for _, test := range versionToIntTests {
		result, err := VersionToInt(test.version)
		if err != nil {
			t.Fatalf("unexpected error: %s\n", err)
		}
		if result != test.expected {
			t.Fatalf("%s: expected: %d got %d\n", test.version, test.expected, result)
		}
	}
}
