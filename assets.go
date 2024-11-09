//go:build !prod
// +build !prod

package main

import (
	"io/fs"
	"paas/internal/embedded"
)

func GetStaticAssets() fs.FS {
	return embedded.NewOsFs()
}
