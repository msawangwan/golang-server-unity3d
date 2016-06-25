package network

import (
	"io"
	"log"
	"os"
)

/* The server logger, prints in various formats to various outputs. */
type ServerLogger struct {
	info       *log.Logger
	debug      *log.Logger
	trace      *log.Logger
	warning    *log.Logger
	alertFatal *log.Logger
}

/* Create an instance of type ServerLogger. */
func NewServerLogger() *ServerLogger {
	logfile, err := os.OpenFile("log-output.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	return &ServerLogger{
		info:       log.New(os.Stdout, "[INFO] ", log.Ltime),
		debug:      log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime),
		warning:    log.New(io.MultiWriter(logfile, os.Stdout), "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile),
		alertFatal: log.New(io.MultiWriter(logfile, os.Stdout), "[ALERT] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

/* Helper function, for printing strings and errors passed to logger methods. */
func printLog(logger interface{}, msgs []interface{}, isFatal bool) {
	if l, ok := logger.(*log.Logger); ok {
		for _, m := range msgs {
			switch msg := m.(type) {
			case string: // print any strings
				if len(msg) > 0 {
					l.Println("[ ", msg, " ]")
				}
			case error: // print any errors
				if msg != nil {
					if isFatal {
						defer l.Fatal("[FATAL] ", msg)
					} else {
						l.Println("[ERROR] ", msg)
					}
				}
			}
		}
	}
}

/* Non-verbose. */
func (sl *ServerLogger) LogInfo(msgs ...interface{}) {
	printLog(sl.info, msgs, false)
}

/* Semi-verbose. */
func (sl *ServerLogger) LogDebug(msgs ...interface{}) {
	printLog(sl.debug, msgs, false)
}

/* Verbose. */
func (sl *ServerLogger) LogWarning(msgs ...interface{}) {
	printLog(sl.warning, msgs, false)
}

/* Fatal. */
func (sl *ServerLogger) LogFatalAlert(msgs ...interface{}) {
	printLog(sl.alertFatal, msgs, true)
}
