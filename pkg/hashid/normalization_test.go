package hashid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizer2(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic normalization",
			input:    "Hello World!",
			expected: "hello-world",
		},
		{
			name:     "Unicode replacement",
			input:    "Hellö Wørld!",
			expected: "hello-world",
		},
		{
			name:     "Special characters",
			input:    "special@#-$-%^-&*-chars",
			expected: "special-dollar-percent-and-chars",
		},
		{
			name:     "Multiple spaces",
			input:    "Multiple   Spaces",
			expected: "multiple-spaces",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  Trim Spaces  ",
			expected: "trim-spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Normalizer(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizerWithSeparator2(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		separator string
		expected  string
	}{
		{
			name:      "Custom separator",
			input:     "Custom Separator",
			separator: "_",
			expected:  "custom_separator",
		},
		{
			name:      "Empty separator",
			input:     "Empty Separator",
			separator: "",
			expected:  "empty-separator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizerWithSeparator(tt.input, tt.separator)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewNormalizer(t *testing.T) {
	t.Run("Custom character map", func(t *testing.T) {
		customCharMap := map[string]string{
			"@": "at",
			"#": "hash",
			"%": "percent",
		}

		n, err := newNormalizer(customCharMap, "-")
		require.NoError(t, err)

		result, err := n.normalize("custom-@-char-#-map-%")
		require.NoError(t, err)
		assert.Equal(t, "custom-at-char-hash-map-percent", result)
	})

	t.Run("Default character map", func(t *testing.T) {
		n, err := newNormalizer(nil, "-")
		require.NoError(t, err)

		result, err := n.normalize("default @ char # map %")
		require.NoError(t, err)
		assert.Equal(t, "default-char-map-percent", result)
	})
}

func TestReplaceUnicodeChars(t *testing.T) {
	testStrings := map[string]string{
		"©":                  "(c)",
		"A81758FFFE04©E4F5":  "A81758FFFE04(c)E4F5",
		"A81758©FFFE04©E4F5": "A81758(c)FFFE04(c)E4F5",
		"₿©円₹FFFE04©E4F5":    "bitcoin(c)yenindian rupeeFFFE04(c)E4F5",
	}

	n, err := newNormalizer(nil, "-")
	require.NoError(t, err)

	for key, val := range testStrings {
		out, err := n.replaceUnicodeChars(key)
		assert.NoError(t, err)
		assert.Equal(t, val, out)
	}
}
