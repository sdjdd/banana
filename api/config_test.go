package main

import (
	"testing"
)

func TestParseSize(t *testing.T) {
	cases := []struct {
		str string
		num int64
	}{
		{"", 0},
		{"1024", 1024},
		{"2K", 1024 * 2},
		{"4MB", 1024 * 1024 * 4},
		{"1.5G", 1024 * 1024 * 1024 * 1.5},
	}
	for _, c := range cases {
		num, err := parseSize(c.str)
		if err != nil {
			t.Fatal(err)
		} else if num != c.num {
			t.Fatalf("want: %d, got: %d", c.num, num)
		}
	}

	if _, err := parseSize("illegal"); err == nil {
		t.Fatal("parse illegal size should returns an error")
	}
}
