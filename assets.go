//go:build !prod
// +build !prod

package main

import (
	"io/fs"
	"paas/app/embedded"
)

func GetStaticAssets() fs.FS {
	return embedded.NewOsFs()
}
