package mydumpster

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/op/go-logging"
)

const MydumpsterLogger = "mydumpster"

// Shared across all the app
var log *logging.Logger

// TODO: configuration
func GetLogger(namespace string) *logging.Logger {
	if log == nil {
		log = logging.MustGetLogger(namespace)
		logging.SetFormatter(logging.MustStringFormatter("▶ [%{level:.1s}] %{message}"))
		logBackend := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
		logging.SetBackend(logBackend)
		logBackend.Color = true
		logging.SetLevel(logging.DEBUG, namespace)
		log.Info(fmt.Sprintf("'%s' Log configured", namespace))
	}

	// Don't needed but sometimes is handy
	return log
}
