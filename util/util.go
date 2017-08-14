package util

import "testing"

func Assert(t *testing.T, expected interface{}, given interface{}) {
	if expected != given {
		t.Errorf("Expected %v but got %v", expected, given)
	}
}

func Contains(needle string, haystack []string) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}
