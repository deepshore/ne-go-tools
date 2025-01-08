package negotools

import (
	"strings"
	"testing"
)

func TestCRC32Checksum(t *testing.T) {
	const inputString string = "func TestCRC32Checksum(t *testing.T)"
	const expected string = "460DF926"
	var actual string = CRC32Checksum(inputString)
	if actual != expected {
		t.Errorf("Expected function name %s, but got %s", expected, actual)
	}
}

func TestGeneratePassword(t *testing.T) {
	var passwordLengths []uint = []uint{3, 12, 42, 97}
	for _, length := range passwordLengths {
		password, _ := GeneratePassword(length, true)
		actualLength := uint(len(password))
		if actualLength != length {
			t.Errorf("Expected function name %v, but got %v", length, actualLength)
		}
		if strings.Contains(password, ",") {
			t.Errorf("Password (%q) contains invalid character ','.", password)
		}
		if strings.Contains(password, ":") {
			t.Errorf("Password (%q) contains invalid character ':'.", password)
		}
	}
}
