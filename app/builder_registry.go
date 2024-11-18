package app

import "sync"

var builders = make(map[string]*ResourceBuilder)
var lock = sync.Mutex{}

func GetBuilder(resourceId string, buildId string) *ResourceBuilder {
	lock.Lock()
	defer lock.Unlock()
	key := resourceId + buildId
	if builder, ok := builders[key]; ok {
		return builder
	}
	return nil
}

func SetBuilder(resourceId string, buildId string, builder *ResourceBuilder) {
	lock.Lock()
	defer lock.Unlock()
	key := resourceId + buildId
	builders[key] = builder
}

func ClearBuilder(resourceId string, buildId string) {
	lock.Lock()
	defer lock.Unlock()
	key := resourceId + buildId
	delete(builders, key)
}
