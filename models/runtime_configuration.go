package models

import (
	"errors"
)

const (
	KARHU_FILE_NAME                    = "karhu.yml"
	KARHU_DEFAULT_RUNTIME_WORKDIR_BASE = "/usr/local/karhu"
	KARHU_DEFAULT_RUNTIME_USER         = "root"
)

type BinaryConfiguration struct {
	User string `yml:"user" bson:"user"`
	Bin  string `yml:"bin" bson:"bin"`
}

type StaticConfiguration struct {
	Src  string `yml:"src" bson:"src"`
	Dest string `yml:"dest" bson:"dest"`
	Mode string `yml:"mode" bson:"mode"`
}

type DependenciesConfiguration struct {
	Name string `yml:"name" bson:"name"`
}
type DependenciesConfigurations []*DependenciesConfiguration

func (dc DependenciesConfigurations) FromString(pkgs []string) DependenciesConfigurations {

	dc = nil

	for _, pkg := range pkgs {
		dc = append(dc, &DependenciesConfiguration{
			Name: pkg,
		})
	}

	return dc
}

func (dc DependenciesConfigurations) ToString() []string {

	var pkgs []string

	for _, dep := range dc {
		pkgs = append(pkgs, dep.Name)
	}

	return pkgs
}

type RuntimeConfiguration struct {
	Workdir      string                     `yml:"workdir" bson:"workdir"`
	User         string                     `yml:"user" bson:"user"`
	Binary       *BinaryConfiguration       `yml:"binary" bson:"binary"`
	Static       []*StaticConfiguration     `yml:"static" bson:"static"`
	Dependencies DependenciesConfigurations `yml:"dependencies" bson:"dependencies"`
}

func (c *RuntimeConfiguration) isValid() error {

	if c.Binary != nil {
		if c.Binary.Bin == "" {
			return errors.New("Invalid app runtime bin value")
		}
	}

	for _, s := range c.Static {

		if len(s.Src) > 0 && s.Src[0] == '/' ||
			len(s.Dest) > 0 && s.Dest[0] == '/' {
			return errors.New("Invalid runtime static files")
		}
	}

	return nil
}
