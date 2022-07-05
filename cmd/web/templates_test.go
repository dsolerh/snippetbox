package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Initialize a new time.Time object and pass it to humanDate function.
	tm := time.Date(2022, 12, 17, 10, 0, 0, 0, time.UTC)
	hd := humanDate(tm)
	expected := "17 Dec 2022 at 10:00"

	if hd != expected {
		t.Errorf("want %q; got %q", expected, hd)
	}
}
