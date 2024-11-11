package caddy

import (
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
)

func Run() error {
	config := BuildConfig()
	return ApplyConfig(config)
}
