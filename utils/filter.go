package utils

import "bytes"

func FilterCW(src []byte) []byte {
	return bytes.Map(cwFilter, src)

}
func cwFilter(r rune) rune {
	// CRLF
	if r == 0x0d || r == 0x0a {
		return r
	}
	// [a-z] -> uppercase
	if r > 0x60 && r < 0x7B {
		return r
	}
	// [ .,?]
	if r == 0x20 || r == 0x2E || r == 0x2c || r == 0x3f {
		return r
	}
	// [A-Z]
	if r > 0x40 && r < 0x5B {
		return r
	}
	// [0-9]
	if r > 0x2f && r < 0x3a {
		return r
	}
	return -1
}
