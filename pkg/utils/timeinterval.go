package utils

import "time"

// Checks wheter a time interval is valid or not
func IsValidTimeFrame(start time.Time, end time.Time) bool {
	if start.Equal(end) {
		return false
	}
	if start.After(end) {
		return false
	}

	return true
}
