package util

import "testing"

func Assert(t *testing.T, expected interface{}, given interface{}) {
    if expected != given {
        t.Errorf("Expected %v but got %v", expected, given)
    }
}
