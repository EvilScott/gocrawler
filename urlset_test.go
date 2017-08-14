package main

import (
	"os"
	"testing"
)

func assert(t *testing.T, expected interface{}, given interface{}) {
	if expected != given {
		t.Errorf("Expected %v but got %v", expected, given)
	}
}

func TestURLSet_AddURL(t *testing.T) {
	urlSet := NewURLSet()
	assert(t, false, urlSet.AddURL("foo"))
	assert(t, true, urlSet.AddURL("foo"))
	assert(t, false, urlSet.AddURL("bar"))
	assert(t, 2, urlSet.set["foo"])
	assert(t, 1, urlSet.set["bar"])
}

func TestURLSet_AddURLs(t *testing.T) {
	urlSet := NewURLSet()
	urlSet.AddURLs([]string{"foo", "foo", "bar"})
	assert(t, 2, urlSet.set["foo"])
	assert(t, 1, urlSet.set["bar"])
}

func TestURLSet_Contains(t *testing.T) {
	urlSet := NewURLSet()
	assert(t, false, urlSet.Contains("foo"))
	urlSet.AddURL("foo")
	assert(t, true, urlSet.Contains("foo"))
}

func TestURLSet_Length(t *testing.T) {
	urlSet := NewURLSet()
	urlSet.AddURL("foo")
	assert(t,1, urlSet.Length())
	urlSet.AddURL("foo")
	urlSet.AddURL("bar")
	assert(t,2, urlSet.Length())
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
