package app

import (
	"github.com/maddalax/htmgo/framework/service"
	"os"
)

type ServiceRegistry struct {
	locator *service.Locator
}

func CreateServiceRegistry(locator *service.Locator) *ServiceRegistry {
	registry := &ServiceRegistry{
		locator: locator,
	}
	service.Set(registry.locator, service.Singleton, func() *ServiceRegistry {
		return registry
	})
	return registry
}

func GetServiceRegistry(locator *service.Locator) *ServiceRegistry {
	return service.Get[ServiceRegistry](locator)
}

func (sr *ServiceRegistry) KvClient() *KvClient {
	return service.Get[KvClient](sr.locator)
}

func (sr *ServiceRegistry) RegisterAgent() {
	agent := NewAgent(sr.locator)
	service.Set[Agent](sr.locator, service.Singleton, func() *Agent {
		return agent
	})
}

func (sr *ServiceRegistry) RegisterServerConfigManager() {
	manager := NewServerConfigManager()
	service.Set[ServerConfigManager](sr.locator, service.Singleton, func() *ServerConfigManager {
		return manager
	})
}

func (sr *ServiceRegistry) RegisterEventHandler() {
	handler := NewEventHandler(sr.locator)
	service.Set[EventHandler](sr.locator, service.Singleton, func() *EventHandler {
		return handler
	})
}

func (sr *ServiceRegistry) RegisterJobRunner() {
	runner := NewIntervalJobRunner(sr.locator)
	service.Set[IntervalJobRunner](sr.locator, service.Singleton, func() *IntervalJobRunner {
		return runner
	})
}

func (sr *ServiceRegistry) RegisterKvClient() {
	client, err := NatsConnect(NatsConnectOptions{
		Host: os.Getenv("NATS_HOST"),
		Port: 4222,
	})
	if err != nil {
		panic(err)
	}
	service.Set[KvClient](sr.locator, service.Singleton, func() *KvClient {
		return client
	})
}

func (sr *ServiceRegistry) RegisterBuilderRegistry() {
	registry := NewBuilderRegistry()
	service.Set[BuilderRegistry](sr.locator, service.Singleton, func() *BuilderRegistry {
		return registry
	})
}

func (sr *ServiceRegistry) RegisterReverseProxy() {
	proxy := CreateReverseProxy(sr.locator)
	service.Set[ReverseProxy](sr.locator, service.Singleton, func() *ReverseProxy {
		return proxy
	})
}

func (sr *ServiceRegistry) RegisterResourceMonitor() {
	monitor := NewMonitor(sr.locator)
	service.Set(sr.locator, service.Singleton, func() *ResourceMonitor {
		return monitor
	})
}

func (sr *ServiceRegistry) RegisterSingleton(proxy func() *ReverseProxy) {
	p := proxy()
	service.Set(sr.locator, service.Singleton, func() *ReverseProxy {
		return p
	})
}

func (sr *ServiceRegistry) RegisterJobMetricsManager() {
	manager := NewJobMetricsManager(sr.locator)
	service.Set(sr.locator, service.Singleton, func() *JobMetricsManager {
		return manager
	})
}

func (sr *ServiceRegistry) GetEventHandler() *EventHandler {
	return service.Get[EventHandler](sr.locator)
}

func (sr *ServiceRegistry) GetJobRunner() *IntervalJobRunner {
	return service.Get[IntervalJobRunner](sr.locator)
}

func (sr *ServiceRegistry) GetBuilderRegistry() *BuilderRegistry {
	return service.Get[BuilderRegistry](sr.locator)
}

func (sr *ServiceRegistry) GetResourceMonitor() *ResourceMonitor {
	return service.Get[ResourceMonitor](sr.locator)
}

func (sr *ServiceRegistry) GetReverseProxy() *ReverseProxy {
	return service.Get[ReverseProxy](sr.locator)
}

func (sr *ServiceRegistry) GetAgent() *Agent {
	return service.Get[Agent](sr.locator)
}

func (sr *ServiceRegistry) GetJobMetricsManager() *JobMetricsManager {
	return service.Get[JobMetricsManager](sr.locator)
}

func (sr *ServiceRegistry) GetServerConfigManager() *ServerConfigManager {
	return service.Get[ServerConfigManager](sr.locator)
}

func (sr *ServiceRegistry) RegisterStartupServices() {
	sr.RegisterKvClient()
	sr.RegisterJobRunner()
	sr.RegisterEventHandler()
	sr.RegisterBuilderRegistry()
	sr.RegisterResourceMonitor()
	sr.RegisterReverseProxy()
	sr.RegisterAgent()
	sr.RegisterJobMetricsManager()
	sr.RegisterServerConfigManager()
}

func (sr *ServiceRegistry) RegisterAgentStartupServices() {
	sr.RegisterKvClient()
	sr.RegisterEventHandler()
	sr.RegisterJobRunner()
	sr.RegisterServerConfigManager()
	sr.RegisterJobMetricsManager()
	sr.RegisterAgent()
}
