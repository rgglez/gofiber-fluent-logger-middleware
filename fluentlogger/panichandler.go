package fluentlogger

import (
	"runtime/debug"
	"github.com/ztrue/tracerr"
	fiber "github.com/gofiber/fiber/v2"
)

//-----------------------------------------------------------------------------

/**
 * PanicLogger logs details on panic to Fluentd. It is intended to be used as a
 * StackTraceHandler for GoFiber's recover.
 *
 * c the context
 * e the error
 */
 func (l *Logger) PanicLogger(c *fiber.Ctx, r interface{}) {
	if l.enabled && l.client != nil {
		logData := map[string]interface{}{
			"method":     c.Method(),
			"path":       c.Path(),
			"client_ip":  c.IP(),
			"user_agent": c.Get("User-Agent"),
			"time_key":   generateTimekey(),
		}

		// Optionally, include the details of the error
		//if err, ok := r.(error); !ok {
		//	logData["error"] = tracerr.SprintSource(err)
		//}

		// Include stack trace if err is a panic
		logData["stacktrace"] = string(debug.Stack())

		// Send to Fluentd
		if err := l.client.Post(l.tag+".panic", logData); err != nil {
			tracerr.PrintSource(err)
		}
	}
}