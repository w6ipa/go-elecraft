package utils

import "bytes"

func CheckAndAdvance(line []byte, x int, input []byte) int {
	// if input match one or more char from current position - advance by the length of input
	// if there is a mistake, move back to begining of the word.
	if x == 0 && input[0] == 0x20 {
		// maybe a trailing space at the start
		input = input[1:]
	}
	if bytes.HasPrefix(line[x:], input) {
		return len(input)
	}
	// advance to the next word as long as what is before is correct
	for i, c := range input {
		if len(line) >= x+i && c == line[x+i] && c == 0x20 {
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
