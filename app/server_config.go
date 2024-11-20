package app

import (
	"fmt"
	"log/slog"
	"os"
	"paas/app/util/filekv"
	"path/filepath"
)

type ServerConfigManager struct {
	cache  map[string]string
	locker *filekv.FileLocker
}

func NewServerConfigManager() *ServerConfigManager {
	return &ServerConfigManager{
		cache:  nil,
		locker: &filekv.FileLocker{},
	}
}

func (m *ServerConfigManager) ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "./"
	}
	abs, err := filepath.Abs(filepath.Join(home, fmt.Sprintf(".%s", AppName), "server_config.json"))
	if err != nil {
		return "./server_config.json"
	}
	return abs
}

func (m *ServerConfigManager) WriteConfig(key string, value string) {
	err := filekv.WriteKeyValue(m.ConfigPath(), key, value, m.locker)
	if err != nil {
		slog.Error("Failed to write server config", slog.String("key", key), slog.String("value", value), slog.String("error", err.Error()))
	}
}

func (m *ServerConfigManager) ClearCache() {
	m.cache = nil
}

func (m *ServerConfigManager) GetConfig(key string) string {
	if m.cache == nil {
		config, err := filekv.ReadKeyValues(m.ConfigPath())
		if err != nil {
			return ""
		}
		m.cache = config
	}
	return m.cache[key]
}
