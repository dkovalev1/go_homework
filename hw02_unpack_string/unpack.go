package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var sb strings.Builder

	var currentSymbol rune   // Last symbol read from input.
	currentSymbolOk := false // If the last_symbol ready to copy to output.
	escape := false          // True if escape is active for the current symbol.

	for _, r := range input {
		switch {
		case escape && unicode.IsDigit(r):
			currentSymbol = r
			currentSymbolOk = true
			escape = false
		case unicode.IsDigit(r):
			if !currentSymbolOk {
				return "", ErrInvalidString
			}
			// Do repeat.
			for count := int(r - '0'); count > 0; count-- {
				sb.WriteRune(currentSymbol)
			}
			currentSymbolOk = false
		case escape && r == '\\':
			escape = false
			currentSymbol, currentSymbolOk = r, true
		case r == '\\':
			sb.WriteRune(currentSymbol)
			escape = true
			currentSymbolOk = false
		default:
			if escape {
				return "", ErrInvalidString
			}
			if currentSymbolOk {
				sb.WriteRune(currentSymbol)
			}
			currentSymbol, currentSymbolOk = r, true
			escape = false
		}
	}
	if currentSymbolOk {
		sb.WriteRune(currentSymbol)
	}
	return sb.String(), nil
}
