package utils

import "bytes"

func CheckAndAdvance(line []byte, x int, input []byte) int {
	// if input match one or more char from current position - advance by the length of input
	// if there is a mistake, move back to begining of the word.

	if len(line) == 0 || len(input) == 0 || x > len(line) {
		return 0
	}
	trimedInput := input
	if x == 0 && input[0] == 0x20 {
		// Elecraft add a space to the start of every session
		trimedInput = input[1:]
	}

	if b, skip := hasPrefix(line[x:], trimedInput, 0x20); b == true {
		return len(trimedInput) + skip
	}
	// advance to the next word as long as what is before is correct
	for i, c := range input {
		if len(line) > x+i && c == line[x+i] && c == 0x20 {
			return i
		}
	}
	// walk back to begining of current word
	for i := x; i > 0; i-- {
		if line[i] == 0x20 {
			return i - x
		}
	}
	return -1 * x
}

// same as bytes.HasPrefix, but skip given characters if first character in input is skip character
// returns the number of skipped characters
func hasPrefix(s []byte, prefix []byte, skip byte) (bool, int) {
	if len(s) < len(prefix) {
		return false, 0
	}
	if prefix[0] != skip {
		return bytes.HasPrefix(s, prefix), 0
	}
	for i := 0; i < len(s); i++ {
		if s[i] != skip {
			return bytes.HasPrefix(s[i-1:], prefix), i - 1
		}
	}
	return false, 0
}
