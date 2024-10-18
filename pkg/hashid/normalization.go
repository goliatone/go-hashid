package hashid

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var (
	removeCharList, _ = regexp.Compile(`[@#:_~.$^()!*+'"\\-]+`)
	spaceRegexp       = regexp.MustCompile(`\s+`)
)

type normalizer struct {
	charMap   map[string]string
	separator string
}

func newNormalizer(charMap map[string]string, separator string) (*normalizer, error) {
	if separator == "" {
		separator = "-"
	}

	var err error
	if charMap == nil {
		charMap, err = GetCharMap()
	}

	if err != nil {
		return nil, err
	}

	return &normalizer{
		charMap:   charMap,
		separator: separator,
	}, nil
}

// Normalizer trims and replaces spaces from the string with the separator:
//  1. Replace unicode chars (by default using charmap.json file)
//  2. remove characters not allowed
//  3. trim leading/trailing spaces
//  4. replaces any redundant whitespaces to single separator chars
//  5. lowercase
func Normalizer(s string) (string, error) {
	return NormalizerWithSeparator(s, "-")
}

// NormalizerWithSeparator will normalize the string
func NormalizerWithSeparator(s, separator string) (string, error) {
	n, err := newNormalizer(nil, separator)
	if err != nil {
		return "", err
	}
	return n.normalize(s)
}

func NormalizerWithCharMap(s string, m map[string]string) (string, error) {
	n, err := newNormalizer(m, "-")
	if err != nil {
		return "", err
	}
	return n.normalize(s)
}

func (n *normalizer) normalize(s string) (string, error) {
	s = unicodeNorm(s)

	var result strings.Builder

	for _, ch := range s {
		char := string(ch)

		appendChar, ok := n.charMap[char]
		if !ok {
			appendChar = char
		}

		if appendChar == n.separator {
			appendChar = " "
		}

		cleanChar := removeCharList.ReplaceAllString(appendChar, "")
		result.WriteString(cleanChar)
	}

	out := result.String()

	out = strings.TrimSpace(out)

	out = spaceRegexp.ReplaceAllString(out, n.separator)

	out = strings.ToLower(out)

	return out, nil
}

func (n *normalizer) replaceUnicodeChars(s string) (string, error) {
	for k, v := range n.charMap {
		s = strings.ReplaceAll(s, k, v)
	}
	return s, nil
}

func removeCharsNotAllowed(s string) string {
	return removeCharList.ReplaceAllString(s, "")
}

func removeSpaces(s string) string {
	return spaceRegexp.ReplaceAllString(s, "-")
}

// Use NFC (Normalization Form C)
// Canonical Composition – precomposed characters
// where possible (e.g., é as a single character
// rather than ‘e’ + combining accent)
func unicodeNorm(s string) string {
	return norm.NFC.String(s)
}
