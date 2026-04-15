package env

import (
	"bufio"
	"os"
	"strings"
)

// MergeMode controls how existing .env values are handled during a sync.
type MergeMode string

const (
	// MergeOverwrite replaces all existing values with vault values.
	MergeOverwrite MergeMode = "overwrite"
	// MergeKeepExisting preserves local values when a key already exists.
	MergeKeepExisting MergeMode = "keep-existing"
)

// Merge combines existing .env file contents with incoming vault secrets
// according to the specified MergeMode. Returns the merged map.
func Merge(existingPath string, incoming map[string]string, mode MergeMode) (map[string]string, error) {
	existing, err := parseEnvFile(existingPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	result := make(map[string]string, len(incoming))
	for k, v := range incoming {
		result[k] = v
	}

	if mode == MergeKeepExisting {
		for k, v := range existing {
			if _, found := result[k]; found {
				result[k] = v
			}
		}
	}

	return result, nil
}

// parseEnvFile reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' and empty lines are ignored.
func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		result[key] = val
	}

	return result, scanner.Err()
}
