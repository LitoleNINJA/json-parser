package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LitoleNINJA/json-parser/cmd/parser"
)

func TestJsonParser(t *testing.T) {
	testDir := "../../test/JSONTestSuite/test_parsing"
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	// fileNames := []string{}

	// for _, fileName := range fileNames {
	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(testDir, fileName)
		isAcceptedTest := true

		// skipping i_ tests
		if strings.HasPrefix(fileName, "y_") {
			isAcceptedTest = true
		} else if strings.HasPrefix(fileName, "n_") {
			isAcceptedTest = false
		}

		t.Run(fileName, func(t *testing.T) {
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", fileName, err)
				return
			}

			_, err = parser.ParseJSON(data)
			if err == nil && !isAcceptedTest {
				t.Errorf("Expected failure for %s | Testcase : %s", fileName, string(data))
			} else if err != nil && isAcceptedTest {
				t.Errorf("Expected success for %s, got %v | Testcase : %s", fileName, err, string(data))
			}
		})
	}
}
