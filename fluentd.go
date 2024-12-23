package fluentlogger

/**
 * safePostToFluentd safely attempts to send the log to Fluentd.
 */
func (l *Logger) safePostToFluentd(tag string, data map[string]interface{}) error {
	// Attempt to post to Fluentd with a timeout.
	err := l.client.Post(tag, data)
	if err != nil {
		// Fluentd is unreachable, return the error.
		return err
	}

	return nil
}
