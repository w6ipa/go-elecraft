package utils

import "bytes"

func CheckAndAdvance(line []byte, x int, input []byte) int {
	// if input match one or more char from current position - advance by the length of input
	// if there is a mistake, move back to begining of the word.

	if len(line) == 0 || len(input) == 0 {
		return 0
	}
	trimedInput := input
	if x == 0 && input[0] == 0x20 {
		// Elecraft add a space to the start of every session
		trimedInput = input[1:]
	}

	// BT skip ahead
	if bytes.Compare([]byte(" BT"), input) == 0 {
		return 1
	}
	if bytes.HasPrefix(line[x:], trimedInput) {
		return len(trimedInput)
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
