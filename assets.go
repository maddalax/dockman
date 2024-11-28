//go:build !prod
// +build !prod

package main

import (
	"dockman/app/embedded"
	"io/fs"
)

func GetStaticAssets() fs.FS {
	return embedded.NewOsFs()
}
