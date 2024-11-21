package app

import (
	"dockside/app/logger"
	"dockside/app/util/filekv"
	"fmt"
	"os"
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
		logger.ErrorWithFields("Failed to write server config", err, map[string]interface{}{
			"key":   key,
			"value": value,
		})
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
