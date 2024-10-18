package hashid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCharMap2(t *testing.T) {
	// Test default charmap
	defaultMap, err := GetCharMap()
	assert.NoError(t, err)
	assert.NotEmpty(t, defaultMap)
	assert.Equal(t, "dollar", defaultMap["$"])

	// Test custom charmap
	customMap := map[string]string{
		"@": "at",
		"#": "hash",
	}
	SetCharMap(customMap)

	retrievedMap, err := GetCharMap()
	assert.NoError(t, err)
	assert.Equal(t, "at", retrievedMap["@"])
	assert.Equal(t, "hash", retrievedMap["#"])

	// Verify that modifying the retrieved map doesn't affect the internal map
	retrievedMap["@"] = "modified"

	nm, err := GetCharMap()
	assert.NoError(t, err)
	assert.NotEqual(t, "modified", nm["@"])

	ResetCharMap()
}

func TestGetCharMap(t *testing.T) {
	testCharMap, err := GetCharMap()
	assert.NoError(t, err)
	assert.NotEmpty(t, testCharMap)

	// Set the test charmap
	SetCharMap(testCharMap)

	// Now test GetCharMap
	charMap, err := GetCharMap()
	assert.NoError(t, err)
	assert.NotEmpty(t, charMap)
	assert.Equal(t, "dollar", charMap["$"])

	ResetCharMap()
}

func TestSetCharMap(t *testing.T) {
	originalMap, err := GetCharMap()
	assert.NoError(t, err)

	customMap := map[string]string{
		"$": "custom_dollar",
		"@": "custom_at",
	}

	SetCharMap(customMap)

	updatedMap, err := GetCharMap()
	assert.NoError(t, err)

	assert.Equal(t, "custom_at", updatedMap["@"])
	assert.Equal(t, "custom_dollar", updatedMap["$"])

	ResetCharMap()
	updatedMap, err = GetCharMap()
	assert.NoError(t, err)

	assert.Equal(t, originalMap["$"], updatedMap["$"])

	ResetCharMap()
}

func TestConcurrentCharMapAccess(t *testing.T) {
	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			GetCharMap()
		}
		done <- true
	}()

	for i := 0; i < 100; i++ {
		customMap := map[string]string{
			"$": "dollar",
			"@": "at",
		}
		SetCharMap(customMap)
	}

	<-done

	ResetCharMap()
}

func TestWithCustomCharMap(t *testing.T) {
	customMap := map[string]string{
		"@": "at",
		"#": "hash",
	}

	opt := WithCustomCharMap(customMap)
	config := defaultOptions()
	opt(&config)

	assert.Equal(t, customMap, config.charMap)

	ResetCharMap()
}
