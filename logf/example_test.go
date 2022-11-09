package logf_test

import (
	"errors"
	"os"

	"github.com/Prastiwar/Go-flow/logf"
)

func Example() {
	// Create new logger with optional settiings like output, formatter, scope fields
	logger := logf.NewLogger(
		logf.WithOutput(os.Stdout),
		logf.WithFormatter(logf.NewTextFormatter()),
		// logf.WithFields(logf.Fields{logf.LogTime: logf.NewTimeField(time.RFC3339)}), // can add log time to message
	)

	// Create logger based on parent logger with additional scope
	logger = logf.WithScope(
		logger,
		logf.Fields{
			"children": true,
		},
	)

	logger.Error("error message")
	logger.Errorf("error occurred: %v", errors.New("invalid call"))

	logger.Info("info message")
	logger.Infof("count: %v", 1)

	logger.Debug("debug message")
	logger.Debugf("debug message: %v", 1)

	// Output:
	// [ERR] error message {"children":true}
	// [ERR] error occurred: invalid call {"children":true}
	// [INFO] info message {"children":true}
	// [INFO] count: 1 {"children":true}
	// [DEBUG] debug message {"children":true}
	// [DEBUG] debug message: 1 {"children":true}
}
