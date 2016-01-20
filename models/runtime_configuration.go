package models

import (
	"errors"
	"strings"
)

const (
	KARHU_FILE_NAME                    = "karhu.yml"
	KARHU_DEFAULT_RUNTIME_WORKDIR_BASE = "/usr/local/karhu"
	KARHU_DEFAULT_RUNTIME_USER         = "root"
)

type StaticFiles []string

func (sf StaticFiles) Src(i int) string {

	s := strings.Split(sf[i], ":")
	return s[0]
}

func (sf StaticFiles) Dest(i int) string {

	s := strings.Split(sf[i], ":")
	if len(s) >= 2 {
		return s[1]
	}
	return s[0]
}

func (sf StaticFiles) Mode(i int) string {

	s := strings.Split(sf[i], ":")
	if len(s) >= 3 {
		return s[2]
	}
	return "0644"
}

type RuntimeConfiguration struct {
	Type    string      `yml:"type" bson:"type"`
	User    string      `yml:"user" bson:"user"`
	Bin     string      `yml:"bin" bson:"bin"`
	Workdir string      `yml:"workdir" bson:"workdir"`
	Static  StaticFiles `yml:"static" bson:"static"`
}

func (c *RuntimeConfiguration) isValid() error {

	switch c.Type {
	case "binary":
		if c.Bin == "" {
			return errors.New("Invalid app runtime bin value")
		}
	default:
		return errors.New("Invalid app type")
	}

	for j := range c.Static {

		if c.Static.Src(j)[0] == '/' ||
			c.Static.Dest(j)[0] == '/' {
			return errors.New("Invalid runtime static files")
		}
	}

	return nil
}
