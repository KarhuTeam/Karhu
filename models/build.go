package models

import (
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/ressources/file"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"
)

type buildMapper struct{}

var BuildMapper = &buildMapper{}

const buildCollection = "build"

func init() {
	col := C(buildCollection)
	defer col.Database.Session.Close()

	// App Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"application_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Build struct {
	Id            bson.ObjectId         `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId         `json:"application_id" bson:"application_id"`
	CommitHash    string                `json:"commit_hash" bson:"commit_id"`
	FilePath      string                `json:"file_path" bson:"file_path"`
	CreatedAt     time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at" bson:"updated_at"`
	RuntimeCfg    *RuntimeConfiguration `json:"-" bson:"runtime_cfg"`
}

type Builds []*Build

func (b *Build) GetApplication() (*Application, error) {

	return ApplicationMapper.FetchOne(b.ApplicationId.Hex())
}

func (b *Build) AttachFile(f multipart.File) error {

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New(errors.Error{
			Label: "internal_error",
			Field: "file",
			Text:  err.Error(),
		})
	}

	if contentType := http.DetectContentType(data); contentType != "application/zip" {
		return errors.New(errors.Error{
			Label: "internal_error",
			Field: "file",
			Text:  "Bad content-type, want application/zip have " + contentType,
		})
	}

	if b.FilePath, err = file.Store("builds", b.Id.Hex()+".zip", data); err != nil {
		return err
	}

	// Check ident
	if err := b.readRuntimeConfig(data); err != nil {
		return err
	}

	return nil
}

func (b *Build) readRuntimeConfig(data []byte) error {

	// Temp work dir
	tmpPath, err := ioutil.TempDir("", bson.NewObjectId().Hex())
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpPath)

	// Write zip file
	zipPath := path.Join(tmpPath, "karhu.zip")
	if err := ioutil.WriteFile(zipPath, data, 0644); err != nil {
		return err
	}

	// Unzip file
	if err := file.Unzip(zipPath, tmpPath); err != nil {
		return err
	}

	// Read karhu file
	data, err = ioutil.ReadFile(path.Join(tmpPath, KARHU_FILE_NAME))
	if err != nil {
		log.Println("models/Build: ReadFile:", err)
		return err
	}

	config := new(RuntimeConfiguration)
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Println("models/Build: Unmarshal:", err)
		return err
	}

	app, err := b.GetApplication()
	if err != nil {
		log.Println("models/Build: GetApplication:", err)
		return err
	}

	// Setup workdir if needed
	if config.Workdir == "" {
		config.Workdir = path.Join(KARHU_DEFAULT_RUNTIME_WORKDIR_BASE, app.Name)
	}

	if config.User == "" {
		config.User = "root"
	}

	if config.Binary != nil && config.Binary.User == "" {
		config.Binary.User = config.User
	}

	if err := config.isValid(); err != nil {
		return err
	}

	b.RuntimeCfg = config

	return nil
}

// // Build creation form
// type BuildCreateForm struct {
// 	Version    string                 `json:"version" valid:"ascii,required"`
// 	CommitHash string                 `json:"commit_hash" valid:"hexadecimal,required"`
// 	CommitUrl  string                 `json:"commit_url" valid:"url,required"`
// 	Tags       []string               `json:"tags" valid:"-"`
// 	Vars       map[string]interface{} `json:"vars" valid:"-"`
// }
//
// // Validator for build creation
// func (f BuildCreateForm) Validate() error {
// 	return validator.Validate(&f)
// }

func (bm *buildMapper) Create(app *Application, commitHash string) *Build {

	return &Build{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		CommitHash:    commitHash,
		CreatedAt:     time.Now(),
	}
}

func (bm *buildMapper) CreateService(app *Application, packages []string) *Build {

	rtmcfg := new(RuntimeConfiguration)
	rtmcfg.Dependencies.FromString(packages)

	return &Build{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		CreatedAt:     time.Now(),
		RuntimeCfg:    rtmcfg,
	}
}

func (bm *buildMapper) Save(b *Build) error {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	return col.Insert(b)
}

func (bm *buildMapper) Update(b *Build) error {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	b.UpdatedAt = time.Now()

	return col.UpdateId(b.Id, b)
}

func (bm *buildMapper) Delete(b *Build) error {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(b.Id)
}

func (bm *buildMapper) FetchAll(app *Application) (Builds, error) {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	var builds Builds
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").All(&builds); err != nil {
		return nil, err
	}

	return builds, nil
}

func (bm *buildMapper) FetchLast(app *Application) (*Build, error) {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	build := new(Build)
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").One(build); err != nil {
		return nil, err
	}

	return build, nil
}

func (bm *buildMapper) FetchOne(app *Application, buildId string) (*Build, error) {

	if !bson.IsObjectIdHex(buildId) {

		return nil, errors.New(errors.Error{
			Label: "invalid_build_id",
			Field: "build_id",
			Text:  "Invalid build id hex",
		})
	}

	col := C(buildCollection)
	defer col.Database.Session.Close()

	build := new(Build)
	if err := col.Find(bson.M{"application_id": app.Id, "_id": bson.ObjectIdHex(buildId)}).One(build); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return build, nil
}
