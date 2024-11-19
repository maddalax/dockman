package util

import (
	"github.com/maddalax/htmgo/framework/h"
	"time"
)

// DelayedPartial delays the execution of a partial by the given duration
// if the partial takes longer than the delay, the delay is ignored
func DelayedPartial(delay time.Duration, f func() *h.Partial) *h.Partial {
	now := time.Now()
	p := f()
	elapsed := time.Since(now)
	if elapsed < delay {
		time.Sleep(delay - elapsed)
	}
	return p
}
