package main

import "testing"

func TestParseConfigFile(t *testing.T) {
	expected := []string{"./repo1", "./repo2", "./repo3"}

	if r := parseConfigFile(); !stringSliceEq(expected, r) {
		t.Errorf("parseConfigFile() = %v, want %v", r, expected)
	}
}

// stringSliceEq compares two string slices and returns true for matching slices
func stringSliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, val := range a {
		if val != b[i] {
			return false
		}
	}

	return true
}
