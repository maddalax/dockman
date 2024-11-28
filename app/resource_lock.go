package app

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

func ResourceStatusLock(locator *service.Locator, resourceId string) *DistributedLock {
	key := fmt.Sprintf("resource-status-lock-%s", resourceId)
	lock := KvFromLocator(locator).NewLock(key, 10*time.Second)
	return lock
}

func ResourcePatchLock(locator *service.Locator, resourceId string) *DistributedLock {
	key := fmt.Sprintf("resource-patch-lock-%s", resourceId)
	lock := KvFromLocator(locator).NewLock(key, 10*time.Second)
	return lock
}

func WithStatusLock[T any](locator *service.Locator, resourceId string, f func(err error) T) T {
	lock := ResourceStatusLock(locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return f(err)
	}
	defer lock.Unlock()
	return f(nil)
}
