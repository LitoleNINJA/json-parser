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

	jsonObject := result.(parser.JsonObject)

	log.Printf("Key : name, Value : %v", jsonObject["name"])
}
