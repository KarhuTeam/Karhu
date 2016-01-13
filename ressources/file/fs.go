package file

import (
	"github.com/gotoolz/env"
	"github.com/gotoolz/errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type fsDriver struct {
	BaseDir string
}

func newFSDriver() storageDriver {

	return &fsDriver{
		BaseDir: env.GetDefault("STORAGE_PATH", "./data"),
	}
}

func (d *fsDriver) store(dir string, name string, data []byte) (string, error) {

	targetDir := path.Join(d.BaseDir, dir)

	// Check for target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Println("ressources/file/fs: failed to create directory:", err)
		return "", errors.New(errors.Error{
			Label: "internal_error",
			Field: "file",
			Text:  err.Error(),
		})
	}

	targetFile := path.Join(targetDir, name)

	if err := ioutil.WriteFile(targetFile, data, 0644); err != nil {
		return "", errors.New(errors.Error{
			Label: "internal_error",
			Field: "file",
			Text:  err.Error(),
		})
	}

	return targetFile, nil
}

func (d *fsDriver) get(path string) ([]byte, error) {

	return ioutil.ReadFile(path)
}
