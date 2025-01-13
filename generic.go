package negotools

import (
	"errors"
	"fmt"
	"hash/crc32"
	"os"

	passwordGenerator "github.com/m1/go-generate-password/generator"
	"gitlab.com/avarf/getenvs"
)

// Calculates the CRC32 checksum (using the ISO 3309-HDLC polynomial setting)
// of a string and returns the resulting number as uppercase HEX string.
func CRC32Checksum(s string) string {
	table := crc32.MakeTable(crc32.IEEE) // ISO 3309 (HDLC) polynomial
	checksum := crc32.Checksum([]byte(s), table)
	return fmt.Sprintf("%X", checksum)
}

// Generates a password of given length, optionally in-/excluding ambigous
// characters ("<>[](){}:;'/|\\,") which can cause issues in some applications
func GeneratePassword(length uint, excludeAmbiguousChars bool) (password string, err error) {

	config := passwordGenerator.Config{
		Length:                     length,
		IncludeSymbols:             true,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    true,
		ExcludeSimilarCharacters:   false,
		ExcludeAmbiguousCharacters: excludeAmbiguousChars,
	}
	generator, err := passwordGenerator.New(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialise password generator: %v", err)
		return "", err
	}

	temp, err := generator.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate password: %v", err)
		return "", err
	}
	password = *temp

	return password, err
}

// Wrappers for 'getenvs' functions to force non-emptiness and existence of ENVs
// Also adds a function to load uints.

func CheckGetEnvString(key string) (value string, err error) {
	stringValue, ok := os.LookupEnv(key)
	if !ok || stringValue == "" {
		fmt.Fprintf(os.Stderr, "Error: ENV %q must be set and not empty!", key)
		return "", errors.New("ValueError")
	}
	value = getenvs.GetEnvString(key, "")
	return value, nil
}

func CheckGetEnvBool(key string) (value bool, err error) {
	stringValue, ok := os.LookupEnv(key)
	if !ok || stringValue == "" {
		fmt.Fprintf(os.Stderr, "Error: ENV %q must be set and not empty.", key)
		return false, errors.New("ValueError")
	}
	value, err = getenvs.GetEnvBool(key, false)
	return value, err
}

func CheckGetEnvInt(key string) (value int, err error) {
	stringValue, ok := os.LookupEnv(key)
	if !ok || stringValue == "" {
		fmt.Fprintf(os.Stderr, "Error: ENV %q must be set and not empty.", key)
		return 0, errors.New("ValueError")
	}
	value, err = getenvs.GetEnvInt(key, 0)
	return value, err
}

func CheckGetEnvUInt(key string) (value uint, err error) {
	stringValue, ok := os.LookupEnv(key)
	if !ok || stringValue == "" {
		fmt.Fprintf(os.Stderr, "Error: ENV %q must be set and not empty.", key)
		return 0, errors.New("ValueError")
	}
	intValue, err := getenvs.GetEnvInt(key, 0)
	if intValue < 0 {
		fmt.Fprintf(os.Stderr, "Error: uint ENV %q must not be negative, is %v", key, intValue)
		return 0, errors.New("ValueError")
	}
	value = uint(intValue)
	return value, err
}

func CheckGetEnvFloat(key string) (value float64, err error) {
	stringValue, ok := os.LookupEnv(key)
	if !ok || stringValue == "" {
		fmt.Fprintf(os.Stderr, "Error: ENV %q must be set and not empty.", key)
		return 0, errors.New("ValueError")
	}
	value, err = getenvs.GetEnvFloat(key, 0.)
	return value, err
}
