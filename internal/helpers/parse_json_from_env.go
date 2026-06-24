package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJSONArrayFromEnv(envVarName string) ([]string, error) {
	envValue := os.Getenv(envVarName)
	if envValue == "" {
		return nil, fmt.Errorf("environment variable %s is not set", envVarName)
	}

	var parsed []string
	if err := json.Unmarshal([]byte(envValue), &parsed); err != nil {
		return nil, fmt.Errorf("error parsing %s as JSON array: %w", envVarName, err)
	}

	return parsed, nil
}
