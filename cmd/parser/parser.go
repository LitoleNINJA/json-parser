package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/LitoleNINJA/json-parser/cmd/tokenizer"
)

type JsonObject map[string]interface{}
type JsonArray []interface{}

var pos int = 0

func Parse(tokens []tokenizer.Token) (interface{}, error) {
	// reset pos
	pos = 0

	parsedJson, err := parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("error while parsing : %v", err)
	}

	return parsedJson, nil
}

func parse(tokens []tokenizer.Token) (interface{}, error) {
	token, err := nextToken(tokens)
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case tokenizer.LeftBrace:
		return parseObject(tokens)
	case tokenizer.LeftBracket:
		return parseArray(tokens)
	case tokenizer.String:
		log.Printf("Parsed String : %s", token.Value)
		return token.Value, nil
	case tokenizer.Number:
		intVal, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing %s to int", token.Value)
		}

		log.Printf("Parsed Number : %d", intVal)
		return intVal, nil
	case tokenizer.Boolean:
		boolVal, err := strconv.ParseBool(token.Value)
		if err != nil {
			return nil, fmt.Errorf("error while parsing %s to bool", token.Value)
		}

		log.Printf("Parsed Bool : %t", boolVal)
		return boolVal, nil
	case tokenizer.Null:
		log.Printf("Parsed Null : %s", token.Value)
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Type)
	}
}

func parseObject(tokens []tokenizer.Token) (JsonObject, error) {
	jsonObject := make(JsonObject)

	for pos < len(tokens) {
		token, err := nextToken(tokens)
		if err != nil {
			return nil, err
		}

		// check for empty object
		if token.Type == tokenizer.RightBrace {
			log.Printf("Parsed Empty Object: %+v", jsonObject)
			return jsonObject, nil
		}

		// need a string
		if token.Type != tokenizer.String {
			return nil, fmt.Errorf("expected string as object key, got %s", token.Type)
		}

		key := token.Value

		// need a colon after key
		token, err = nextToken(tokens)
		if err != nil {
			return nil, err
		}

		if token.Type != tokenizer.Colon {
			return nil, fmt.Errorf("expected colon after key, got %s", token.Type)
		}

		// read value recursively
		value, err := parse(tokens)
		if err != nil {
			return nil, err
		}

		jsonObject[key] = value

		// need comma or right brace to continue
		token, err = nextToken(tokens)
		if err != nil {
			return nil, err
		}

		if token.Type == tokenizer.RightBrace {
			log.Printf("Parsed Object: %+v", jsonObject)
			return jsonObject, nil
		} else if token.Type == tokenizer.Comma {
			continue
		} else {
			return nil, fmt.Errorf("expected comma after object value, got %s", token.Type)
		}
	}

	return nil, fmt.Errorf("unexpected end of tokens while parsing object")
}

func parseArray(tokens []tokenizer.Token) (JsonArray, error) {
	jsonArray := make(JsonArray, 0)

	for pos < len(tokens) {
		token, err := nextToken(tokens)
		if err != nil {
			return nil, err
		}

		// check for empty array
		if token.Type == tokenizer.RightBracket {
			log.Printf("Parsed Empty Array : %v", jsonArray)
			return jsonArray, nil
		}

		// read value recursively
		pos--
		value, err := parse(tokens)
		if err != nil {
			return nil, err
		}

		jsonArray = append(jsonArray, value)

		// need comma or right brace to continue
		token, err = nextToken(tokens)
		if err != nil {
			return nil, err
		}

		if token.Type == tokenizer.RightBracket {
			log.Printf("Parsed Array : %v", jsonArray)
			return jsonArray, nil
		} else if token.Type == tokenizer.Comma {
			continue
		} else {
			return nil, fmt.Errorf("expected comma or ] after array value, got %s", token.Type)
		}
	}

	return nil, fmt.Errorf("unexpected end of tokens while parsing array")
}

func nextToken(tokens []tokenizer.Token) (tokenizer.Token, error) {
	if pos < len(tokens) {
		token := tokens[pos]
		pos++

		return token, nil
	} else {
		return tokenizer.Token{}, fmt.Errorf("end of tokens")
	}
}
