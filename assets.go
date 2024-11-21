//go:build !prod
// +build !prod

package main

import (
	"dockside/app/embedded"
	"io/fs"
)

func GetStaticAssets() fs.FS {
	return embedded.NewOsFs()
}
