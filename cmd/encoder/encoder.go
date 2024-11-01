package encoder

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/LitoleNINJA/json-parser/cmd/customError"
)

func EncodeJSON(jsonData interface{}, showLogs bool) ([]byte, error) {
	if !showLogs {
		log.SetOutput(io.Discard)
	}

	encodedJson, err := encode(reflect.ValueOf(jsonData), "")
	if err != nil {
		return nil, customError.NewError(fmt.Errorf("error while encoding JSON : %v", err))
	}

	return encodedJson, err
}

func encode(value reflect.Value, indent string) ([]byte, error) {
	// Handle nil values
	if !value.IsValid() {
		return []byte("null"), nil
	}

	switch value.Kind() {
	case reflect.Map:
		return encodeMap(&value, indent)
	case reflect.Interface:
		return encode(value.Elem(), indent)
	case reflect.Array, reflect.Slice:
		return encodeArray(&value)
	case reflect.String:
		return encodeString(&value)
	case reflect.Float64:
		return encodeNumber(&value)
	case reflect.Bool:
		return encodeBool(&value)
	default:
		return nil, customError.NewError(fmt.Errorf("unsupported type: %s", value.Type()))
	}
}

func encodeMap(value *reflect.Value, indent string) ([]byte, error) {
	var mapEncoding bytes.Buffer

	mapEncoding.Write([]byte("{\n"))
	indent += "\t"
	for i, k := range value.MapKeys() {
		if i > 0 {
			mapEncoding.Write([]byte(",\n"))
		}

		v := value.MapIndex(k)
		// log.Printf("Key: %s, Value: %s", k, v)

		encodedKey, err := encode(k, indent)
		if err != nil {
			return nil, customError.NewError(fmt.Errorf("error while encoding map key '%s' : %v", k, err))
		}

		encodedValue, err := encode(v, indent)
		if err != nil {
			return nil, customError.NewError(fmt.Errorf("error while encoding map value '%s' : %v", k, err))
		}

		mapEncoding.Write([]byte(indent))
		mapEncoding.Write(encodedKey)
		mapEncoding.Write([]byte(": "))
		mapEncoding.Write(encodedValue)
	}

	mapEncoding.Write([]byte("\n" + strings.TrimSuffix(indent, "\t") + "}"))

	log.Printf("Encoded Map: %s", mapEncoding.Bytes())
	return mapEncoding.Bytes(), nil
}

func encodeArray(value *reflect.Value) ([]byte, error) {
	var arrayEncoding bytes.Buffer

	arrayEncoding.WriteByte('[')
	for i := 0; i < value.Len(); i++ {
		if i > 0 {
			arrayEncoding.Write([]byte(", "))
		}

		encodedValue, err := encode(value.Index(i), "")
		if err != nil {
			return nil, customError.NewError(fmt.Errorf("error while encoding array value '%s' : %v", value.Index(i), err))
		}

		arrayEncoding.Write(encodedValue)
	}
	arrayEncoding.WriteByte(']')

	log.Printf("Encoded Array: %s", arrayEncoding.Bytes())
	return arrayEncoding.Bytes(), nil
}

func encodeString(value *reflect.Value) ([]byte, error) {
	var stringEncoder bytes.Buffer

	stringEncoder.WriteByte('"')
	for _, r := range value.String() {
		switch r {
		case '"':
			stringEncoder.WriteString(`\"`)
		case '\\':
			stringEncoder.WriteString(`\\`)
		case '\n':
			stringEncoder.WriteString(`\n`)
		case '\r':
			stringEncoder.WriteString(`\r`)
		case '\t':
			stringEncoder.WriteString(`\t`)
		default:
			if r < 0x20 {
				stringEncoder.WriteString(fmt.Sprintf(`\u%04x`, r))
			} else {
				stringEncoder.WriteRune(r)
			}
		}
	}

	stringEncoder.WriteByte('"')
	return stringEncoder.Bytes(), nil
}

func encodeNumber(value *reflect.Value) ([]byte, error) {
	number := value.Float()

	// Convert to scientific notation if number is very large or very small
	abs := math.Abs(number)
	if abs < 1e-6 || abs >= 1e21 {
		// Use 'e' notation with 6 decimal places
		return []byte(strconv.FormatFloat(number, 'e', 1, 64)), nil
	}

	// Use standard notation for normal numbers
	// The -1 precision argument tells FormatFloat to use the smallest number
	// of decimal places needed to represent the number exactly
	return []byte(strconv.FormatFloat(number, 'f', -1, 64)), nil
}

func encodeBool(value *reflect.Value) ([]byte, error) {
	return []byte(strconv.FormatBool(value.Bool())), nil
}
