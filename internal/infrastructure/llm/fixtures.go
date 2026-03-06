package llm

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFixture reads a JSON fixture file and unmarshals it into the given type.
func LoadFixture[T any](path string) (T, error) {
	var result T
	data, err := os.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("llm.LoadFixture: %w", err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("llm.LoadFixture: %w", err)
	}
	return result, nil
}
