package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var sb strings.Builder

	// Assume that 0 is invalid rune value and will never happen in input string
	var current_symbol rune = 0        // last symbol read from input
	var current_symbol_ok bool = false // if the last_symbol ready to copy to output
	var escape bool = false            // true if escape is active for the current symbol

	for _, r := range input {
		switch {
		case escape && unicode.IsDigit(r):
			current_symbol = r
			current_symbol_ok = true
			escape = false
		case unicode.IsDigit(r):
			if !current_symbol_ok {
				return "", ErrInvalidString
			}
			// Do repeat
			for count := int(r - '0'); count > 0; count-- {
				sb.WriteRune(current_symbol)
			}
			current_symbol_ok = false
		case escape && r == '\\':
			escape = false
			current_symbol, current_symbol_ok = r, true
		case r == '\\':
			sb.WriteRune(current_symbol)
			escape = true
			current_symbol_ok = false
		default:
			if escape {
				return "", ErrInvalidString
			}
			if current_symbol_ok {
				sb.WriteRune(current_symbol)
			}
			current_symbol, current_symbol_ok = r, true
			escape = false
		}
	}
	if current_symbol_ok {
		sb.WriteRune(current_symbol)
	}
	return sb.String(), nil
}
