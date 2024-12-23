package fluentlogger

/*
Copyright 2024 Rodolfo González González

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"log"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/ztrue/tracerr"
)

//*****************************************************************************

type LoggerConfig struct {
	Enabled bool   // whether the middleware is enabled
	Tag     string // the tag to be used for the messages
}

//-----------------------------------------------------------------------------

/**
 * Logger is a struct that holds the Fluentd logger instance and configuration
 */
type Logger struct {
	client  *fluent.Fluent
	tag     string
	enabled bool
}

//-----------------------------------------------------------------------------

/**
 * New initializes a Fluentd logger and returns a middleware
 */
func New(config LoggerConfig, client *fluent.Fluent) *Logger {
	return &Logger{
		client:  client,
		tag:     config.Tag,
		enabled: config.Enabled,
	}
}

//-----------------------------------------------------------------------------

/**
 * Logger logs each request to Fluentd
 */
func (l *Logger) Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // Process the request
		latency := time.Since(start)

		if l.enabled && l.client != nil {
			// Log data to Fluentd
			logData := map[string]interface{}{
				"method":        c.Method(),
				"path":          c.Path(),
				"status":        c.Response().StatusCode(),
				"latency_ms":    latency.Milliseconds(),
				"client_ip":     c.IP(),
				"user_agent":    c.Get("User-Agent"),
				"response_size": len(c.Response().Body()),
				"time_key":      generateTimekey(),
			}
			if err != nil {
				logData["error"] = tracerr.SprintSource(err)
			}

			// Send the log to Fluentd asynchronously in a goroutine.
			go func() {
				// Safely attempt to post to Fluentd.
				if postErr := l.safePostToFluentd(logData); postErr != nil {
					// If Fluentd fails, fallback to logging to console (or file).
					log.Printf("Fluentd log failed: %v, using fallback mechanism.", postErr)
				}
			}()
		}

		return err
	}
}

//-----------------------------------------------------------------------------

/**
 * safePostToFluentd safely attempts to send the log to Fluentd.
 */
func (l *Logger) safePostToFluentd(data map[string]interface{}) error {
	// Attempt to post to Fluentd with a timeout.
	err := l.client.Post(l.tag, data)
	if err != nil {
		// Fluentd is unreachable, return the error.
		return err
	}

	return nil
}
