package debug

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// Enable sets the DEBUG env var to true
// and makes the logger to log at debug level.
func Enable() {
	os.Setenv("DEBUG", "1")
	logrus.SetLevel(logrus.DebugLevel)
}

// IsEnabled checks whether the debug flag is set or not.
func IsEnabled() bool {
	return os.Getenv("DEBUG") != ""
}
