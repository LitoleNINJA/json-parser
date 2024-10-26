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

func ParseJSON(fileData []byte) (interface{}, error) {

	log.Print("Starting tokenizer ... ")
	tokens, err := tokenizer.TokenizeJSON(fileData)
	if err != nil {
		return nil, fmt.Errorf("error while tokenizing : %v", err)
	}

	// reset pos
	pos = 0

	log.Print("Starting parser ...")
	parsedJson, err := parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("error while parsing : %v", err)
	}

	if pos < len(tokens) {
		return nil, fmt.Errorf("unexpected tokens remaining after parsing")
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
		// number of form 0...
		if token.Value[0] == '0' && len(token.Value) > 1 && token.Value[1] != '.' && token.Value[1] != 'e' && token.Value[1] != 'E' {
			return nil, fmt.Errorf("invalid number %s", token.Value)
		}

		floarVal, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing %s to float", token.Value)
		}

		log.Printf("Parsed Number : %f", floarVal)
		return floarVal, nil
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
			// object not empty
			if len(jsonObject) != 0 {
				return jsonObject, fmt.Errorf("unexpected , at end of object")
			}
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
			return nil, fmt.Errorf("expected colon after key, got %s", token.Value)
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
			return nil, fmt.Errorf("expected comma after object value, got %s", token.Value)
		}
	}

	return nil, fmt.Errorf("unexpected end of tokens while parsing object")
}

func parseArray(tokens []tokenizer.Token) (JsonArray, error) {
	jsonArray := make(JsonArray, 0)
	log.Print("Array parse ...")
	for {
		token, err := nextToken(tokens)
		if err != nil {
			return nil, err
		}

		// check for empty array
		if token.Type == tokenizer.RightBracket {
			// array not empty
			if len(jsonArray) != 0 {
				return jsonArray, fmt.Errorf("unexpected , at end of array")
			}
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
			return nil, fmt.Errorf("expected comma or ] after array value, got %s", token.Value)
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
