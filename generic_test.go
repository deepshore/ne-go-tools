package negotools

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TESTCHECKGETENV_String(t *testing.T) {
	_, err := CheckGetEnvString("TESTCHECKGETENV_STRING")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_STRING", "")
	_, err = CheckGetEnvString("TESTCHECKGETENV_STRING")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_STRING", "SPAM")
	value, err := CheckGetEnvString("TESTCHECKGETENV_STRING")
	assert.NoError(t, err)
	assert.Equal(t, value, "SPAM")
	os.Unsetenv("TESTCHECKGETENV_STRING")
}

func TESTCHECKGETENV_Bool(t *testing.T) {
	_, err := CheckGetEnvBool("TESTCHECKGETENV_BOOL")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_BOOL", "")
	_, err = CheckGetEnvBool("TESTCHECKGETENV_BOOL")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_BOOL", "eggs")
	_, err = CheckGetEnvBool("TESTCHECKGETENV_BOOL")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_BOOL", "true")
	value, err := CheckGetEnvBool("TESTCHECKGETENV_BOOL")
	assert.NoError(t, err)
	assert.Equal(t, value, true)
	os.Setenv("TESTCHECKGETENV_BOOL", "false")
	value, err = CheckGetEnvBool("TESTCHECKGETENV_BOOL")
	assert.NoError(t, err)
	assert.Equal(t, value, false)
	os.Unsetenv("TESTCHECKGETENV_STRING")
}

func TESTCHECKGETENV_Int(t *testing.T) {
	_, err := CheckGetEnvInt("TESTCHECKGETENV_INT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_INT", "")
	_, err = CheckGetEnvInt("TESTCHECKGETENV_INT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_INT", "3.1415")
	_, err = CheckGetEnvInt("TESTCHECKGETENV_INT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_INT", "42")
	value, err := CheckGetEnvInt("TESTCHECKGETENV_INT")
	assert.NoError(t, err)
	assert.Equal(t, value, 42)
	os.Setenv("TESTCHECKGETENV_INT", "-24")
	value, err = CheckGetEnvInt("TESTCHECKGETENV_INT")
	assert.NoError(t, err)
	assert.Equal(t, value, -24)
	os.Unsetenv("TESTCHECKGETENV_INT")
}

func TESTCHECKGETENV_UInt(t *testing.T) {
	_, err := CheckGetEnvUInt("TESTCHECKGETENV_UINT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_UINT", "")
	_, err = CheckGetEnvUInt("TESTCHECKGETENV_UINT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_UINT", "3.1415")
	_, err = CheckGetEnvUInt("TESTCHECKGETENV_UINT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_UINT", "-24")
	_, err = CheckGetEnvUInt("TESTCHECKGETENV_UINT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_UINT", "42")
	value, err := CheckGetEnvUInt("TESTCHECKGETENV_UINT")
	assert.NoError(t, err)
	assert.Equal(t, value, uint(42))
	os.Setenv("TESTCHECKGETENV_UINT", "0")
	value, err = CheckGetEnvUInt("TESTCHECKGETENV_UINT")
	assert.NoError(t, err)
	assert.Equal(t, value, uint(0))
	os.Unsetenv("TESTCHECKGETENV_UINT")
}

func TESTCHECKGETENV_Float(t *testing.T) {
	_, err := CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_FLOAT", "")
	_, err = CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_FLOAT", "true")
	_, err = CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.Error(t, err)
	os.Setenv("TESTCHECKGETENV_FLOAT", "-24")
	value, err := CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.NoError(t, err)
	assert.Equal(t, value, -24.)
	os.Setenv("TESTCHECKGETENV_FLOAT", "42")
	value, err = CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.NoError(t, err)
	assert.Equal(t, value, 42.)
	os.Setenv("TESTCHECKGETENV_FLOAT", "6.283185")
	value, err = CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.NoError(t, err)
	assert.Equal(t, value, 6.283185)
	os.Setenv("TESTCHECKGETENV_FLOAT", "0")
	value, err = CheckGetEnvFloat("TESTCHECKGETENV_FLOAT")
	assert.NoError(t, err)
	assert.Equal(t, value, 0.)
	os.Unsetenv("TESTCHECKGETENV_FLOAT")
}
