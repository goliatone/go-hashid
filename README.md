# hashid

`hashid` is a Go library for generating deterministic, globally unique identifiers (UUIDs) from input strings using various hashing algorithms. It supports UUID versions 3, 5, and 8, enabling use cases such as test fixtures, entity identification, cache key generation, and privacy-preserving identifiers.

**Key Features**:

- Deterministic output: the same input string always produces the same UUID.
- Support for multiple hashing algorithms (MD5, SHA1, SHA256, HMAC-SHA256).
- Configurable string normalization (e.g., case folding, trimming whitespace).
- UUID versions 3, 5 (RFC 4122), and 8 (RFC 9562).
- Thread-safe and no external dependencies.


## Installation

```bash
go get github.com/goliatone/hashid
```

## Usage

### Library

You can import the library and use it in your projects.

```go
import "github.com/goliatone/hashid"

func main() {
    // Basic usage with defaults (MD5, UUID v3)
    uuid, err := hashid.New("my-input-string")

    // Using SHA1 (UUID v5)
    uuid, err = hashid.New("my-input-string",
        hashid.WithHashAlgorithm(hashid.SHA1))

    // Using HMAC-SHA256 (UUID v8)
    key := []byte("secret-key")
    uuid, err = hashid.New("my-input-string",
        hashid.WithHashAlgorithm(hashid.HMAC_SHA256),
        hashid.WithHMACKey(key))

    // Disable normalization
    uuid, err = hashid.New("My-Input-String",
        hashid.WithNormalization(false))
}
```

### CLI

```bash
# Basic usage
hashid "my-input-string"

# Using SHA1
hashid -hash sha1 "my-input-string"

# Using HMAC
hashid -hash hmac -key mysecret "my-input-string"

# No need to specify hash if key given
hashid -key mysecret "my-input-string"

# Custom normalization
hashid -normalize upper "my-input-string"
```

## Implementation Details

- Supports MD5 (UUID v3), SHA1 (UUID v5), and HMAC-SHA256 (UUID v8) algorithms
- Implements RFC 4122 for UUID versions 3 and 5
- Implements RFC 9562 for UUID version 8 (custom format)
- Thread-safe
- No external dependencies
- Configurable string normalization
- Error handling for invalid configurations

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request

## License

MIT License - see LICENSE file for details
