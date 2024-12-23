package fluentlogger

import (
	"time"
)

//-----------------------------------------------------------------------------

func generateTimekey() string {
	// Get current time
	now := time.Now()

	// Format time according to your desired granularity
	// Example: 1-minute granularity
	timekey := now.Format("2006-01-02T15:04")

	return timekey
}