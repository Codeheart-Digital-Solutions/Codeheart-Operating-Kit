package portfolio

import "time"

const (
	fileShareRetryAttempts = 40
	fileShareRetryDelay    = 5 * time.Millisecond
)
