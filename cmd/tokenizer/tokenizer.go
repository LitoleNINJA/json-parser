package tokenizer

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/LitoleNINJA/json-parser/cmd/customError"
)

type Token struct {
	Type  string
	Value string
}

const (
	LeftBrace    = "LEFT_BRACE"
	RightBrace   = "RIGHT_BRACE"
	LeftBracket  = "LEFT_BRACKET"
	RightBracket = "RIGHT_BRACKET"
	Colon        = "COLON"
	Comma        = "COMMA"
	String       = "STRING"
	Number       = "NUMBER"
	Boolean      = "BOOLEAN"
	Null         = "NULL"
	EOF          = "EOF"
)

var lineNumber int = 1
var column int = 0

func TokenizeJSON(fileData []byte) ([]Token, error) {

	// reset lineNumber and column
	lineNumber = 1
	column = 0

	var tokens []Token
	reader := bytes.NewReader(fileData)

	for {
		ch, err := readRune(reader)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("EOF !")
				break
			}

			log.Printf("Error while reading rune at line : %d, column : %d!", lineNumber, column)
			return tokens, err
		}

		// Only allow specific whitespace characters as per JSON spec (formfeed returns true in go unicode.IsSpace())
		// Space, tab, CR, LF are allowed. Form feed and other whitespace are not.
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			if ch == '\n' {
				lineNumber++
				column = 0
			}
			continue
		}

		// convert to token
		token, err := makeToken(ch, reader)
		if err != nil {
			return tokens, customError.NewError(fmt.Errorf("error at line : %d, column : %d => %v", lineNumber, column, err))
		}

		tokens = append(tokens, token)
	}

	log.Printf("Tokens : %+v", tokens)
	return tokens, nil
}

func makeToken(ch rune, reader *bytes.Reader) (Token, error) {
	var token Token

	switch ch {
	case '{':
		token.Type, token.Value = LeftBrace, "{"
	case '}':
		token.Type, token.Value = RightBrace, "}"
	case '[':
		token.Type, token.Value = LeftBracket, "["
	case ']':
		token.Type, token.Value = RightBracket, "]"
	case ':':
		token.Type, token.Value = Colon, ":"
	case ',':
		token.Type, token.Value = Comma, ","
	case '"':
		value, err := readString(reader)
		if err != nil {
			return token, customError.NewError(fmt.Errorf("error reading string : %v", err))
		}

		token.Type, token.Value = String, value
	default:
		if unicode.IsDigit(ch) || ch == '-' {
			unreadRune(reader)
			value, err := readNumber(reader)
			if err != nil {
				return token, customError.NewError(fmt.Errorf("error reading number : %v", err))
			}

			token.Type, token.Value = Number, value
		} else if ch == 't' {
			err := readLiteral(reader, "rue")
			if err != nil {
				return token, customError.NewError(fmt.Errorf("error reading boolean : %v", err))
			}

			log.Printf("Read Bool : true")
			token.Type, token.Value = Boolean, "true"
		} else if ch == 'f' {
			err := readLiteral(reader, "alse")
			if err != nil {
				return token, customError.NewError(fmt.Errorf("error reading boolean : %v", err))
			}

			log.Printf("Read Bool : false")
			token.Type, token.Value = Boolean, "false"
		} else if ch == 'n' {
			err := readLiteral(reader, "ull")
			if err != nil {
				return token, customError.NewError(fmt.Errorf("error reading boolean : %v", err))
			}

			log.Printf("Read Null : null")
			token.Type, token.Value = Null, "null"
		} else {
			return token, customError.NewError(fmt.Errorf("unknown literal : %c", ch))
		}
	}

	return token, nil
}

