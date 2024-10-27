package main

import (
	"log"
	"os"

	"github.com/LitoleNINJA/json-parser/cmd/encoder"
	"github.com/LitoleNINJA/json-parser/cmd/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Input JSON filename missing !")
	}

	fileName := os.Args[1]

	fileData, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading file : %v", err)
	}

	// decode JSON
	var result interface{}
	err = parser.ParseJSON(fileData, &result)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	// you can cast the returned interface into any of the below type if compatible
	switch result := result.(type) {
	case map[string]interface{}:
		log.Printf("JSON Object: %+v", result)
	case []interface{}:
		log.Println("JSON Array:", result)
	case float64:
		log.Println("JSON Number:", result)
	case string:
		log.Println("JSON String:", result)
	case bool:
		log.Println("JSON Bool:", result)
	default:
		log.Fatalf("Unknown JSON structure")
	}

	// encode JSON
	json, err := encoder.EncodeJSON(result)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	log.Printf("Encoded JSON : %s", json)

	err = writeToFile("test-encode.json", json)
	if err != nil {
		log.Fatalf("Error Writing to file: %v", err)
	}
}

func writeToFile(fileName string, json []byte) error {
	// Create or open the file
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(json)
	if err != nil {
		return err
	}

	return nil
}
