package application

import (
	"errors"
	"github.com/karhuteam/karhu/ressources/file"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	KARHU_FILE_NAME                    = "karhu.yml"
	KARHU_APP_NAME_REGEXP              = "[a-zA-Z0-9_-]+"
	KARHU_DEFAULT_RUNTIME_WORKDIR_BASE = "/usr/local/karhu"
	KARHU_DEFAULT_RUNTIME_USER         = "root"
)

var appNameRgx = regexp.MustCompile(KARHU_APP_NAME_REGEXP)

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

type IdentifierRuntime struct {
	Type    string      `yml:"type"`
	User    string      `yml:"user"`
	Bin     string      `yml:"bin"`
	Workdir string      `yml:"workdir"`
	Static  StaticFiles `yml:"static"`
}

type Identifier struct {
	Name    string            `yml:"name"`
	Version string            `yml:"version"`
	Runtime IdentifierRuntime `yml:"runtime"`
}

func Read(data []byte) (*Identifier, error) {

	// Temp work dir
	tmpPath, err := ioutil.TempDir("", bson.NewObjectId().Hex())
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpPath)

	// Write zip file
	zipPath := path.Join(tmpPath, "karhu.zip")
	if err := ioutil.WriteFile(zipPath, data, 0644); err != nil {
		return nil, err
	}

	// Unzip file
	if err := file.Unzip(zipPath, tmpPath); err != nil {
		return nil, err
	}

	// Read karhu file
	data, err = ioutil.ReadFile(path.Join(tmpPath, KARHU_FILE_NAME))
	if err != nil {
		log.Println("ressources/application: ReadFile:", err)
		return nil, err
	}

	ident := new(Identifier)
	if err := yaml.Unmarshal(data, ident); err != nil {
		log.Println("ressources/application: Unmarshal:", err)
		return nil, err
	}

	if err := ident.isValid(); err != nil {
		return nil, err
	}

	// Setup workdir if needed
	if ident.Runtime.Workdir == "" {
		ident.Runtime.Workdir = path.Join(KARHU_DEFAULT_RUNTIME_WORKDIR_BASE, ident.Name)
	}

	if ident.Runtime.User == "" {
		ident.Runtime.User = KARHU_DEFAULT_RUNTIME_USER
	}

	return ident, nil
}

func (i *Identifier) isValid() error {

	if !appNameRgx.MatchString(i.Name) {
		return errors.New("Invalid app name: " + KARHU_APP_NAME_REGEXP)
	}

	switch i.Runtime.Type {
	case "binary":
		if i.Runtime.Bin == "" {
			return errors.New("Invalid app runtime bin value")
		}
	default:
		return errors.New("Invalid app type")
	}

	for j := range i.Runtime.Static {

		if i.Runtime.Static.Src(j)[0] == '/' ||
			i.Runtime.Static.Dest(j)[0] == '/' {
			return errors.New("Invalid runtime static files")
		}
	}

	return nil
}
