package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseJSONArrayFromEnv retrieves an environment variable containing a JSON array
// and parses it into a Go slice of strings. Returns an error if parsing fails.
func ParseJSONArrayFromEnv(envVarName string) ([]string, error) {
	envValue := os.Getenv(envVarName)
	if envValue == "" {
		return nil, fmt.Errorf("Environment variable %s does not exist", envVarName)
	}

	var parsedArray []string
	err := json.Unmarshal([]byte(envValue), &parsedArray); 
	if err != nil {
		return nil, fmt.Errorf("Error parsing environment variable %s as JSON array: %w", envVarName, err)
	}

	return parsedArray, nil
}
