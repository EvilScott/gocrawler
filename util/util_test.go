package util

import "testing"

func TestAssertEquals(t *testing.T) {
	AssertEquals(t, true, true, "AssertEquals")
}

func TestAssertEqualSlice(t *testing.T) {
	AssertEqualSlice(t, []string{"a", "b"}, []string{"a", "b"}, "AssertEqualSlice")
}
