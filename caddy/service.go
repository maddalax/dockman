package caddy

import (
	"paas/httpjson"
)

func GetConfig() (*Config, error) {
	config, err := httpjson.Get[Config]("http://localhost:2019/config")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ApplyConfig(config *Config) error {
	_, err := httpjson.Post[Config]("http://localhost:2019/load", config)
	if err != nil {
		return err
	}
	return nil
}
