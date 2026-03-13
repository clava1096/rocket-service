package path

import (
	"os"
	"path/filepath"
)

func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get current directory: " + err.Error())
	}

	for {
		_, err = os.Stat(filepath.Join(dir, "go.work"))

		if err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root")
		}

		dir = parent
	}
}
