// Package hashid provides utilities for generating deterministic
// globally unique identifiers (UUIDs) using various hashing algorithms.
//
// It supports multiple hashing algorithms such as MD5, SHA1, SHA256, and HMAC-SHA256,
// with optional string normalization. The library allows generation of UUIDs based on
// specific attributes, ensuring consistency across different systems.
//
// Usage:
//
//	id, err := hashid.New("user@example.com", hashid.WithHashAlgorithm(hashid.SHA256))
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println(id)
package hashid
