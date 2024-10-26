package main

import (
	"log"
	"os"

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

	result, err := parser.ParseJSON(fileData)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	switch result := result.(type) {
	case parser.JsonObject:
		log.Printf("JSON Object: %+v", result)
	case parser.JsonArray:
		log.Println("JSON Array:", result)
	case float64:
		log.Println("JSON Number:", result)
	default:
		log.Fatalf("Unknown JSON structure")
	}
}
