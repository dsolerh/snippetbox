package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	testCases := []struct {
		desc string
		tm   time.Time
		want string
	}{
		{
			desc: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			desc: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			desc: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*6)),
			want: "17 Dec 2020 at 09:00",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			hd := humanDate(tC.tm)
			if hd != tC.want {
				t.Errorf("want %q; got %q", tC.want, hd)
			}
		})
	}
}
