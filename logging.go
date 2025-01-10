package negotools

import (
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func SetLoglevel(logLevel string) (err error) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
		return
	}
	log.SetLevel(level)
	return
}

// stripPath removes the package path prefix from the function name
func stripPath(fullName string) string {
	parts := strings.Split(fullName, "/")
	return parts[len(parts)-1]
}

// GetFunction feturns the function name by functionLvl. 1 = one level up, 2 = two levels up from current function
func GetFunctionName(functionLvl int) string {
	// pc = program counter
	pc, _, _, _ := runtime.Caller(functionLvl)
	return stripPath(runtime.FuncForPC(pc).Name())
}

// debugLog logs debug messages with contextual information.
func LogTrace(message string, fields ...interface{}) {
	logFields := log.Fields{
		"Function": GetFunctionName(2),
	}
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		logFields[key] = value
	}
	log.WithFields(logFields).Debug(message)
}

// debugLog logs debug messages with contextual information.
func LogDebug(message string, fields ...interface{}) {
	logFields := log.Fields{
		"Function": GetFunctionName(2),
	}
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		logFields[key] = value
	}
	log.WithFields(logFields).Debug(message)
}

// LogInfo logs info messages with contextual information.
func LogInfo(message string, fields ...interface{}) {
	logFields := log.Fields{
		"Function": GetFunctionName(2),
	}
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		logFields[key] = value
	}
	log.WithFields(logFields).Info(message)
}

// LogError logs an error with contextual information.
func LogError(message string, err error, fields ...interface{}) {
	logFields := log.Fields{
		"Function": GetFunctionName(2),
		"Error":    err,
	}
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		logFields[key] = value
	}
	log.WithFields(logFields).Error(message)
}

// LogWarning logs an warning with contextual information.
func LogWarning(message string, err error, fields ...interface{}) {
	logFields := log.Fields{
		"Function": GetFunctionName(2),
		"Error":    err,
	}
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		logFields[key] = value
	}
	log.WithFields(logFields).Warn(message)
}
