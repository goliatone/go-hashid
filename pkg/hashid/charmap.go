package hashid

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed charmap.json
var defaultCharMapFS embed.FS

var (
	charMap     map[string]string
	charMapOnce sync.Once
	initError   error
)

func GetCharMap() (map[string]string, error) {
	charMapOnce.Do(initDefaultCharMap)

	if initError != nil {
		return nil, initError
	}

	// Create a copy to prevent external modification
	copy := make(map[string]string, len(charMap))
	for k, v := range charMap {
		copy[k] = v
	}
	return copy, nil
}

func SetCharMap(mapping map[string]string) {
	charMap = mapping
	initError = nil
}

func ResetCharMap() error {
	initError = nil
	initDefaultCharMap()
	return initError
}

func initDefaultCharMap() {
	defaultCharMapData, err := defaultCharMapFS.ReadFile("charmap.json")
	if err != nil {
		initError = fmt.Errorf("failed to open default charmap: %w", err)
		return
	}

	charMap, err = loadCharMap(defaultCharMapData)
	if err != nil {
		initError = fmt.Errorf("failed to load default charmap: %w", err)
	}
}

func loadCharMap(data []byte) (map[string]string, error) {
	var mapping map[string]string
	err := json.Unmarshal(data, &mapping)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal charmap: %w", err)
	}
	return mapping, nil
}
