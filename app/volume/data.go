package volume

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetPersistentVolumePath() string {
	storeDir := "/data/dockman"
	useAbs := true

	if runtime.GOOS == "windows" {
		storeDir = "C:/data/dockman"
	}

	if runtime.GOOS == "darwin" {
		dir, err := os.UserHomeDir()
		if err == nil {
			storeDir = filepath.Join(dir, ".dockman")
		} else {
			storeDir = "./dockman"
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
