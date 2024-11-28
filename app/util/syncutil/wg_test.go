package syncutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWaitGroupWithConcurrency(t *testing.T) {
	wg := NewWaitGroupWithConcurrency(5)
	assert.Equal(t, 5, wg.concurrency)
	assert.NotNil(t, wg.semaphore)

	count := 0
	for i := 0; i < 25; i++ {
		wg.Add()
		go func() {
			defer wg.Done()
			count++
		}()
	}
	wg.Wait()
	assert.Equal(t, 25, count)
}
