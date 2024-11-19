package util

import "time"

func WaitFor(timeout time.Duration, interval time.Duration, predicate func() bool) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			return false
		case <-ticker.C:
			if predicate() {
				return true
			}
		}
	}
}