func readString(reader *bytes.Reader) (string, error) {
	var value strings.Builder

	for {
		ch, err := readRune(reader)
		if err != nil {
			return value.String(), customError.NewError(fmt.Errorf("error while reading rune : %v", err))
		}

		// check end of string
		if ch == '"' {
			break
		}

		// handle escaped sequence
		if ch == '\\' {
			escapedString, err := readEscapedSequence(reader)
			if err != nil {
				return "", customError.NewError(fmt.Errorf("error while reading escaped sequece : %v", err))
			}

			value.WriteString(escapedString)
			continue
		}

		// Check for unescaped control characters and line breaks
		if ch <= 0x1F {
			var description string
			switch ch {
			case 0x00:
				description = "null"
			case 0x0A:
				description = "newline"
			case 0x0D:
				description = "carriage return"
			case 0x09:
				description = "tab"
			case 0x0C:
				description = "form feed"
			case 0x08:
				description = "backspace"
			default:
				description = fmt.Sprintf("unknown control character %#U", ch)
			}
			return "", customError.NewError(fmt.Errorf("unescaped %s in string literal", description))
		}

		value.WriteRune(ch)
	}

	log.Printf("Read string : %s", value.String())
	return value.String(), nil
}

func readEscapedSequence(reader *bytes.Reader) (string, error) {
	ch, err := readRune(reader)
	if err != nil {
		return "", err
	}

	switch ch {
	case '"', '\\', '/':
		return string(ch), nil
	case 'b':
		return "\b", nil
	case 'f':
		return "\f", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 't':
		return "\t", nil
	case 'u':
		return readUnicodeEscaped(reader)
	default:
		return "", customError.NewError(fmt.Errorf("invalid escape sequence \\%c", ch))
	}
}

func readUnicodeEscaped(reader *bytes.Reader) (string, error) {
	var hexStr string
	// read the HEX part
	for i := 0; i < 4; i++ {
		ch, err := readRune(reader)
		if err != nil {
			return "", customError.NewError(fmt.Errorf("incomplete unicode escape sequence"))
		}

		// check if ch belongs to hex
		if !isHexDigit(ch) {
			return "", customError.NewError(fmt.Errorf("invalid hex digit in unicode escape: %c", ch))
		}

		hexStr += string(ch)
	}

	// convert hex to int
	codePoint, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return "", customError.NewError(fmt.Errorf("invalid unicode escape sequence: \\u%s", hexStr))
	}

	// Handle surrogate pairs
	if codePoint >= 0xD800 && codePoint <= 0xDBFF {
		// This is a high surrogate, must be followed by a low surrogate
		ch, err := readRune(reader)
		if err != nil {
			return "", customError.NewError(fmt.Errorf("incomplete surrogate pair"))
		}

		// Check for the following low surrogate
		if ch == '\\' {
			ch2, err := readRune(reader)
			if err != nil || ch2 != 'u' {
				return "", customError.NewError(fmt.Errorf("expected low surrogate pair"))
			}

			lowSurrogateHex := ""
			for i := 0; i < 4; i++ {
				ch, err := readRune(reader)
				if err != nil {
					return "", customError.NewError(fmt.Errorf("incomplete low surrogate unicode escape sequence"))
				}

				// check if ch belongs to hex
				if !isHexDigit(ch) {
					return "", customError.NewError(fmt.Errorf("invalid hex digit in low surrogate unicode escape: %c", ch))
				}

				lowSurrogateHex += string(ch)
			}

			// convert hex to int
			lowSurrogate, err := strconv.ParseUint(lowSurrogateHex, 16, 32)
			if err != nil {
				return "", customError.NewError(fmt.Errorf("invalid unicode escape sequence in low surrogate: \\u%s", lowSurrogateHex))
			}

			// Convert the surrogate pair to a single Unicode character
			if lowSurrogate < 0xDC00 || lowSurrogate > 0xDFFF {
				return "", customError.NewError(fmt.Errorf("invalid low surrogate %#U", lowSurrogate))
			}

			finalCode := ((codePoint - 0xD800) << 10) + (lowSurrogate - 0xDC00) + 0x10000
			return string(rune(finalCode)), nil
		}

		// If not followed by a low surrogate, it's an error in JSON
		return "", customError.NewError(fmt.Errorf("high surrogate not followed by low surrogate"))
	}

	// Handle lone low surrogates (invalid in JSON)
	if codePoint >= 0xDC00 && codePoint <= 0xDFFF {
		return "", customError.NewError(fmt.Errorf("lone low surrogate"))
	}

	return string(rune(codePoint)), nil
}

