package hashid

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
)

type HashAlgorithm string

const (
	MD5         HashAlgorithm = "md5"
	SHA1        HashAlgorithm = "sha1"
	SHA256      HashAlgorithm = "sha256"
	HMAC_SHA256 HashAlgorithm = "hmac"
)

type options struct {
	hashAlgo    HashAlgorithm
	normalize   bool
	normalizer  func(string) (string, error)
	uuidVersion int
	hmacKey     []byte
	charMap     map[string]string
}

type Option func(*options)

func defaultOptions() options {
	return options{
		hashAlgo:    MD5,
		normalize:   true,
		normalizer:  Normalizer,
		uuidVersion: 3,
		hmacKey:     nil,
		charMap:     nil,
	}
}

func WithHashAlgorithm(algo HashAlgorithm) Option {
	return func(o *options) {
		o.hashAlgo = algo
		switch algo {
		case SHA1:
			o.uuidVersion = 5
		case HMAC_SHA256:
			o.uuidVersion = 8
		default:
			o.uuidVersion = 3
		}
	}
}

func WithCustomNormalizer(normalizer func(string) (string, error)) Option {
	return func(o *options) {
		o.normalizer = normalizer
	}
}

func WithHMACKey(key []byte) Option {
	return func(o *options) {
		o.hmacKey = key
		o.hashAlgo = HMAC_SHA256
	}
}

func WithNormalization(normalize bool) Option {
	return func(o *options) {
		o.normalize = normalize
	}
}

// WithUUIDVersion allows explicitly setting the UUID version
func WithUUIDVersion(version int) Option {
	return func(o *options) {
		o.uuidVersion = version
	}
}

func WithCustomCharMap(mapping map[string]string) Option {
	return func(o *options) {
		o.charMap = mapping
	}
}

func New(input string, opts ...Option) (string, error) {
	config := defaultOptions()

	for _, opt := range opts {
		opt(&config)
	}

	if config.hashAlgo == HMAC_SHA256 && config.hmacKey == nil {
		return "", fmt.Errorf("HMAC key is required when using HMAC_SHA256")
	}

	uuidVersion := 0
	switch config.uuidVersion {
	case 3, 5, 8:
		uuidVersion = config.uuidVersion
	case 0:
		uuidVersion = 3
	default:
		return "", fmt.Errorf("UUID version should be one of 3, 5, 8")
	}

	var normalizer func(string) (string, error)
	if config.charMap != nil {
		n, err := newNormalizer(config.charMap, "-")
		if err != nil {
			return "", err
		}
		normalizer = n.normalize
	} else {
		normalizer = config.normalizer
	}

	var err error

	if config.normalize {
		input, err = normalizer(input)
	}

	if err != nil {
		return "", fmt.Errorf("normalization error: %w", err)
	}

	hasher, err := getHasher(config.hashAlgo, config.hmacKey)
	if err != nil {
		return "", err
	}
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	return formatUUID(hash, uuidVersion), nil
}

func getHasher(algo HashAlgorithm, key []byte) (hash.Hash, error) {
	switch algo {
	case SHA1:
		return sha1.New(), nil
	case SHA256:
		return sha256.New(), nil
	case HMAC_SHA256:
		if key == nil {
			return nil, fmt.Errorf("HMAC key is required when using HMAC_SHA256")
		}
		return hmac.New(sha256.New, key), nil
	default:
		return md5.New(), nil
	}
}

// Valid UUID version 3, following the format
// xxxxxxxx-xxxx-3xxx-yxxx-xxxxxxxxxxxx
// Where:
// - The 13th character is always 3 (indicating a version 3 UUID).
// - The 17th character is one of 8, 9, a, or b.
func formatUUID(hash []byte, version int) string {
	var versionByte uint8
	if version == 8 {
		// For version 8, we set the version but don't modify any other bits
		// as per RFC 9562, allowing for custom formats
		versionByte = hash[6]&0x0F | 0x80
	} else {
		// For other versions (3, 5), use traditional version formatting
		versionByte = hash[6]&0x0F | uint8(version<<4)
	}

	// Set the variant to RFC 4122 (the 2 most significant bits should be 10)
	variantByte := hash[8]&0x3F | 0x80

	// Return the UUID formatted string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		hash[0:4],
		hash[4:6],
		[]byte{versionByte, hash[7]},
		[]byte{variantByte, hash[9]},
		hash[10:16],
	)
}
