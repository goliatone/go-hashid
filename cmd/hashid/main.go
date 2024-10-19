package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goliatone/hashid/pkg/hashid"
	"github.com/goliatone/hashid/pkg/version"
)

type config struct {
	algorithm   string
	hmacKey     string
	noNormalize bool
	uuidVersion int
	showVersion bool
	charmapFile string
}

func main() {
	conf := parseFlags()

	if conf.showVersion {
		version.Print(os.Stdout)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		fmt.Fprint(os.Stderr, "Error: Input string is required\n\n")
		usage()
		os.Exit(1)
	}

	input := strings.Join(flag.Args(), " ")

	options := []hashid.Option{}

	switch strings.ToLower(conf.algorithm) {
	case "md5":
		options = append(options, hashid.WithHashAlgorithm(hashid.MD5))
	case "sha1":
		options = append(options, hashid.WithHashAlgorithm(hashid.SHA1))
	case "sha256":
		options = append(options, hashid.WithHashAlgorithm(hashid.SHA256))
	case "hmac":
		if conf.hmacKey == "" {
			fmt.Fprintln(os.Stderr, "Error: HMAC key is required when using HMAC algorithm")
			os.Exit(1)
		}
		options = append(options,
			hashid.WithHashAlgorithm(hashid.HMAC_SHA256),
			hashid.WithHMACKey([]byte(conf.hmacKey)))
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported hashing algorithm: %s\n", conf.algorithm)
		os.Exit(1)
	}

	if conf.noNormalize {
		options = append(options, hashid.WithNormalization(false))
	}

	switch conf.uuidVersion {
	case 3, 5, 8:
		options = append(options, hashid.WithUUIDVersion(conf.uuidVersion))
	case 0:
		options = append(options, hashid.WithUUIDVersion(3))
	default:
		fmt.Printf("Unsupported UUID version: %d\n", conf.uuidVersion)
		os.Exit(1)
	}

	uuid, err := hashid.New(input, options...)
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(uuid)
}

func parseFlags() config {
	conf := config{}

	flag.StringVar(&conf.algorithm, "hash", "md5", "Hashing algorithm (md5, sha1, sha256, hmac)")
	flag.StringVar(&conf.hmacKey, "key", "", "HMAC key (required when using hmac algorithm)")
	flag.BoolVar(&conf.noNormalize, "no-normalize", false, "Disable string normalization")
	flag.IntVar(&conf.uuidVersion, "uuid-version", 0, "Force specific UUID version (3, 5, or 8)")
	flag.BoolVar(&conf.showVersion, "version", false, "Show version information")
	flag.StringVar(&conf.charmapFile, "charmap", "", "Path to custom character mapping JSON file")

	flag.Usage = usage

	flag.Parse()
	return conf
}

func loadCustomCharMap(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("filed to read charmap file: %w", err)
	}
	var mapping map[string]string
	if err := json.Unmarshal(data, &mapping); err != nil {
		return nil, fmt.Errorf("filed to parse charmap file: %w", err)
	}
	return mapping, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: hashid [options] <input-string>

Options:
  -hash string
        Hashing algorithm (md5, sha1, sha256, hmac) (default "md5")
  -key string
        HMAC key (required when using hmac algorithm)
  -no-normalize
        Disable string normalization
  -uuid-version int
        Force specific UUID version (3, 5, or 8) (default 3)
  -version
        Show version information

Examples:
  hashid "user@example.com"
  hashid -hash sha1 "user@example.com"
  hashid -hash hmac -key mysecret "user@example.com"
  hashid -no-normalize "user@example.com"
  hashid -normalize upper "user@example.com"
  hashid -uuid-version 8 "user@example.com"

Version:
  %s

`, version.GetVersion())
}