func isHexDigit(ch rune) bool {
	return unicode.IsDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func readNumber(reader *bytes.Reader) (string, error) {
	var value strings.Builder
	isNegative := false
	isFloat := false
	isExp := false
	expSignAllowed := false

	for {
		ch, err := readRune(reader)
		// Check end of number
		if ch == ',' || unicode.IsSpace(ch) || ch == ']' || ch == '}' || err != nil {
			unreadRune(reader)
			break
		}

		// handle - sign
		if ch == '-' {
			// check if - sign is allowed
			if value.Len() == 0 && !isNegative {
				isNegative = true
				value.WriteString("-")
			} else if isExp && expSignAllowed {
				value.WriteString("-")
				expSignAllowed = false
			} else {
				return "", customError.NewError(fmt.Errorf("unexpected symbol '-'"))
			}

			// - should be followed by digit
			nextCh, err := nextRune(reader)
			if err != nil || !unicode.IsDigit(nextCh) {
				return "", customError.NewError(fmt.Errorf("invalid number format : %s", value.String()+string(ch)))
			}
			continue
		}

		// handle .
		if ch == '.' {
			if isFloat || isExp {
				return "", customError.NewError(fmt.Errorf("unexpected symbol '.'"))
			}
			isFloat = true
			value.WriteString(".")

			// . should be followed by digit
			nextCh, err := nextRune(reader)
			if err != nil || !unicode.IsDigit(nextCh) {
				return "", customError.NewError(fmt.Errorf("invalid number format : %s", value.String()))
			}
			continue
		}

		// handle e / E
		if ch == 'e' || ch == 'E' {
			if isExp {
				return "", customError.NewError(fmt.Errorf("invalid use of 'e' in number"))
			}
			isExp = true
			expSignAllowed = true

			value.WriteRune(ch)
			continue
		}

		// handle + sign
		if ch == '+' {
			if isExp && expSignAllowed {
				value.WriteString("+")
				expSignAllowed = false
			} else {
				return "", customError.NewError(fmt.Errorf("unexpected symbol '+'"))
			}

			// + should be followed by digit
			nextCh, err := nextRune(reader)
			if err != nil || !unicode.IsDigit(nextCh) {
				return "", customError.NewError(fmt.Errorf("invalid number format : %s", value.String()+string(ch)))
			}
			continue
		}

		// Reset expSignAllowed once a digit follows 'e'/'E'
		if unicode.IsDigit(ch) {
			// reset expSign after digit
			if expSignAllowed {
				expSignAllowed = false
			}

			value.WriteRune(ch)
			continue
		} else {
			return "", customError.NewError(fmt.Errorf("unexpected symbol '%c'", ch))
		}
	}

	result := value.String()

	// number should not have leading 0
	if value.Len() > 1 && strings.HasPrefix(result, "0") && !isExp {
		return result, customError.NewError(fmt.Errorf("invalid number format : %s", result))
	}

	// - sign should not be followed by 0 prefixed number
	if value.Len() > 2 && strings.HasPrefix(result, "-0") && !isFloat && !isExp {
		return "", customError.NewError(fmt.Errorf("invalid number format %s", result))
	}

	log.Printf("Read number : %s", result)
	return result, nil
}

func readLiteral(reader *bytes.Reader, expected string) error {
	for _, e := range expected {
		ch, err := readRune(reader)
		if err != nil {
			return customError.NewError(fmt.Errorf("error while reading rune : %v", err))
		}

		if ch != e {
			return customError.NewError(fmt.Errorf("invalid literal : %c", ch))
		}
	}

	return nil
}

func readRune(reader *bytes.Reader) (rune, error) {
	ch, _, err := reader.ReadRune()
	column++

	return ch, err
}

func nextRune(reader *bytes.Reader) (rune, error) {
	ch, err := readRune(reader)
	if err != nil {
		return ch, err
	}
	unreadRune(reader)

	return ch, nil
}

func unreadRune(reader *bytes.Reader) {
	reader.UnreadRune()
	column--
}
