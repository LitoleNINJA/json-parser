package tokenizer

import (
	"bytes"
	"fmt"
	"log"
	"unicode"
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

func readRune(reader *bytes.Reader) (rune, error) {
	ch, _, err := reader.ReadRune()
	column++

	return ch, err
}

func unreadRune(reader *bytes.Reader) {
	reader.UnreadRune()
	column--
}

func Tokenize(fileData []byte) ([]Token, error) {

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
				// tokens = append(tokens, Token{
				// 	Type:  EOF,
				// 	Value: "EOF",
				// })
				break
			}

			log.Printf("Error while reading rune at line : %d, column : %d!", lineNumber, column)
			return tokens, err
		}

		if unicode.IsSpace(ch) {
			if ch == '\n' {
				// log.Printf("Newline !")
				lineNumber++
				column = 0
			}
			// log.Printf("Space skipped !")
			continue
		}
		// log.Printf("char : %c", ch)

		token, err := makeToken(ch, reader)
		if err != nil {
			return tokens, fmt.Errorf("error at line : %d, column : %d => %v", lineNumber, column, err)
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
			return token, fmt.Errorf("error reading string : %v", err)
		}

		token.Type, token.Value = String, value
	default:
		if unicode.IsDigit(ch) || ch == '-' {
			unreadRune(reader)
			value, err := readNumber(reader)
			if err != nil {
				return token, fmt.Errorf("error reading number : %v", err)
			}

			token.Type, token.Value = Number, value
		} else if ch == 't' {
			err := readLiteral(reader, "rue")
			if err != nil {
				return token, fmt.Errorf("error reading boolean : %v", err)
			}

			log.Printf("Read Bool : true")
			token.Type, token.Value = Boolean, "true"
		} else if ch == 'f' {
			err := readLiteral(reader, "alse")
			if err != nil {
				return token, fmt.Errorf("error reading boolean : %v", err)
			}

			log.Printf("Read Bool : false")
			token.Type, token.Value = Boolean, "false"
		} else if ch == 'n' {
			err := readLiteral(reader, "ull")
			if err != nil {
				return token, fmt.Errorf("error reading boolean : %v", err)
			}

			log.Printf("Read Null : null")
			token.Type, token.Value = Null, "null"
		} else {
			return token, fmt.Errorf("unknown literal : %c", ch)
		}
	}

	return token, nil
}

func readString(reader *bytes.Reader) (string, error) {
	var value string

	for {
		ch, err := readRune(reader)
		if err != nil {
			return value, fmt.Errorf("error while reading rune : %v", err)
		}

		// log.Printf("char : %c", ch)
		if ch == '"' {
			log.Printf("Read string : %s", value)
			break
		}

		value += string(ch)
	}

	return value, nil
}

func readNumber(reader *bytes.Reader) (string, error) {
	var value string
	isNegative := ""

	for {
		ch, err := readRune(reader)
		if err != nil {
			return value, fmt.Errorf("error reading rune : %v", err)
		}

		// log.Printf("Char : %c", ch)
		if ch == ',' || unicode.IsSpace(ch) || ch == ']' {
			unreadRune(reader)
			log.Printf("Read number : %s", value)
			break
		} else if ch == '-' {
			if isNegative != "" {
				return value, fmt.Errorf("unexpected symbol -")
			}
			isNegative = "-"
		} else if !unicode.IsDigit(ch) {
			return value, fmt.Errorf("invalid number format")
		}

		value += string(ch)
	}

	return isNegative + value, nil
}

func readLiteral(reader *bytes.Reader, expected string) error {
	for _, e := range expected {
		ch, err := readRune(reader)
		if err != nil {
			return fmt.Errorf("error while reading rune : %v", err)
		}

		if ch != e {
			return fmt.Errorf("invalid literal : %c", ch)
		}
	}

	return nil
}
