package app

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

func GetStatusLock(locator *service.Locator, resourceId string) *Lock {
	key := fmt.Sprintf("resource-status-lock-%s", resourceId)
	lock := GetClientFromLocator(locator).NewLock(key, 10*time.Second)
	return lock
}

func GetPatchLock(locator *service.Locator, resourceId string) *Lock {
	key := fmt.Sprintf("resource-patch-lock-%s", resourceId)
	lock := GetClientFromLocator(locator).NewLock(key, 10*time.Second)
	return lock
}

func WithStatusLock[T any](locator *service.Locator, resourceId string, f func(err error) T) T {
	lock := GetStatusLock(locator, resourceId)
	err := lock.Lock()
	if err != nil {
		return f(err)
	}
	defer lock.Unlock()
	return f(nil)
}