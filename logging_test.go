package negotools

import (
	"errors"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/assert"
)

func TestGetFunctionName(t *testing.T) {
	expected := "ne-go-tools.TestGetFunctionName"
	actual := stripPath(GetFunctionName(1))
	if actual != expected {
		t.Errorf("Expected function name %s, but got %s", expected, actual)
	}
}

func TestAnotherFunction(t *testing.T) {
	expected := "ne-go-tools.TestAnotherFunction"
	actual := stripPath(GetFunctionName(1))

	if actual != expected {
		t.Errorf("Expected function name %s, but got %s", expected, actual)
	}
}

func TestLogDebug(t *testing.T) {
	hook := test.NewGlobal()

	// Set log level to Debug
	log.SetLevel(log.DebugLevel)

	LogDebug("This is a debug message", "Key1", "Value1", "Key2", "Value2")

	// Verify that there is exactly one log entry
	assert.Equal(t, 1, len(hook.Entries), "Expected exactly one log entry")

	// Get the last log entry
	entry := hook.LastEntry()

	// Verify the log entry details
	assert.Equal(t, log.DebugLevel, entry.Level)
	assert.Equal(t, "This is a debug message", entry.Message)
	assert.Equal(t, "ne-go-tools.TestLogDebug", entry.Data["Function"])
	assert.Equal(t, "Value1", entry.Data["Key1"])
	assert.Equal(t, "Value2", entry.Data["Key2"])

	// Clean up the hook
	hook.Reset()
}

func TestLogError(t *testing.T) {
	hook := test.NewGlobal()
	err := errors.New("test error")

	// Set log level to Error (though not strictly necessary here)
	log.SetLevel(log.ErrorLevel)

	LogError("This is an error message", err, "Key1", "Value1", "Key2", "Value2")

	// Verify that there is exactly one log entry
	assert.Equal(t, 1, len(hook.Entries), "Expected exactly one log entry")

	// Get the last log entry
	entry := hook.LastEntry()

	// Verify the log entry details
	assert.Equal(t, log.ErrorLevel, entry.Level)
	assert.Equal(t, "This is an error message", entry.Message)
	assert.Equal(t, "ne-go-tools.TestLogError", entry.Data["Function"])
	assert.Equal(t, err, entry.Data["Error"])
	assert.Equal(t, "Value1", entry.Data["Key1"])
	assert.Equal(t, "Value2", entry.Data["Key2"])

	// Clean up the hook
	hook.Reset()
}
