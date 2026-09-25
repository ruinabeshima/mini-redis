/*
	Table driven tests for Parse() function
*/

package resp

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {

	// Declare and initialise test cases
	tests := []struct {
		name          string
		input         []byte
		expectedVal   Value
		expectedBytes int
		expectedError error
	}{}

	// Loop through each test case
	for _, tt := range tests {

		// Create subtest with given name
		t.Run(tt.name, func(t *testing.T) {
			val, bytesConsumed, err := Parse(tt.input)

			// Assert error
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error: %v, actual error: %v", tt.expectedError, err)
			}

			// Assert bytes consumed
			if tt.expectedBytes != bytesConsumed {
				t.Errorf("Expected bytes consumed: %d, actual bytes consumed: %d", tt.expectedBytes, bytesConsumed)
			}

			// Assert parsed value
			if tt.expectedError == nil {
				if tt.expectedVal.Type != val.Type {
					t.Errorf("Expected type: %v, actual type: %v", tt.expectedVal.Type, val.Type)
				}

				if val.Str != tt.expectedVal.Str {
					t.Errorf("Expected string content: %q, actual string content: %q", tt.expectedVal.Str, val.Str)
				}

				if val.Int != tt.expectedVal.Int {
					t.Errorf("Expected integer value: %d, actual integer value: %d", tt.expectedVal.Int, val.Int)
				}

				if len(val.Array) != len(tt.expectedVal.Array) {
					t.Errorf("Expected array length: %d, actual array length: %d", len(tt.expectedVal.Array), len(val.Array))
				}

				if val.IsNull != tt.expectedVal.IsNull {
					t.Errorf("Expected IsNull value: %t, actual IsNull value: %t", val.IsNull, tt.expectedVal.IsNull)
				}
			}

		})
	}

}
