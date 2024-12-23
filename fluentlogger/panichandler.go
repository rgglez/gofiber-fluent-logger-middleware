package fluentlogger

import (
	"log"
	"runtime/debug"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/ztrue/tracerr"
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
		if err, ok := r.(error); !ok {
			logData["error"] = tracerr.SprintSource(err)
		}

		// Include stack trace if err is a panic
		logData["stacktrace"] = string(debug.Stack())

		// Send the log to Fluentd asynchronously in a goroutine.
		go func() {
			// Safely attempt to post to Fluentd.
			if postErr := l.safePostToFluentd(l.tag+".panic", logData); postErr != nil {
				// If Fluentd fails, fallback to logging to console (or file).
				log.Printf("Fluentd log failed: %v, using fallback mechanism.", postErr)
			}
		}()
	} else {
		log.Printf("Panic occurred but no logger available: %v", r)
	}
}