package util

import (
    "testing"
    "reflect"
)

func AssertEquals(t *testing.T, expected, given interface{}, message string) {
    if expected != given {
        t.Errorf("%s :: Expected %v but got %v", message, expected, given)
    }
}

func AssertEqualSlice(t *testing.T, expected, given []string, message string) {
    if !reflect.DeepEqual(expected, given) {
        t.Errorf("%s :: Expected %v but got %v", message, expected, given)
    }
}
