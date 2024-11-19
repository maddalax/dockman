package app

import (
	"github.com/maddalax/htmgo/framework/service"
	"sync"
)

type BuilderRegistry struct {
	lock     sync.Mutex
	builders map[string]*ResourceBuilder
}

func NewBuilderRegistry() *BuilderRegistry {
	return &BuilderRegistry{
		builders: make(map[string]*ResourceBuilder),
	}
}

func (r *BuilderRegistry) GetBuilder(resourceId string, buildId string) *ResourceBuilder {
	r.lock.Lock()
	defer r.lock.Unlock()
	key := resourceId + buildId
	if builder, ok := r.builders[key]; ok {
		return builder
	}
	return nil
}

func (r *BuilderRegistry) SetBuilder(resourceId string, buildId string, builder *ResourceBuilder) {
	r.lock.Lock()
	defer r.lock.Unlock()
	key := resourceId + buildId
	r.builders[key] = builder
}

func (r *BuilderRegistry) ClearBuilder(resourceId string, buildId string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	key := resourceId + buildId
	delete(r.builders, key)
}

func GetBuilderRegistry(locator *service.Locator) *BuilderRegistry {
	return service.Get[BuilderRegistry](locator)
}
