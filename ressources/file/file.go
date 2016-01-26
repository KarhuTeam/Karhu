package file

import (
	"archive/zip"
	"github.com/gotoolz/env"
	"io"
	"log"
	"os"
	"path/filepath"
)

type storageDriver interface {
	store(string, string, []byte) (string, error)
	get(string) ([]byte, error)
}

var driver storageDriver

func init() {

	switch env.GetDefault("STORAGE_DRIVER", "fs") {
	case "fs":
		driver = newFSDriver()
	default:
		panic("Invalid storage driver")
	}
}

func Store(directory, name string, data []byte) (string, error) {

	return driver.store(directory, name, data)
}

func Get(path string) ([]byte, error) {

	return driver.get(path)
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {

			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				log.Fatal(err)
				return err
			}

		} else {

			// Ensure directory exist
			if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				log.Fatal(err)
				return err
			}

			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
