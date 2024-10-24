package main

import (
	"log"
	"os"

	"github.com/LitoleNINJA/json-parser/cmd/parser"
	"github.com/LitoleNINJA/json-parser/cmd/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Input JSON filename missing !")
	}

	fileName := os.Args[1]
	// log.Printf("Filename : %s", fileName)

	fileData, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading file : %v", err)
	}

	log.Print("Starting tokenizer ... ")
	tokens, err := tokenizer.Tokenize(fileData)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting parser ...")
	result, err := parser.Parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	switch result := result.(type) {
	case parser.JsonObject:
		log.Printf("JSON Object: %+v", result)
	case parser.JsonArray:
		log.Println("JSON Array:", result)
	default:
		log.Fatalf("Unknown JSON structure")
	}
}
