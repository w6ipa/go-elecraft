package utils

import (
	"bytes"
	"testing"
)

func TestFilter(t *testing.T) {
	text := []byte{0xeb, 0xbb, 0xbf, 0x6f, 0x77, 0x20, 0x2a, 0x0d, 0x0a, 0x2a}
	expected := []byte{0x6f, 0x77, 0x20, 0x0d, 0x0a}
	res := FilterCW(text)
	if bytes.Compare(res, expected) != 0 {
		t.Errorf("Expected %+v result %+v", expected, res)
	}
}
