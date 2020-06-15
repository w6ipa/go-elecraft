package utils

import (
	"strings"
	"testing"
)

func TestVerify(t *testing.T) {
	//                 1         2         3         4         5
	//       0123456789012345678901234567890123456789012345678901234567890
	line := "THE PROJECT GUTENBERG EBOOK OF METAMORPHOSIS,  BY FRANZ KAFKA"

	var testCases = []struct {
		input string
		oldX  int
		newX  int
	}{
		{
			input: "UTENBERG",
			oldX:  strings.Index(line, "GUTENBERG") + 1,
			newX:  strings.Index(line, " EB"),
		},
		{
			input: " XUT",
			oldX:  strings.Index(line, " GUTENBERG"),
			newX:  strings.Index(line, " GUTENBERG"),
		},
		{
			input: "X",
			oldX:  2,
			newX:  0,
		},
		{
			input: " T",
			oldX:  0,
			newX:  1,
		},
		{
			input: "X",
			oldX:  strings.Index(line, "C"),
			newX:  strings.Index(line, " PRO"),
		},
		{
			input: " BY",
			oldX:  strings.Index(line, "  BY"),
			newX:  strings.Index(line, " F"),
		},
	}

	for i, c := range testCases {
		dx := CheckAndAdvance([]byte(line), c.oldX, []byte(c.input))
		if c.oldX+dx != c.newX {
			t.Fatalf("%d: input %s - Expected %d got %d", i, c.input, c.newX, c.oldX+dx)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	var testCases = []struct {
		s      string
		prefix string
		x      int
		t      bool
	}{
		{
			s:      "  BY",
			prefix: " BY",
			x:      1,
			t:      true,
		},
		{
			s:      "    BY",
			prefix: " BY",
			x:      2,
			t:      true,
		},
		{
			s:      "   BY",
			prefix: " BT",
			x:      0,
			t:      false,
		},
		{
			s:      "BY",
			prefix: "BY",
			x:      0,
			t:      true,
		},
		{
			s:      "BY",
			prefix: "BX",
			x:      0,
			t:      false,
		},
	}

	for i, c := range testCases {
		b, x := hasPrefix([]byte(c.s), []byte(c.prefix), 0x20)
		if b != c.t && x != c.x {
			t.Fatalf("%d: s: %s - Expected %d/%t got %d/%t", i, c.s, c.x, c.t, x, b)
		}
	}
}
