package volume

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetPersistentVolumePath() string {
	storeDir := "/data/dockside"
	useAbs := true

	if runtime.GOOS == "windows" {
		storeDir = "C:/data/dockside"
	}

	if runtime.GOOS == "darwin" {
		dir, err := os.UserHomeDir()
		if err == nil {
			storeDir = filepath.Join(dir, ".dockside")
		} else {
			storeDir = "./dockside"
			useAbs = false
		}
	}

	if !useAbs {
		return storeDir
	}

	abs, err := filepath.Abs(storeDir)

	if err != nil {
		panic(err)
	}

	return abs
}
