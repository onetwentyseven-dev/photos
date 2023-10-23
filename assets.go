package photos

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets
var assets embed.FS

const assetsDirectory = "assets"

func AssetFS(environment Environment) fs.FS {

	if !environment.IsProduction() {
		return getLocalFS()
	}

	subFS, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(fmt.Sprintf("fs.Sub: %s", err))
	}

	return subFS

}

func getLocalFS() fs.FS {

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get current working directory: %s", err))
	}

	return os.DirFS(filepath.Join(cwd, assetsDirectory))

}
