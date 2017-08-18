package util

import (
    "testing"
    "reflect"
)

func AssertEquals(t *testing.T, expected, given interface{}) {
    if expected != given {
        t.Errorf("Expected %v but got %v", expected, given)
    }
}

func AssertEqualSlice(t *testing.T, expected, given []string) {
    if !reflect.DeepEqual(expected, given) {
        t.Errorf("Expected %v but got %v", expected, given)
    }
}
