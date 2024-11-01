package encoder_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LitoleNINJA/json-parser/cmd/encoder"
	"github.com/LitoleNINJA/json-parser/cmd/parser"
)

func TestJSONEncoder(t *testing.T) {
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

			var result any
			// decode and encode
			parserErr := parser.ParseJSON(data, &result, true)
			_, encoderErr := encoder.EncodeJSON(result, true)

			if isAcceptedTest && (parserErr != nil || encoderErr != nil) {
				t.Errorf("Expected success for Testcase : %s\nParseErr: %v, EncodeErr: %v", string(data), parserErr, encoderErr)
			} else if !isAcceptedTest && parserErr == nil && encoderErr == nil {
				t.Errorf("Expected failure for Testcase : %s", string(data))
			}
		})
	}
}
