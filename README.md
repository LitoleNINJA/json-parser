# JSON Parser & Encoder
### Go package to serialize and de-serialize JSON 

This project is a simple JSON parser and encoder written in Go. It provides a set of functions to parse JSON data into a Go data structure and vice versa. The project includes the following components:

- **Parser**: The parser package is responsible for parsing JSON data and converting it into a Go data structure.

- **Encoder**: The encoder package is responsible for encoding a Go data structure back into a JSON string.

## Features

- Parses JSON data into a Go data structure (map, slice, etc.)

- Encodes Go data structures back into JSON format

- Supports all JSON data types (string, number, boolean, null, object, array)

- Handles invalid JSON data and provides detailed error messages

- Includes a comprehensive test suite to ensure correctness

## Usage

### Parsing JSON
```
func ParseJSON(fileData []byte, result any, showLogs bool) error
```
- **fileData []byte**: The JSON data to be parsed, in the form of a byte slice.

- **result any**: A pointer to a variable that will hold the parsed JSON data. This can be of any type, such as map[string]interface{}, []interface{}, or a custom struct.

- **showLogs bool**: A flag that determines whether to show log messages during the parsing process. If set to false, the function will suppress all log output.

<br/>

```
import (
    "github.com/LitoleNINJA/json-parser/cmd/parser"
)

var result any
err := parser.ParseJSON([]byte(jsonData), &result, false)
if err != nil {
    // Handle error
}

// result now contains the parsed JSON data
```

### Encoding JSON
```
func EncodeJSON(jsonData interface{}, showLogs bool) ([]byte, error)
```
- **jsonData interface{}**: The Go data structure to be encoded, such as a map[string]interface{}, []interface{}, or a custom struct.

- **showLogs bool**: A flag that determines whether to show log messages during the encoding process. If set to false, the function will suppress all log output.

- **[]byte**: The encoded JSON data as a byte slice.

<br />

```
import (
    "github.com/LitoleNINJA/json-parser/cmd/encoder"
)

jsonData := map[string]interface{}{
    "name": "John Doe",
    "age":  30,
    "hobbies": []string{"reading", "hiking"},
}

encodedJSON, err := encoder.EncodeJSON(jsonData, false)
if err != nil {
    // Handle error
}

// encodedJSON now contains the JSON representation of the data
```

## Future Improvements

- Optimize the performance of the parser and encoder (Buffer Pooling, Pre-allocation, Optimized String Handling, Faster Type Checking, Memory Optimization)

- Add support for more advanced JSON features (e.g., comments, trailing commas)