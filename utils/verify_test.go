package utils

import (
	"strings"
	"testing"
)

func TestVerify(t *testing.T) {
	line := "THE PROJECT GUTENBERG EBOOK OF METAMORPHOSIS,  BY FRANZ KAFKA"
	x := strings.Index(line, "GUTENBERG") + 1
	dx := CheckAndAdvance([]byte(line), x, []byte("UTENBERG"))
	if dx != 8 {
		t.Fatalf("Expected 3 got %d", dx)
	}
	x = strings.Index(line, " GUTENBERG")
	dx = CheckAndAdvance([]byte(line), x, []byte(" XUT"))
	if dx != 0 {
		t.Fatalf("Expected 1 got %d", dx)
	}
	x = 2
	dx = CheckAndAdvance([]byte(line), x, []byte("X"))
	if dx+x != 0 {
		t.Fatalf("Expected 0 got %d", dx)
	}
	x = 0
	dx = CheckAndAdvance([]byte(line), x, []byte(" T"))
	if dx != 1 {
		t.Fatalf("Expected 1 got %d", dx)
	}
	x = 10
	dx = CheckAndAdvance([]byte(line), x, []byte("X"))
	if x+dx != 3 {
		t.Fatalf("Expected -7 got %d", dx)
	}

	x = strings.Index(line, "  BY")
	dx = CheckAndAdvance([]byte(line), x, []byte(" BY"))
	if dx != 4 {
		t.Fatalf("Expected 4 got %d", dx)
	}
}
