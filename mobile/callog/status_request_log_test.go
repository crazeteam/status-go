package callog

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/log"

	"github.com/stretchr/testify/require"

	"github.com/status-im/status-go/logutils/requestlog"
)

func TestRemoveSensitiveInfo(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic test",
			input:    `{"username":"user1","password":"secret123","mnemonic":"mnemonic123 xyz"}`,
			expected: `{"username":"user1","password":"***","mnemonic":"***"}`,
		},
		{
			name:     "uppercase password field",
			input:    `{"USERNAME":"user1","PASSWORD":"secret123"}`,
			expected: `{"USERNAME":"user1","PASSWORD":"***"}`,
		},
		{
			name:     "password field with spaces",
			input:    `{"username":"user1", "password" : "secret123"}`,
			expected: `{"username":"user1", "password":"***"}`,
		},
		{
			name:     "multiple password fields",
			input:    `{"password":"secret123","data":{"nested_password":"nested_secret"}}`,
			expected: `{"password":"***","data":{"nested_password":"***"}}`,
		},
		{
			name:     "no password field",
			input:    `{"username":"user1","email":"user1@example.com"}`,
			expected: `{"username":"user1","email":"user1@example.com"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeSensitiveInfo(tc.input)
			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}

func TestCall(t *testing.T) {
	// Create a temporary file for logging
	tempLogFile, err := os.CreateTemp(t.TempDir(), "TestCall*.log")
	require.NoError(t, err)

	// Enable request logging
	logger, err := requestlog.CreateRequestLogger(tempLogFile.Name())
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test case 1: Normal execution
	testFunc := func(param string) string {
		return "test result: " + param
	}
	testParam := "test input"
	expectedResult := "test result: test input"

	result := CallWithResponse(logger, testFunc, testParam)

	// Check the result
	if result != expectedResult {
		t.Errorf("Expected result %s, got %s", expectedResult, result)
	}

	// Read the log file
	logData, err := os.ReadFile(tempLogFile.Name())
	require.NoError(t, err)
	logOutput := string(logData)

	// Check if the log contains expected information
	expectedLogParts := []string{getShortFunctionName(testFunc), "params", testParam, "resp", expectedResult}
	for _, part := range expectedLogParts {
		if !strings.Contains(logOutput, part) {
			t.Errorf("Log output doesn't contain expected part: %s", part)
		}
	}

	// Create a mock logger to capture log output
	mockLogger := log.New()
	mockLogger.SetHandler(log.FuncHandler(func(r *log.Record) error {
		logOutput += r.Msg + fmt.Sprintf("%s", r.Ctx...)
		return nil
	}))

	// Test case 2: Panic -> recovery -> re-panic
	oldRootHandler := log.Root().GetHandler()
	defer log.Root().SetHandler(oldRootHandler)
	log.Root().SetHandler(mockLogger.GetHandler())
	// Clear log output for next test
	logOutput = ""

	e := "test panic"
	panicFunc := func() {
		panic(e)
	}

	require.PanicsWithValue(t, e, func() {
		Call(logger, panicFunc)
	})

	// Check if the panic was logged
	if !strings.Contains(logOutput, "panic found in call") {
		t.Errorf("Log output doesn't contain panic information")
	}
	if !strings.Contains(logOutput, e) {
		t.Errorf("Log output doesn't contain panic message")
	}
	if !strings.Contains(logOutput, "stacktrace") {
		t.Errorf("Log output doesn't contain stacktrace")
	}
}

func initializeApplication(requestJSON string) string {
	return ""
}

func TestGetFunctionName(t *testing.T) {
	fn := getShortFunctionName(initializeApplication)
	require.Equal(t, "initializeApplication", fn)
}