package negotools

import (
	"fmt"
	"hash/crc32"
	"os"

	passwordGenerator "github.com/m1/go-generate-password/generator"
)

func CRC32Checksum(s string) string {
	table := crc32.MakeTable(crc32.IEEE) // ISO 3309 (HDLC) polynomial
	checksum := crc32.Checksum([]byte(s), table)
	return fmt.Sprintf("%X", checksum)
}

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
