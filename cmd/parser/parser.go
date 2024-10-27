package parser

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/LitoleNINJA/json-parser/cmd/customError"
	"github.com/LitoleNINJA/json-parser/cmd/tokenizer"
)

var pos int = 0

func ParseJSON(fileData []byte, result any) error {

	// Get reflect.Value of the result pointer
	resultValue := reflect.ValueOf(result)
	if result == nil || resultValue.Kind() != reflect.Ptr {
		return customError.NewError(fmt.Errorf("result must be a pointer"))
	}

	log.Print("Starting tokenizer ... ")
	tokens, err := tokenizer.TokenizeJSON(fileData)
	if err != nil {
		return customError.NewError(fmt.Errorf("error while tokenizing : %v", err))
	}

	// reset pos
	pos = 0

	log.Print("Starting parser ...")
	parsedJson, err := parse(tokens)
	if err != nil {
		return customError.NewError(fmt.Errorf("error while parsing : %v", err))
	}

	if pos < len(tokens) {
		return customError.NewError(fmt.Errorf("unexpected tokens remaining after parsing"))
	}

	// Convert and assign the parsed value to the result
	if err := assignParsedValue(resultValue.Elem(), parsedJson); err != nil {
		return customError.NewError(fmt.Errorf("error assigning parsed value: %v", err))
	}

	return nil
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

func parseObject(tokens []tokenizer.Token) (map[string]interface{}, error) {
	jsonData := make(map[string]interface{})

	for pos < len(tokens) {
		token, err := nextToken(tokens)
		if err != nil {
			return nil, err
		}

		// check for empty object
		if token.Type == tokenizer.RightBrace {
			// object not empty
			if len(jsonData) != 0 {
				return jsonData, fmt.Errorf("unexpected , at end of object")
			}
			log.Printf("Parsed Empty Object: %+v", jsonData)
			return jsonData, nil
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

		jsonData[key] = value

		// need comma or right brace to continue
		token, err = nextToken(tokens)
		if err != nil {
			return nil, err
		}

		if token.Type == tokenizer.RightBrace {
			log.Printf("Parsed Object: %+v", jsonData)
			return jsonData, nil
		} else if token.Type == tokenizer.Comma {
			continue
		} else {
			return nil, fmt.Errorf("expected comma after object value, got %s", token.Value)
		}
	}

	return nil, fmt.Errorf("unexpected end of tokens while parsing object")
}

func parseArray(tokens []tokenizer.Token) ([]interface{}, error) {
	jsonData := make([]interface{}, 0)
	log.Print("Array parse ...")
	for {
		token, err := nextToken(tokens)
		if err != nil {
			return nil, err
		}

		// check for empty array
		if token.Type == tokenizer.RightBracket {
			// array not empty
			if len(jsonData) != 0 {
				return jsonData, fmt.Errorf("unexpected , at end of array")
			}
			log.Printf("Parsed Empty Array : %v", jsonData)
			return jsonData, nil
		}

		// read value recursively
		pos--
		value, err := parse(tokens)
		if err != nil {
			return nil, err
		}

		jsonData = append(jsonData, value)

		// need comma or right brace to continue
		token, err = nextToken(tokens)
		if err != nil {
			return nil, err
		}

		if token.Type == tokenizer.RightBracket {
			log.Printf("Parsed Array : %v", jsonData)
			return jsonData, nil
		} else if token.Type == tokenizer.Comma {
			continue
		} else {
			return nil, fmt.Errorf("expected comma or ] after array value, got %s", token.Value)
		}
	}
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

// Helper function to assign parsed value to the result with type checking
func assignParsedValue(result reflect.Value, parsedJson interface{}) error {
	// Handle null values
	if parsedJson == nil {
		result.Set(reflect.Zero(result.Type()))
		return nil
	}

	parsedType := reflect.TypeOf(parsedJson)
	resultType := result.Type()

	switch {
	case parsedType.AssignableTo(resultType):
		// Direct assignment if types are compatible
		result.Set(reflect.ValueOf(parsedJson))
		return nil

	case resultType.Kind() == reflect.Interface:
		// Assign to interface{}
		result.Set(reflect.ValueOf(parsedJson))
		return nil

	case resultType.Kind() == reflect.Map && parsedType.Kind() == reflect.Map:
		// Handle map assignments
		if !result.IsValid() || result.IsNil() {
			result.Set(reflect.MakeMap(resultType))
		}

		parsedMap := parsedJson.(map[string]interface{})
		for key, value := range parsedMap {
			mapKey := reflect.ValueOf(key)
			mapValue := reflect.New(resultType.Elem()).Elem()
			if err := assignParsedValue(mapValue, value); err != nil {
				return fmt.Errorf("error assigning map value for key %s: %v", key, err)
			}
			result.SetMapIndex(mapKey, mapValue)
		}
		return nil

	case resultType.Kind() == reflect.Slice && parsedType.Kind() == reflect.Slice:
		// Handle slice assignments
		parsedSlice := parsedJson.([]interface{})
		newSlice := reflect.MakeSlice(resultType, len(parsedSlice), len(parsedSlice))

		for i, value := range parsedSlice {
			if err := assignParsedValue(newSlice.Index(i), value); err != nil {
				return fmt.Errorf("error assigning slice value at index %d: %v", i, err)
			}
		}
		result.Set(newSlice)
		return nil

	default:
		return fmt.Errorf("cannot assign %v to %v", parsedType, resultType)
	}
}
