package file

import (
	"github.com/gotoolz/env"
)

type storageDriver interface {
	store(string, string, []byte) (string, error)
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
